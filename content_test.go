package webster_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/cwinters8/webster"
)

func TestContent(t *testing.T) {
	dir := "testdata/hello"
	c := webster.Content{
		HTML: "<h1>Hello, World!</h1>",
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
	if got != c.HTML {
		t.Errorf("wanted text `%s`; got `%s`", c.HTML, got)
	} else {
		t.Log(got)
	}

	t.Run("load_saved", func(t *testing.T) {
		content, err := webster.Load(c.Path)
		if err != nil {
			t.Fatalf("failed to load: %v", err)
		}
		if content.HTML != c.HTML {
			t.Errorf("wanted text `%s`; got `%s`", c.HTML, content.HTML)
		}
	})

	t.Run("load_nonexistent", func(t *testing.T) {
		file := "testdata/fake/index.html"
		remove := func(t *testing.T) {
			if err := os.RemoveAll(file); err != nil {
				t.Fatalf("failed to remove file: %v", err)
			}
		}
		// ensure file doesn't already exist
		remove(t)
		content, err := webster.Load(file)
		if err != nil {
			t.Fatalf("failed to load content: %v", err)
		}
		// clean up
		defer remove(t)

		if content.Path != file {
			t.Errorf("wanted file path `%s`; got `%s`", file, content.Path)
		}
		if _, err := os.Stat(content.Path); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				t.Errorf("file not created")
			} else {
				t.Fatalf("failed to stat file: %v", err)
			}
		}
		if len(content.HTML) > 0 {
			t.Errorf("wanted content text to be empty; got `%s`", content.HTML)
		}

		// validate file can be written to
		text := "<p>Some fake data</p>"
		content.HTML = text
		if err := content.Save(); err != nil {
			t.Fatalf("failed to save new content: %v", err)
		}
		out, err := os.ReadFile(content.Path)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}
		got := string(out)
		if got != text {
			t.Errorf("wanted text `%s`; got `%s`", text, got)
		}
	})
}
