package core

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// forbidden are imports the core must never take: the hexagon depends on its
// ports, and adapters depend on the core — never the other way round.
var forbidden = []string{
	"aprilpollo/internal/adapters",
	"github.com/gofiber/",
	"gorm.io/",
	"github.com/jackc/",
	"github.com/redis/",
	"github.com/go-playground/",
	"net/http",
	"database/sql",
}

func TestCoreDependsOnNoAdapter(t *testing.T) {
	fset := token.NewFileSet()
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}
		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, imp := range f.Imports {
			p, _ := strconv.Unquote(imp.Path.Value)
			for _, bad := range forbidden {
				if strings.HasPrefix(p, bad) {
					t.Errorf("%s imports %q; the core must not depend on %s", path, p, bad)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
