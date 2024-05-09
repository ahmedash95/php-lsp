package workspacescanner_test

import (
	workspacescanner "ahmedash95/php-lsp-server/pkg/workspace_scanner"
	"fmt"
	"os"
	"path"
	"testing"
)

func TestNewScanner(t *testing.T) {
	gp := os.Getenv("GOPATH")
	path := path.Join(gp, "src/ahmedash95/php-lsp-server")

	tests := map[string]struct {
		path     string
		ext      []string
		want     *workspacescanner.Scanner
		expected []string
	}{
		"Can get all php files in a directory": {
			fmt.Sprintf("%s/tests/pkg/workspace_scanner/tmp", path),
			[]string{".php"},
			&workspacescanner.Scanner{
				Path: fmt.Sprintf("%s/tests/pkg/workspace_scanner/tmp", path),
			},
			[]string{
				"file.php",
				"src/assets/file1.php",
				"src/file1.php",
				"src/file2.php",
			},
		},
		"Can get svg files in a directory": {
			fmt.Sprintf("%s/tests/pkg/workspace_scanner/tmp", path),
			[]string{".svg"},
			&workspacescanner.Scanner{
				Path: fmt.Sprintf("%s/tests/pkg/workspace_scanner/tmp", path),
			},
			[]string{
				"src/assets/logo.svg",
			},
		},
		"Can get all files in a directory": {
			fmt.Sprintf("%s/tests/pkg/workspace_scanner/tmp", path),
			[]string{".php", ".svg"},
			&workspacescanner.Scanner{
				Path: fmt.Sprintf("%s/tests/pkg/workspace_scanner/tmp", path),
			},
			[]string{
				"file.php",
				"src/assets/file1.php",
				"src/assets/logo.svg",
				"src/file1.php",
				"src/file2.php",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := workspacescanner.NewScanner(tt.path)
			files := s.Scan(tt.ext)

			if len(files) != len(tt.expected) {
				t.Errorf("Expected %d files, got %d", len(tt.expected), len(files))
			}

			for i, file := range tt.expected {
				if files[i] != file {
					t.Errorf("Expected %s, got %s", file, files[i])
				}
			}
		})
	}
}
