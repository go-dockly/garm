package dbg

import "runtime/debug"

// read commit hash from build info
func init() {
	info, _ := debug.ReadBuildInfo()
	if info != nil {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				commit = setting.Value
			}
		}
	}
}

// name of the compiled binary. Populated through LD flags
var name string

func Name() string {
	if name == "" {
		return "garm"
	}

	return name
}

// version contains the version of the compiled binaries. Populated through LD flags
var version string

// Version returns the version of the binary. <unknown> is returned if version is not set
func Version() string {
	if version == "" {
		return "v1"
	}

	return version
}

var commit string

// commit hash of the binary. <unknown> is returned if commit hash is not set.
func Commit() string {
	if commit == "" {
		return "<unknown>"
	}

	return commit
}
