package testdata

import (
	"os"
	"path/filepath"
)

// GoldenFile represents a golden file for testing
type GoldenFile struct {
	Dir string
}

// NewGoldenFile creates a new GoldenFile handler
func NewGoldenFile(dir string) *GoldenFile {
	return &GoldenFile{Dir: dir}
}

// GetPath returns the full path to a golden file
func (g *GoldenFile) GetPath(name string) string {
	return filepath.Join(g.Dir, name)
}

// Load loads a golden file's content
func (g *GoldenFile) Load(name string) ([]byte, error) {
	return os.ReadFile(g.GetPath(name))
}

// LoadString loads a golden file as string
func (g *GoldenFile) LoadString(name string) (string, error) {
	data, err := g.Load(name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Save saves content to a golden file
func (g *GoldenFile) Save(name string, content []byte) error {
	path := g.GetPath(name)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}

// SaveString saves string content to a golden file
func (g *GoldenFile) SaveString(name, content string) error {
	return g.Save(name, []byte(content))
}

// Exists checks if a golden file exists
func (g *GoldenFile) Exists(name string) bool {
	_, err := os.Stat(g.GetPath(name))
	return err == nil
}

// Update updates a golden file with new content
func (g *GoldenFile) Update(name string, content []byte) error {
	return g.Save(name, content)
}

// UpdateString updates a golden file with new string content
func (g *GoldenFile) UpdateString(name, content string) error {
	return g.Save(name, []byte(content))
}

// GoldenFileTester is a helper for golden file testing
type GoldenFileTester struct {
	GoldenDir string
	Update    bool
}

// NewGoldenFileTester creates a new GoldenFileTester
func NewGoldenFileTester(goldenDir string) *GoldenFileTester {
	return &GoldenFileTester{
		GoldenDir: goldenDir,
		Update:    false,
	}
}

// SetUpdate sets the update mode
func (t *GoldenFileTester) SetUpdate(update bool) {
	t.Update = update
}

// GetExpected returns the expected content from golden file
func (t *GoldenFileTester) GetExpected(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(t.GoldenDir, name))
}

// GetActual returns the actual content
func (t *GoldenFileTester) GetActual(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(t.GoldenDir, name))
}

// Compare compares golden file content with actual
func (t *GoldenFileTester) Compare(name string, actual []byte) (bool, error) {
	expected, err := t.GetExpected(name)
	if err != nil {
		return false, err
	}
	return string(expected) == string(actual), nil
}

// UpdateGolden updates the golden file with actual content
func (t *GoldenFileTester) UpdateGolden(name string, actual []byte) error {
	return os.WriteFile(filepath.Join(t.GoldenDir, name), actual, 0644)
}
