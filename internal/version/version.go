// Package version is a singleton module which stores project build information.
package version

import (
	"fmt"
)

type buildInfo struct {
	buildNumber string
	buildDate   string
	buildCommit string
}

var bi buildInfo

func init() {
	valueIsNotAvailable := "N/A"

	bi = buildInfo{
		buildNumber: valueIsNotAvailable,
		buildDate:   valueIsNotAvailable,
		buildCommit: valueIsNotAvailable,
	}
}

// Set should be called from the main function to make application version details available for other app modules
func Set(buildVersion, buildDate, buildCommit string) {
	if len(buildVersion) > 0 {
		bi.buildNumber = buildVersion
	}

	if len(buildDate) > 0 {
		bi.buildDate = buildDate
	}

	if len(buildCommit) > 0 {
		bi.buildCommit = buildCommit
	}
}

// BuildVersion sets version of the application
func BuildVersion() string {
	return bi.buildNumber
}

// BuildDate sets date of the build
func BuildDate() string {
	return bi.buildDate
}

// BuildCommit sets last commit id
func BuildCommit() string {
	return bi.buildCommit
}

// PrintConsole prints application version details including build information in the terminal
func PrintConsole() {
	fmt.Printf("Build version: %s\n", BuildVersion())
	fmt.Printf("Build date: %s\n", BuildDate())
	fmt.Printf("Build commit: %s\n", BuildCommit())
}
