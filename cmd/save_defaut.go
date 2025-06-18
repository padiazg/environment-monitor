/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/padiazg/environment-monitor-daemon/internals/models/settings"
	"github.com/spf13/cobra"
)

// saveDefautCmd represents the saveDefaut command
var saveDefaultCmd = &cobra.Command{
	Use:   "save-default",
	Short: "Save an example configuration file",
	Long:  `Save a configuration file with example values`,
	Run: func(cmd *cobra.Command, args []string) {
		file_name, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Printf("getting file flag value: %+v", err)
			return
		}

		ex := &settings.Settings{}
		(*ex) = (*s) // lazzy deep copy

		// Fill some example values
		ex.Url = "http://api/path"
		ex.Sensor = "dummy"
		ex.ApiKey = "abcdef0123456789"
		ex.Source = "MySensor/1"
		ex.Description = "My first sensor"
		ex.Lat = -25.376054
		ex.Lon = -57.545052
		ex.Interval = 5 * time.Second
		ex.ZH07 = &settings.ZH07Settings{
			Mode: "qa",
			Port: "/dev/serial0",
		}
		ex.Dummy = &settings.DummySettings{
			PM1:  1.0,
			PM25: 2.5,
			PM10: 10.0,
		}

		ex.Save(file_name)
	},
}

func init() {
	rootCmd.AddCommand(saveDefaultCmd)
	saveDefaultCmd.Flags().StringP("file", "f", ".environment-monitor.default.yaml", "Save a sampe for the configuration file")
}
