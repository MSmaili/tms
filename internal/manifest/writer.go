package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func Write(workspace *Workspace, path string) error {
	extendedPath := expandPath(path)

	if err := os.MkdirAll(filepath.Dir(extendedPath), 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	var data []byte
	var err error
	ext := filepath.Ext(extendedPath)

	switch ext {
	case ".yaml", ".yml":
		data, err = yaml.Marshal(workspace)
		if err != nil {
			return fmt.Errorf("marshal yaml: %w", err)
		}
	case ".json":
		data, err = json.MarshalIndent(workspace, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal json: %w", err)
		}
	default:
		return fmt.Errorf("unsupported format: %s (use .yaml, .yml, or .json)", ext)
	}

	if err := os.WriteFile(extendedPath, data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
