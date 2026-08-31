package config

import (
	"os"
	"path/filepath"
	"testing"
)

type testConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	err := os.WriteFile(path, []byte("host: localhost\nport: 8080\nname: test\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := Load[testConfig](path)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	if cfg.Host != "localhost" {
		t.Fatalf("Host: got %q, want %q", cfg.Host, "localhost")
	}
	if cfg.Port != 8080 {
		t.Fatalf("Port: got %d, want 8080", cfg.Port)
	}
	if cfg.Name != "test" {
		t.Fatalf("Name: got %q, want %q", cfg.Name, "test")
	}
}

func TestLoadMultipleFiles(t *testing.T) {
	dir := t.TempDir()

	path1 := filepath.Join(dir, "config1.yaml")
	err := os.WriteFile(path1, []byte("host: host1\nport: 1111\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	path2 := filepath.Join(dir, "config2.yaml")
	err = os.WriteFile(path2, []byte("name: from_file2\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := Load[testConfig](path1, path2)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	if cfg.Host != "host1" {
		t.Fatalf("Host: got %q, want %q", cfg.Host, "host1")
	}
	if cfg.Port != 1111 {
		t.Fatalf("Port: got %d, want 1111", cfg.Port)
	}
	if cfg.Name != "from_file2" {
		t.Fatalf("Name: got %q, want %q", cfg.Name, "from_file2")
	}
}

func TestLoadFileNotFound(t *testing.T) {
	_, err := Load[testConfig]("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestLoadDirectory(t *testing.T) {
	dir := t.TempDir()
	_, err := Load[testConfig](dir)
	if err == nil {
		t.Fatal("expected error for directory path")
	}
}

func TestLoadNoExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config")
	err := os.WriteFile(path, []byte("host: localhost\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Load[testConfig](path)
	if err == nil {
		t.Fatal("expected error for file without extension")
	}
}
