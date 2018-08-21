// Copyright © 2018 ThreeComma.io <hello@threecomma.io>

package cmd

import (
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/threecommaio/audit/pkg"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "audit",
	Short: "ThreeComma Audit Toolkit",
	Long: `ThreeComma Audit Toolkit

This toolkit helps collect various sets of information from a host. This includes various components from the linux kernel (cpuinfo, sysctl, proc), running processes, versions of software, software configuration, etc.
It's designed to produce a human readable audit file for further analysis.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		uid := os.Getuid()
		console, _ := cmd.Flags().GetBool("console")
		upload, _ := cmd.Flags().GetString("upload")

		if uid != 0 {
			log.Fatal("this script must be run as the user [root]")
		}

		audit.Create(console, upload)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("console", "c", false, "print to console instead of file")
	rootCmd.Flags().StringP("upload", "u", "", "Accepts a Client-Token to upload to google cloud storage")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".tmp" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tmp")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
