// Command probecheck is ADR-007's structural (layer 2) enforcement: every concrete type that
// implements one of internal/ports' driven-adapter interfaces (BookmarkReader, BookmarkWriter,
// BackupService, UsageLogger) must also implement ports.Prober (Probe() error). Compile-time
// interface embedding alone doesn't catch this for BookmarkReader/BookmarkWriter, which don't
// embed Prober -- hence this go/ast (go/types) walk, run in CI against `./...`.
//
// Adapters are discovered by structural interface satisfaction (go/types.Implements) against
// whatever packages are passed on the command line, never by a hardcoded type-name list, so a
// new driven adapter is covered automatically the moment it's added.
package main

import (
	"fmt"
	"go/types"
	"os"

	"golang.org/x/tools/go/packages"
)

const portsPkgPath = "github.com/gizawNahom/bookmark-cli/internal/ports"

func main() {
	patterns := os.Args[1:]
	if len(patterns) == 0 {
		patterns = []string{"./..."}
	}

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo |
			packages.NeedSyntax | packages.NeedImports | packages.NeedDeps,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "probecheck: failed to load packages: %v\n", err)
		os.Exit(1)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	portsPkg := findPackage(pkgs, portsPkgPath)
	if portsPkg == nil {
		fmt.Fprintf(os.Stderr, "probecheck: package %q not found among loaded packages -- "+
			"run against a pattern that includes it (e.g. ./...)\n", portsPkgPath)
		os.Exit(1)
	}

	driverIfaces, err := driverInterfaces(portsPkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "probecheck: %v\n", err)
		os.Exit(1)
	}
	prober, err := lookupInterface(portsPkg, "Prober")
	if err != nil {
		fmt.Fprintf(os.Stderr, "probecheck: %v\n", err)
		os.Exit(1)
	}

	var violations []string
	for _, pkg := range pkgs {
		if pkg.Types == nil || pkg.PkgPath == portsPkgPath {
			continue
		}
		for _, name := range pkg.Types.Scope().Names() {
			obj, ok := pkg.Types.Scope().Lookup(name).(*types.TypeName)
			if !ok {
				continue
			}
			named, ok := obj.Type().(*types.Named)
			if !ok {
				continue
			}
			if _, isInterface := named.Underlying().(*types.Interface); isInterface {
				continue
			}

			candidate := candidateType(named)
			var implementsDriverIface string
			for ifaceName, iface := range driverIfaces {
				if types.Implements(candidate, iface) {
					implementsDriverIface = ifaceName
					break
				}
			}
			if implementsDriverIface == "" {
				continue
			}
			if !types.Implements(candidate, prober) {
				violations = append(violations, fmt.Sprintf(
					"%s.%s implements ports.%s but not ports.Prober (missing Probe() error)",
					pkg.PkgPath, name, implementsDriverIface))
			}
		}
	}

	if len(violations) > 0 {
		fmt.Fprintln(os.Stderr, "probecheck: ADR-007 structural check failed:")
		for _, v := range violations {
			fmt.Fprintln(os.Stderr, "  "+v)
		}
		os.Exit(1)
	}

	fmt.Println("probecheck: OK -- every driven-adapter-interface implementation defines Probe() error")
}

// candidateType always checks *T: Go's method-set rule makes *T's method set a superset of T's
// (it includes both pointer- and value-receiver methods), so checking *T alone is safe regardless
// of which receiver form the type's methods were declared with -- matching how this codebase's
// adapters are actually constructed and used (New*() returns a pointer).
func candidateType(named *types.Named) types.Type {
	return types.NewPointer(named)
}

func driverInterfaces(portsPkg *packages.Package) (map[string]*types.Interface, error) {
	names := []string{"BookmarkReader", "BookmarkWriter", "BackupService", "UsageLogger"}
	result := make(map[string]*types.Interface, len(names))
	for _, name := range names {
		iface, err := lookupInterface(portsPkg, name)
		if err != nil {
			return nil, err
		}
		result[name] = iface
	}
	return result, nil
}

func lookupInterface(pkg *packages.Package, name string) (*types.Interface, error) {
	obj := pkg.Types.Scope().Lookup(name)
	if obj == nil {
		return nil, fmt.Errorf("ports.%s not found in %s", name, pkg.PkgPath)
	}
	named, ok := obj.Type().(*types.Named)
	if !ok {
		return nil, fmt.Errorf("ports.%s is not a named type", name)
	}
	iface, ok := named.Underlying().(*types.Interface)
	if !ok {
		return nil, fmt.Errorf("ports.%s is not an interface", name)
	}
	return iface, nil
}

func findPackage(pkgs []*packages.Package, pkgPath string) *packages.Package {
	for _, pkg := range pkgs {
		if pkg.PkgPath == pkgPath {
			return pkg
		}
	}
	return nil
}
