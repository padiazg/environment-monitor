package help

import (
	"fmt"
	"os"
	"path/filepath"
)

func SettingsHelp() {
	// appName, _ := os.Executable()
	fmt.Printf(`1. Run "%s save-default"
2. Edit the ".environment-monitor.default.yaml" file according to your needs
3. Rename the file to ".environment-monitor.yaml" so the app could find it, or
use the "-f .environment-monitor.default.yaml" flag to tell the app to use it
`, filepath.Base(os.Args[0]))
}
