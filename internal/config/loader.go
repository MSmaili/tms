package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Loader interface {
	Load() (*Config, error)
}

type FileLoader struct {
	Path string
}

func NewFileLoader(path string) *FileLoader {
	return &FileLoader{Path: path}
}

func (l *FileLoader) Load() (*Config, error) {
	data, err := os.ReadFile(l.Path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var raw Config
	if err = json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err = l.validate(&raw); err != nil {
		return nil, err
	}

	normalized, err := l.normalize(&raw)
	if err != nil {
		return nil, err
	}

	return normalized, nil

}

func (l *FileLoader) validate(cfg *Config) error {
	if cfg.Sessions == nil {
		return fmt.Errorf("sessions block missing")
	}

	for name, windows := range cfg.Sessions {
		if name == "" {
			return fmt.Errorf("session name cannot be empty")
		}
		if len(windows) == 0 {
			return fmt.Errorf("session '%s' has no windows", name)
		}
		for _, w := range windows {
			if w.Path == "" {
				return fmt.Errorf("window in session '%s' missing path", name)
			}
		}
	}

	return nil
}

func (l *FileLoader) normalize(cfg *Config) (*Config, error) {
	out := &Config{Sessions: map[string]WindowList{}}

	for name, windows := range cfg.Sessions {
		normalized := make(WindowList, len(windows))
		for i, w := range windows {
			w.Path = expandPath(w.Path)
			if w.Name == "" {
				w.Name = inferNameFromPath(w.Path)
			}
			normalized[i] = w
		}
		out.Sessions[name] = normalized
	}

	return out, nil
}

func expandPath(p string) string {
	p = os.ExpandEnv(p)

	if strings.HasPrefix(p, "~") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, strings.TrimPrefix(p, "~"))
	}
	return p
}
