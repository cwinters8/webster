package webster

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
)

type Content struct {
	Text     string `form:"editor"`
	Path     string
	tempFile *os.File
}

var (
	ErrEmptyText = errors.New("content is an empty string")
	ErrParseHTML = errors.New("content failed to parse as HTML")
)

func (c *Content) WriteTemp() error {
	tmpl, err := c.Parse()
	if err != nil {
		return err
	}
	if c.tempFile == nil {
		dir, err := os.MkdirTemp("tmp", "")
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		c.tempFile, err = os.Create(fmt.Sprintf("%s/%s", dir, path.Clean(c.Path)))
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
	}
	// execute template and write to disk
	if err := tmpl.Execute(c.tempFile, nil); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}

func (c *Content) Save() error {
	tmpl, err := c.Parse()
	if err != nil {
		return err
	}
	dir := path.Dir(c.Path)
	// ensure the directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	f, err := os.OpenFile(c.Path, os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", c.Path, err)
	}
	if err := tmpl.Execute(f, nil); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}

func (c *Content) Parse() (*template.Template, error) {
	if len(c.Text) == 0 {
		return nil, ErrEmptyText
	}
	tmpl, err := template.New("t").Parse(c.Text)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrParseHTML, err)
	}
	return tmpl, nil
}
