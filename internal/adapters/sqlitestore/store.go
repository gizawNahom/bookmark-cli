// Package sqlitestore is the SQLiteBookmarkStore driven adapter (ADR-002/ADR-004/ADR-007):
// implements both ports.BookmarkReader and ports.BookmarkWriter against a WAL-mode SQLite file
// via the pure-Go modernc.org/sqlite driver (no cgo -- keeps the single-static-binary story).
package sqlitestore

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/gizawNahom/bookmark-cli/internal/core"
)

// Store implements ports.BookmarkReader, ports.BookmarkWriter, and ports.Prober.
type Store struct {
	Path string // absolute path to the bookmarks.db file under ${data_dir}
}

// NewStore constructs a Store for the given database file path. Does not open a connection or
// touch the filesystem -- that happens in Probe(), per "wire then probe then use".
func NewStore(path string) *Store {
	return &Store{Path: path}
}

const schema = `
CREATE TABLE IF NOT EXISTS bookmarks (
	id       TEXT PRIMARY KEY,
	url      TEXT NOT NULL UNIQUE,
	tag      TEXT NOT NULL DEFAULT '',
	saved_at TEXT NOT NULL
);
`

// open creates the data directory if needed, opens a WAL-mode connection with a busy_timeout
// guard (~5000ms, ADR-004), and ensures the schema exists.
func (s *Store) open() (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return nil, fmt.Errorf("creating data directory: %w", err)
	}

	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)", s.Path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("initializing schema: %w", err)
	}

	return db, nil
}

// Probe verifies: (1) the DB file/directory is creatable and writable, (2) WAL mode actually
// engages (PRAGMA journal_mode returns "wal", not silently falling back), (3) a write+fsync+
// read-back round-trip on a throwaway row actually persists (brief.md Section 11/12).
func (s *Store) Probe() error {
	db, err := s.open()
	if err != nil {
		return err
	}
	defer db.Close()

	var mode string
	if err := db.QueryRow("PRAGMA journal_mode").Scan(&mode); err != nil {
		return fmt.Errorf("checking journal_mode: %w", err)
	}
	if !strings.EqualFold(mode, "wal") {
		return fmt.Errorf("expected WAL journal mode, got %q", mode)
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS probe_check (id INTEGER PRIMARY KEY, value TEXT)"); err != nil {
		return fmt.Errorf("preparing probe table: %w", err)
	}
	if _, err := db.Exec("INSERT INTO probe_check (value) VALUES (?)", "probe"); err != nil {
		return fmt.Errorf("probe write: %w", err)
	}
	var readBack string
	if err := db.QueryRow("SELECT value FROM probe_check ORDER BY id DESC LIMIT 1").Scan(&readBack); err != nil {
		return fmt.Errorf("probe read-back: %w", err)
	}
	if readBack != "probe" {
		return fmt.Errorf("probe round-trip mismatch: wrote %q, read %q", "probe", readBack)
	}
	if _, err := db.Exec("DELETE FROM probe_check"); err != nil {
		return fmt.Errorf("probe cleanup: %w", err)
	}
	return nil
}

// FindByID resolves a single bookmark by its short id.
func (s *Store) FindByID(id string) (core.Record, bool, error) {
	db, err := s.open()
	if err != nil {
		return core.Record{}, false, err
	}
	defer db.Close()

	rec, err := scanRecord(db.QueryRow("SELECT id, url, tag, saved_at FROM bookmarks WHERE id = ?", id))
	if err == sql.ErrNoRows {
		return core.Record{}, false, nil
	}
	if err != nil {
		return core.Record{}, false, err
	}
	return rec, true, nil
}

// Search returns candidate records for the given query. Ranking/filtering against the query is
// performed by the caller (core.RankMatches, wired in step 02-01) -- this adapter's job is
// fetching candidates, not ranking them.
func (s *Store) Search(query string) ([]core.Record, error) {
	return s.All()
}

