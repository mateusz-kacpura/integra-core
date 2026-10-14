package contracts

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestRemovedLegacyPackagesStayRemoved(t *testing.T) {
	root := ".."
	for _, name := range []string{"domain", "license"} {
		if _, err := os.Stat(filepath.Join(root, name)); !os.IsNotExist(err) {
			t.Fatalf("legacy package %s must remain absent; stat error: %v", name, err)
		}
	}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, imported := range file.Imports {
			dependency, err := strconv.Unquote(imported.Path.Value)
			if err != nil {
				return err
			}
			if !strings.HasSuffix(path, "_test.go") && (dependency == "google.golang.org/grpc" || strings.HasPrefix(dependency, "google.golang.org/grpc/")) {
				t.Errorf("%s: shared adapters must not directly import gRPC", path)
			}
			for _, name := range []string{"domain", "license"} {
				legacy := "github.com/mateusz-kacpura/integra-core/pkg/" + name
				if dependency == legacy || strings.HasPrefix(dependency, legacy+"/") {
					t.Errorf("%s imports removed package %s", path, dependency)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestParentModuleDoesNotDependOnNestedAdapters(t *testing.T) {
	root := filepath.Clean("../..")
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != root {
				if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
					return filepath.SkipDir
				}
				for _, name := range []string{"auth", "kafka", "database"} {
					if filepath.Base(path) == name && filepath.Dir(path) == filepath.Join(root, "pkg") {
						return filepath.SkipDir
					}
				}
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, spec := range file.Imports {
			dependency, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return err
			}
			for _, name := range []string{"auth", "kafka", "database"} {
				nested := "github.com/mateusz-kacpura/integra-core/pkg/" + name
				if dependency == nested || strings.HasPrefix(dependency, nested+"/") {
					t.Errorf("%s reintroduces parent-to-child dependency %s", path, dependency)
				}
			}
			rel, _ := filepath.Rel(root, path)
			if filepath.Dir(rel) == filepath.Join("pkg", "rbac") && !strings.HasSuffix(path, "_test.go") && dependency != "sync" && dependency != "sync/atomic" && dependency != "strings" {
				t.Errorf("RBAC core has an adapter dependency: %s", dependency)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
