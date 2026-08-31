package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

func loadSettings(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	settings := make(map[string]any)
	err = yaml.Unmarshal(data, &settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func checkPath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("the path[%s] is a directory", path)
	}

	fileExt := filepath.Ext(path)
	if len(fileExt) == 0 {
		return fmt.Errorf("file[%s] ext not found", path)
	}

	return nil
}

func Load[T any](paths ...string) (*T, error) {
	for _, p := range paths {
		err := checkPath(p)
		if err != nil {
			return nil, err
		}
	}

	merged := make(map[string]any)
	for _, p := range paths {
		settings, err := loadSettings(p)
		if err != nil {
			return nil, err
		}

		for k, v := range settings {
			merged[k] = v
		}
	}

	cfg := new(T)
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: cfg,
	})
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(merged)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
