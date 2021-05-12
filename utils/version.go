package utils

import "fmt"

var (
	BuildVersion string
	BuildTime    string
	BuildName    string
)

func Version() string {
	return fmt.Sprintf("%s: v%s (%s)\n", BuildName, BuildVersion, BuildTime)
}
