/*
Copyright © 2025 Pato Diaz (padiazg@gmail.com)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	// "log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/padiazg/environment-monitor-daemon/internals/models/settings"
	"github.com/padiazg/environment-monitor-daemon/internals/monitor"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultInterval = 120 * time.Second
)

var (
	cfgFile string
	s       = &settings.Settings{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "environment-monitor-daemon",
	Short: "An air-quality monitor daemon",
	Long: `An air-quality monitoring application intended to be used to colaborate 
with the [Aire Libre](http://airelib.re/) project.

This application is meant to run in any SBC that has serial port headers or implemente I2c`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// s.Show()

		if err := s.Validate(); err != nil {
			log.Fatalf("Validating settings: %v", err)
		}

		monitor.Start(s)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() { initConfig(".environment-monitor", "em") })

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.environment-monitor.yaml)")

	rootCmd.Flags().StringP("url", "u", "", "Request URL")
	viper.BindPFlag("url", rootCmd.Flag("url"))

	rootCmd.Flags().StringP("apikey", "a", "", "API key")
	viper.BindPFlag("apikey", rootCmd.Flag("apikey"))

	rootCmd.Flags().StringP("source", "s", "", "Name used to identify the monitor")
	viper.BindPFlag("source", rootCmd.Flag("source"))

	rootCmd.Flags().StringP("description", "d", "", "User friendly name to identify the monitor")
	viper.BindPFlag("description", rootCmd.Flag("description"))

	rootCmd.Flags().StringP("sensor", "e", "", "Sensor model used for readings [SPS30|ZH07|PMS007]")
	viper.BindPFlag("sensor", rootCmd.Flag("sensor"))

	rootCmd.Flags().Float64("lat", 0.0, "Latitude")
	viper.BindPFlag("lat", rootCmd.Flag("lat"))

	rootCmd.Flags().Float64("lon", 0.0, "Longitude")
	viper.BindPFlag("lon", rootCmd.Flag("lon"))

	rootCmd.Flags().DurationP("interval", "i", defaultInterval, "Interval between readings")
	viper.BindPFlag("interval", rootCmd.Flag("interval"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig(configName, envPrefix string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name defined by configName (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(configName)
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix(envPrefix)

	setDefaults()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	bindEnvs(s)

	if err := viper.Unmarshal(s); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	// s.Show()

	// if err := s.Validate(); err != nil {
	// 	log.Fatalf("Validating settings: %v", err)
	// }
}

func setDefaults() {
	viper.SetDefault("interval", defaultInterval)
	// viper.SetDefault("zh07.mode", "qa")
	// viper.SetDefault("zh07.port", "/dev/serial0")
}

// bindEnvs creates the environment variable bindings for the given struct, also aliases for proper
// binding of environment variables and values from .env files and other structured config files.
func bindEnvs(i interface{}, parts ...string) {
	ifv := reflect.ValueOf(i)
	ift := reflect.TypeOf(i)

	// received a pointer, dereference it
	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
		ift = ift.Elem()
	}

	for x := 0; x < ift.NumField(); x++ {
		t := ift.Field(x)
		v := ifv.Field(x)

		if !t.IsExported() {
			// fmt.Printf("  not exported: %s\n", t.Name)
			continue
		}

		// fmt.Printf("field: %s, value: %s\n", t.Name, v.String())
		switch v.Kind() {
		case reflect.Struct:
			// fmt.Printf("  struct: %s\n", t.Name)
			bindEnvs(v.Interface(), append(parts, t.Name)...)

		case reflect.Ptr:
			if v.IsNil() {
				// fmt.Printf("  nil pointer: %s\n", t.Name)
				continue
			}
			// fmt.Printf("  pointer: %s\n", t.Name)
			bindEnvs(v.Interface(), append(parts, t.Name)...)

		default:
			var (
				envKey   = strings.ToUpper(strings.Join(append(parts, t.Name), "_"))
				key      = strings.Join(append(parts, t.Name), ".")
				envAlias = strings.ToLower(envKey)
			)

			// set the env binding
			if err := viper.BindEnv(key, envKey); err != nil {
				log.Fatalf("config: unable to bind env: %s", err.Error())
			}

			viper.RegisterAlias(envAlias, key)

			// fmt.Printf("  key: %s => %s => %s\n", key, envKey, envAlias)
		}
	}
}
