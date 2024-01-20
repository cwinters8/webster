package webster

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
)

type Content struct {
	ID       uuid.UUID `db:"id"`
	HTML     string    `db:"html"`
	Path     string    `db:"path"`
	Versions map[uuid.UUID]Version
	tempFile *os.File
}

type Version struct {
	ID        uuid.UUID `db:"id"`
	Timestamp time.Time `db:"timestamp"`
}

var (
	ErrEmptyText = errors.New("content is an empty string")
	ErrParseHTML = errors.New("content failed to parse as HTML")
)

// file must be a qualified path
//
// if file does not exist, a new one is created
func Load(file string) (*Content, error) {
	if len(file) == 0 {
		return nil, fmt.Errorf("file must not be empty")
	}
	contents, err := os.ReadFile(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			dir := path.Dir(file)
			if dir != "." {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return nil, fmt.Errorf("failed to create directory `%s`: %w", dir, err)
				}
			}
			if _, err := os.Create(file); err != nil {
				return nil, fmt.Errorf("failed to create file `%s`: %w", file, err)
			}
			return &Content{Path: file}, nil
		}
		return nil, fmt.Errorf("failed to read file `%s`: %w", file, err)
	}
	return &Content{
		HTML: string(contents),
		Path: file,
	}, nil
}

// may not be needed
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
	if len(c.HTML) == 0 {
		return nil, ErrEmptyText
	}
	tmpl, err := template.New("t").Parse(c.HTML)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrParseHTML, err)
	}
	return tmpl, nil
}
