package wordlist

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "words.json")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadFile(t *testing.T) {
	path := writeTemp(t, `{"words": ["shit", "son of a bitch"]}`)

	got, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"shit", "son of a bitch"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestLoadFileMissing(t *testing.T) {
	if _, err := LoadFile("does-not-exist.json"); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadFileBadJSON(t *testing.T) {
	path := writeTemp(t, `{"words": [`)
	if _, err := LoadFile(path); err == nil {
		t.Error("expected error for bad JSON")
	}
}
