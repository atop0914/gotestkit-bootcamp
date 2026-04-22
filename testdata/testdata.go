// Package testdata provides test data loading utilities.
package testdata

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
)

// Load reads a file from the testdata directory
func Load(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join("testdata", name))
}

// LoadJSON loads and parses a JSON file into v
func LoadJSON(name string, v interface{}) error {
	data, err := Load(name)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// LoadXML loads and parses an XML file into v
func LoadXML(name string, v interface{}) error {
	data, err := Load(name)
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}

// MustLoad loads a file or panics
func MustLoad(name string) []byte {
	data, err := Load(name)
	if err != nil {
		panic("testdata: " + err.Error())
	}
	return data
}

// MustLoadJSON loads JSON or panics
func MustLoadJSON(name string, v interface{}) {
	if err := LoadJSON(name, v); err != nil {
		panic("testdata: " + err.Error())
	}
}

// ReadFile reads a file relative to the testdata directory
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// Paths returns paths of files in testdata matching pattern
func Paths(pattern string) ([]string, error) {
	return filepath.Glob(filepath.Join("testdata", pattern))
}

// FileContent returns the content of a file as a string
func FileContent(name string) string {
	data, err := Load(name)
	if err != nil {
		return ""
	}
	return string(data)
}

// NormalizeJSON normalizes JSON for comparison
func NormalizeJSON(s string) string {
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return s
	}
	b, _ := json.Marshal(v)
	return string(b)
}

// NormalizeXML normalizes XML for comparison
func NormalizeXML(s string) string {
	var v interface{}
	if err := xml.Unmarshal([]byte(s), &v); err != nil {
		return s
	}
	b, _ := json.Marshal(v)
	return string(b)
}

// CompareFiles compares two files content
func CompareFiles(name1, name2 string) (bool, error) {
	d1, err := Load(name1)
	if err != nil {
		return false, err
	}
	d2, err := Load(name2)
	if err != nil {
		return false, err
	}
	return string(d1) == string(d2), nil
}

// AppendLine appends a line to a file in testdata
func AppendLine(name, line string) error {
	f, err := os.OpenFile(filepath.Join("testdata", name),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(line + "\n")
	return err
}

// FindFiles finds files with extension in testdata
func FindFiles(ext string) ([]string, error) {
	pattern := filepath.Join("testdata", "*"+ext)
	return filepath.Glob(pattern)
}

// IsEmpty checks if testdata file is empty
func IsEmpty(name string) bool {
	data, err := Load(name)
	if err != nil {
		return true
	}
	return len(strings.TrimSpace(string(data))) == 0
}
