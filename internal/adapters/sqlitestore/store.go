// Package sqlitestore is the SQLiteBookmarkStore driven adapter (ADR-002/ADR-004/ADR-007):
// implements both ports.BookmarkReader and ports.BookmarkWriter against a WAL-mode SQLite file
// via the pure-Go modernc.org/sqlite driver (no cgo -- keeps the single-static-binary story).
//
// SCAFFOLD: true -- DISTILL RED scaffold (nw-acceptance-designer). Replaced with a real
// implementation during DELIVER GREEN phase.
package sqlitestore

import "bookmark-cli/internal/core"

// Store implements ports.BookmarkReader, ports.BookmarkWriter, and ports.Prober.
type Store struct {
	Path string // absolute path to the bookmarks.db file under ${data_dir}
}

// NewStore constructs a Store for the given database file path. Does not open a connection or
// touch the filesystem -- that happens in Probe(), per "wire then probe then use".
func NewStore(path string) *Store {
	return &Store{Path: path}
}

// Probe verifies: (1) the DB file/directory is creatable and writable, (2) WAL mode actually
// engages (PRAGMA journal_mode returns "wal", not silently falling back), (3) a write+fsync+
// read-back round-trip on a throwaway row actually persists (brief.md Section 11).
func (s *Store) Probe() error {
	panic("sqlitestore.Store.Probe not yet implemented -- RED scaffold")
}

// FindByID resolves a single bookmark by its short id.
func (s *Store) FindByID(id string) (core.Record, bool, error) {
	panic("sqlitestore.Store.FindByID not yet implemented -- RED scaffold")
}

// Search performs an FTS5 ranked search across URL/tag content.
func (s *Store) Search(query string) ([]core.Record, error) {
	panic("sqlitestore.Store.Search not yet implemented -- RED scaffold")
}

// All returns every stored record (used by core.RankMatches callers and bm stats-adjacent flows).
func (s *Store) All() ([]core.Record, error) {
	panic("sqlitestore.Store.All not yet implemented -- RED scaffold")
}

// Execute performs the single mutation described by plan (New | Duplicate | TagUpdate) inside a
// WAL-mode, busy_timeout-guarded transaction. This is the only method on the entire adapter that
// writes to the bookmarks table.
func (s *Store) Execute(plan core.SavePlan) (core.Record, error) {
	panic("sqlitestore.Store.Execute not yet implemented -- RED scaffold")
}