// All returns every stored record (used by core.RankMatches callers and bm stats-adjacent flows).
func (s *Store) All() ([]core.Record, error) {
	db, err := s.open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, url, tag, saved_at FROM bookmarks ORDER BY saved_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []core.Record
	for rows.Next() {
		rec, err := scanRecord(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, rows.Err()
}

// Execute performs the single mutation described by plan (New | Duplicate | TagUpdate) inside a
// WAL-mode, busy_timeout-guarded transaction. This is the only method on the entire adapter that
// writes to the bookmarks table.
func (s *Store) Execute(plan core.SavePlan) (core.Record, error) {
	db, err := s.open()
	if err != nil {
		return core.Record{}, err
	}
	defer db.Close()

	switch plan.Kind {
	case core.PlanNew:
		return insertNew(db, plan)
	case core.PlanDuplicate, core.PlanTagUpdate:
		// No mutation: the pure core already decided this URL is already saved. Look up and
		// return the existing record unchanged -- Execute never writes a second row for a plan
		// that isn't PlanNew (ADR-006 / v3.15.1 precedent: duplicate detection must never
		// silently write).
		rec, err := scanRecord(db.QueryRow(
			"SELECT id, url, tag, saved_at FROM bookmarks WHERE id = ?", plan.ExistingID,
		))
		if err != nil {
			return core.Record{}, fmt.Errorf("looking up existing bookmark: %w", err)
		}
		return rec, nil
	default:
		return core.Record{}, fmt.Errorf("sqlitestore: unsupported plan kind %q", plan.Kind)
	}
}

// maxIDCollisionRetries bounds insertNew's retry loop. newBookmarkID's random space makes a
// single collision already unlikely at brief.md's stated 10,000-bookmark scale; this only
// guards against the residual chance of one, not a normal/expected occurrence.
const maxIDCollisionRetries = 5

// insertNew writes a PlanNew record, regenerating the id and retrying on a bookmarks.id
// collision -- newBookmarkID draws from a random space, so a fresh id has no reason to collide
// again. A non-id constraint failure (e.g. the url UNIQUE constraint, which PlanNew should never
// hit under normal operation) is returned immediately rather than retried, since retrying with a
// new id wouldn't change the url.
func insertNew(db *sql.DB, plan core.SavePlan) (core.Record, error) {
	savedAt := time.Now().UTC()
	var lastErr error
	for range maxIDCollisionRetries {
		id := newBookmarkID()
		if _, err := db.Exec(
			"INSERT INTO bookmarks (id, url, tag, saved_at) VALUES (?, ?, ?, ?)",
			id, plan.URL, plan.Tag, savedAt.Format(time.RFC3339),
		); err != nil {
			if !strings.Contains(err.Error(), "bookmarks.id") {
				return core.Record{}, fmt.Errorf("inserting new bookmark: %w", err)
			}
			lastErr = err
			continue
		}
		return core.Record{ID: id, URL: plan.URL, Tag: plan.Tag, SavedAt: savedAt}, nil
	}
	return core.Record{}, fmt.Errorf("inserting new bookmark: id collided %d times in a row: %w", maxIDCollisionRetries, lastErr)
}

// scanner abstracts over *sql.Row and *sql.Rows so scanRecord serves both FindByID and All.
type scanner interface {
	Scan(dest ...any) error
}

func scanRecord(row scanner) (core.Record, error) {
	var rec core.Record
	var savedAt string
	if err := row.Scan(&rec.ID, &rec.URL, &rec.Tag, &savedAt); err != nil {
		return core.Record{}, err
	}
	parsed, err := time.Parse(time.RFC3339, savedAt)
	if err != nil {
		return core.Record{}, fmt.Errorf("parsing saved_at: %w", err)
	}
	rec.SavedAt = parsed
	return rec, nil
}

// newBookmarkID generates a short, human-typeable hex id (e.g. "a1b2c3d4"). 4 random bytes (32
// bits) keeps the birthday-bound collision probability at brief.md's stated 10,000-bookmark scale
// under ~1.2% -- insertNew's retry loop covers that residual chance; a 2-byte id (the original
// size) put collisions above 50% by ~300 saves and had no retry, so saves silently started failing
// well before reaching any real-world store size.
func newBookmarkID() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand.Read on the standard reader does not fail in practice; a timestamp-derived
		// fallback keeps this function total without panicking the imperative shell.
		return hex.EncodeToString([]byte(time.Now().Format("20060102150405.000000")))
	}
	return hex.EncodeToString(b)
}
