package testdata

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temp testdata file
	tmpDir := t.TempDir()
	testdataDir := filepath.Join(tmpDir, "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test file
	testFile := filepath.Join(testdataDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hello world"), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to tmp dir so Load() can find testdata dir
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Test Load
	data, err := Load("test.txt")
	if err != nil {
		t.Errorf("Load() error = %v", err)
		return
	}
	if string(data) != "hello world" {
		t.Errorf("Load() = %q, want %q", string(data), "hello world")
	}
}

func TestLoadJSON(t *testing.T) {
	tmpDir := t.TempDir()
	testdataDir := filepath.Join(tmpDir, "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test JSON file
	testFile := filepath.Join(testdataDir, "test.json")
	jsonContent := `{"name":"test","value":123}`
	if err := os.WriteFile(testFile, []byte(jsonContent), 0644); err != nil {
		t.Fatal(err)
	}

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	var result TestStruct
	err := LoadJSON("test.json", &result)
	if err != nil {
		t.Errorf("LoadJSON() error = %v", err)
		return
	}
	if result.Name != "test" || result.Value != 123 {
		t.Errorf("LoadJSON() = %+v, want {Name:test, Value:123}", result)
	}
}

func TestLoadXML(t *testing.T) {
	tmpDir := t.TempDir()
	testdataDir := filepath.Join(tmpDir, "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test XML file
	testFile := filepath.Join(testdataDir, "test.xml")
	xmlContent := `<test><name>test</name><value>123</value></test>`
	if err := os.WriteFile(testFile, []byte(xmlContent), 0644); err != nil {
		t.Fatal(err)
	}

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	type TestStruct struct {
		Name  string `xml:"name"`
		Value int    `xml:"value"`
	}

	var result TestStruct
	err := LoadXML("test.xml", &result)
	if err != nil {
		t.Errorf("LoadXML() error = %v", err)
		return
	}
	if result.Name != "test" || result.Value != 123 {
		t.Errorf("LoadXML() = %+v, want {Name:test, Value:123}", result)
	}
}

func TestMustLoad(t *testing.T) {
	tmpDir := t.TempDir()
	testdataDir := filepath.Join(tmpDir, "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(testdataDir, "must.txt")
	if err := os.WriteFile(testFile, []byte("must content"), 0644); err != nil {
		t.Fatal(err)
	}

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	data := MustLoad("must.txt")
	if string(data) != "must content" {
		t.Errorf("MustLoad() = %q, want %q", string(data), "must content")
	}
}

func TestPaths(t *testing.T) {
	tmpDir := t.TempDir()
	testdataDir := filepath.Join(tmpDir, "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test files
	files := []string{"file1.txt", "file2.txt", "file3.json"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(testdataDir, f), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	paths, err := Paths("*.txt")
	if err != nil {
		t.Errorf("Paths() error = %v", err)
		return
	}
	if len(paths) != 2 {
		t.Errorf("Paths() returned %d files, want 2", len(paths))
	}
}

func TestFileContent(t *testing.T) {
	tmpDir := t.TempDir()
	testdataDir := filepath.Join(tmpDir, "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(testdataDir, "content.txt")
	if err := os.WriteFile(testFile, []byte("file content here"), 0644); err != nil {
		t.Fatal(err)
	}

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	content := FileContent("content.txt")
	if content != "file content here" {
		t.Errorf("FileContent() = %q, want %q", content, "file content here")
	}
}

func TestNormalizeJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", `{"a":1}`, `{"a":1}`},
		{"with space", `{"a": 1}`, `{"a":1}`},
		{"invalid returns original", `not json`, `not json`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeJSON(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeJSON() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCompareFiles(t *testing.T) {
	tmpDir := t.TempDir()
	testdataDir := filepath.Join(tmpDir, "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Create two identical files
	file1 := filepath.Join(testdataDir, "file1.txt")
	file2 := filepath.Join(testdataDir, "file2.txt")
	content := []byte("same content")

	if err := os.WriteFile(file1, content, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, content, 0644); err != nil {
		t.Fatal(err)
	}

	same, err := CompareFiles("file1.txt", "file2.txt")
	if err != nil {
		t.Errorf("CompareFiles() error = %v", err)
		return
	}
	if !same {
		t.Errorf("CompareFiles() = false, want true for identical files")
	}

	// Modify file2
	if err := os.WriteFile(file2, []byte("different"), 0644); err != nil {
		t.Fatal(err)
	}

	diff, err := CompareFiles("file1.txt", "file2.txt")
	if err != nil {
		t.Errorf("CompareFiles() error = %v", err)
		return
	}
	if diff {
		t.Errorf("CompareFiles() = true, want false for different files")
	}
}

func TestFindFiles(t *testing.T) {
	tmpDir := t.TempDir()
	testdataDir := filepath.Join(tmpDir, "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test files
	files := map[string]string{
		"data.json":  `{}`,
		"data.xml":   `<xml/>`,
		"other.txt":  "text",
		"nested/data.json": `{"nested":true}`,
	}

	for path, content := range files {
		fullPath := filepath.Join(testdataDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	jsonFiles, err := FindFiles(".json")
	if err != nil {
		t.Errorf("FindFiles() error = %v", err)
		return
	}
	// Note: FindFiles(".json") = testdata/*.json, does not recurse into subdirs
	if len(jsonFiles) != 1 {
		t.Errorf("FindFiles(.json) returned %d files, want 1 (no recursion)", len(jsonFiles))
	}
}

func TestIsEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	testdataDir := filepath.Join(tmpDir, "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Create empty file
	emptyFile := filepath.Join(testdataDir, "empty.txt")
	if err := os.WriteFile(emptyFile, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	// Create non-empty file
	nonEmpty := filepath.Join(testdataDir, "nonempty.txt")
	if err := os.WriteFile(nonEmpty, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	if !IsEmpty("empty.txt") {
		t.Errorf("IsEmpty(empty.txt) = false, want true")
	}

	if IsEmpty("nonempty.txt") {
		t.Errorf("IsEmpty(nonempty.txt) = true, want false")
	}

	// Non-existent file should be considered empty
	if !IsEmpty("nonexistent.txt") {
		t.Errorf("IsEmpty(nonexistent.txt) = false, want true")
	}
}

// Golden file tests
func TestGoldenFile_GetPath(t *testing.T) {
	gf := NewGoldenFile("/some/dir")
	path := gf.GetPath("test.golden")
	expected := "/some/dir/test.golden"
	if path != expected {
		t.Errorf("GetPath() = %q, want %q", path, expected)
	}
}

func TestGoldenFile_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	gf := NewGoldenFile(tmpDir)

	content := []byte("golden content here")
	name := "test.golden"

	// Save
	err := gf.Save(name, content)
	if err != nil {
		t.Errorf("Save() error = %v", err)
		return
	}

	// Load
	data, err := gf.Load(name)
	if err != nil {
		t.Errorf("Load() error = %v", err)
		return
	}
	if string(data) != string(content) {
		t.Errorf("Load() = %q, want %q", string(data), string(content))
	}
}

func TestGoldenFile_Exists(t *testing.T) {
	tmpDir := t.TempDir()
	gf := NewGoldenFile(tmpDir)

	// Create a file
	if err := gf.SaveString("exists.golden", "content"); err != nil {
		t.Fatal(err)
	}

	if !gf.Exists("exists.golden") {
		t.Errorf("Exists(exists.golden) = false, want true")
	}

	if gf.Exists("nonexistent.golden") {
		t.Errorf("Exists(nonexistent.golden) = true, want false")
	}
}

func TestGoldenFile_Update(t *testing.T) {
	tmpDir := t.TempDir()
	gf := NewGoldenFile(tmpDir)

	// Create initial file
	gf.SaveString("update.golden", "original")

	// Update
	newContent := "updated content"
	if err := gf.UpdateString("update.golden", newContent); err != nil {
		t.Errorf("UpdateString() error = %v", err)
		return
	}

	// Verify
	loaded, _ := gf.LoadString("update.golden")
	if loaded != newContent {
		t.Errorf("After update, LoadString() = %q, want %q", loaded, newContent)
	}
}

func TestGoldenFileTester_Compare(t *testing.T) {
	tmpDir := t.TempDir()

	// Create golden file
	goldenContent := []byte("expected output")
	goldenPath := filepath.Join(tmpDir, "output.golden")
	if err := os.WriteFile(goldenPath, goldenContent, 0644); err != nil {
		t.Fatal(err)
	}

	tester := NewGoldenFileTester(tmpDir)

	// Test matching content
	match, err := tester.Compare("output.golden", goldenContent)
	if err != nil {
		t.Errorf("Compare() error = %v", err)
		return
	}
	if !match {
		t.Errorf("Compare() = false, want true for matching content")
	}

	// Test non-matching content
	diffContent := []byte("actual output")
	match, err = tester.Compare("output.golden", diffContent)
	if err != nil {
		t.Errorf("Compare() error = %v", err)
		return
	}
	if match {
		t.Errorf("Compare() = true, want false for non-matching content")
	}
}

func TestGoldenFileTester_UpdateGolden(t *testing.T) {
	tmpDir := t.TempDir()
	tester := NewGoldenFileTester(tmpDir)

	// Create initial golden file
	initialContent := []byte("initial")
	goldenPath := filepath.Join(tmpDir, "update.golden")
	if err := os.WriteFile(goldenPath, initialContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Update golden file
	newContent := []byte("updated golden content")
	err := tester.UpdateGolden("update.golden", newContent)
	if err != nil {
		t.Errorf("UpdateGolden() error = %v", err)
		return
	}

	// Verify
	updated, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(updated) != string(newContent) {
		t.Errorf("Golden file content = %q, want %q", string(updated), string(newContent))
	}
}
