package version

import (
	"strings"
	"testing"
)

func TestVersionDetail(t *testing.T) {
	result := VersionDetail("testapp")
	if !strings.Contains(result, "testapp") {
		t.Errorf("VersionDetail should contain app name: %s", result)
	}
	if !strings.Contains(result, "Version:") {
		t.Errorf("VersionDetail should contain Version: %s", result)
	}
}
