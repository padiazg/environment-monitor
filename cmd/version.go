/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/padiazg/environment-monitor-daemon/internals/models/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows build info and version",
	Long:  `Shows a splash screen with build and version info`,
	Run: func(cmd *cobra.Command, args []string) {
		version.Splash()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
