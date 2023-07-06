package version

import (
	"encoding/json"
)

var (
	APP_NAME string

	// Version release version of the provider
	Version string

	// GitCommitSHA is the Git SHA of the latest tag/release
	Commit string

	// GitBranch
	Branch string

	BuildTime string

	// DevVersion string for the development version
	DevVersion = "dev"
)

type VersionInfo struct {
	AppName   string
	Version   string
	Branch    string
	Commit    string
	BuildTime string
}

func NewVersionInfo() *VersionInfo {
	return &VersionInfo{
		AppName:   APP_NAME,
		Version:   BuildVersion(),
		Branch:    Branch,
		Commit:    Commit,
		BuildTime: BuildTime,
	}
}

// BuildVersion returns current version of the provider
func BuildVersion() string {
	if len(Version) == 0 {
		return DevVersion
	}
	return Version
}

func String() string {
	b, _ := json.Marshal(NewVersionInfo())
	return string(b)
}

func StringPretty() string {
	b, _ := json.MarshalIndent(NewVersionInfo(), "", "    ")
	return string(b)
}
