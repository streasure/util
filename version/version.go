package version

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	Version   string
	GoVersion string
	Built     string
	GitCommit string
	OSArch    string
)

func VersionDetail(appName string) string {
	if len(strings.TrimSpace(Version)) == 0 {
		Version = "0.0.1"
	}
	if len(strings.TrimSpace(GoVersion)) == 0 {
		GoVersion = runtime.Version()
	}
	if len(strings.TrimSpace(GitCommit)) == 0 {
		GitCommit = "unknown"
	}
	if len(strings.TrimSpace(Built)) == 0 {
		Built = "0000-00-00 00:00:00"
	}
	if len(strings.TrimSpace(OSArch)) == 0 {
		OSArch = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	}
	return fmt.Sprintf("<AppName:%s Version:%s GoVersion:%s Built:%s GitCommit:%s OSArch:%s>", appName, Version, GoVersion, Built, GitCommit, OSArch)
}
