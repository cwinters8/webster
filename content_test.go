package craftly_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cwinters8/craftly"
)

func TestSave(t *testing.T) {
	dir := "testdata/hello"
	c := craftly.Content{
		Text: "<h1>Hello, World!</h1>",
		Path: fmt.Sprintf("%s/index.html", dir),
	}
	if err := c.Save(); err != nil {
		t.Fatalf("failed to save: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatalf("failed to clean up directory %s: %v", dir, err)
		}
	}()
	out, err := os.ReadFile(c.Path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", c.Path, err)
	}
	got := string(out)
	if got != c.Text {
		t.Errorf("wanted text `%s`; got `%s`", c.Text, got)
	} else {
		t.Log(got)
	}
}
