package workspacescanner

import (
	"ahmedash95/php-lsp-server/pkg/logger"
	"os"
	"path/filepath"
)

type Scanner struct {
	Path string
}

func NewScanner(path string) *Scanner {
	return &Scanner{
		Path: path,
	}
}

func (s *Scanner) GetFileContent(path string) string {
	absolutepath := filepath.Join(s.Path, path)
	content, err := os.ReadFile(absolutepath)
	if err != nil {
		logger.GetLogger().Printf("Error reading file: %s", err)
		return ""
	}
	return string(content)
}

func (s *Scanner) Scan(ext []string) []string {
	files := make([]string, 0)
	s.walkDir(s.Path, ext, &files)
	return files
}

func (s *Scanner) walkDir(dirPath string, ext []string, files *[]string) {
	extMap := make(map[string]bool)
	for _, e := range ext {
		extMap[e] = true
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // propagate the error further (e.g., permissions error)
		}
		if !info.IsDir() {
			if extMap[filepath.Ext(info.Name())] {
				relativePath, err := filepath.Rel(s.Path, path)
				if err != nil {
					logger.GetLogger().Printf("Error getting relative path: %s", err)
				} else {
					*files = append(*files, relativePath)
				}
			}
		}
		return nil
	})

	if err != nil {
		logger.GetLogger().Printf("Error walking directory: %s", err)
	}
}
