package helper

import "runtime/debug"

var GitCommit = func() string {
	// Build with go build -buildvcs
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return ""
}()
