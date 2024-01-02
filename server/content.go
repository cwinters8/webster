package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/google/uuid"
)

type Content struct {
	Text string `form:"editor"`
	Path string
}

func (c *Content) WriteTemp() error {
	if len(c.Text) == 0 {
		return fmt.Errorf("c.Text is an empty string")
	}
	tmpl, err := template.New("t").Parse(c.Text)
	if err != nil {
		return fmt.Errorf("content failed to parse as HTML: %w", err)
	}
	dir, err := os.MkdirTemp("tmp", "")
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	id := uuid.New()
	c.Path = fmt.Sprintf("%s/%s.html", dir, id.String())
	f, err := os.Create(c.Path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	// execute template and write to disk
	if err := tmpl.Execute(f, nil); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}
