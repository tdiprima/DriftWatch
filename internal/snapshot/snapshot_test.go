package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileAtomicReplacesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, baselineFileName)

	if err := os.WriteFile(path, []byte("old baseline"), 0600); err != nil {
		t.Fatalf("writing original file: %v", err)
	}

	if err := writeFileAtomic(path, []byte("new baseline"), 0600); err != nil {
		t.Fatalf("writeFileAtomic() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading replaced file: %v", err)
	}
	if string(got) != "new baseline" {
		t.Fatalf("replaced file = %q, want %q", got, "new baseline")
	}

	matches, err := filepath.Glob(filepath.Join(dir, baselineFileName+".*.tmp"))
	if err != nil {
		t.Fatalf("globbing temp files: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("temporary files left behind: %v", matches)
	}
}
