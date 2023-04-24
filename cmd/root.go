/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

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
	"bufio"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var filterVar string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goji",
	Short: "A Jira CLI interface",
	Long:  `Used to list and view Jira tickets from the CLI.`,
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goji.yaml)")
}

// buildConfig iteractively creates a config file
func buildConfig(pathName string, fileName string) {
	fullPath := path.Join(pathName, fileName)
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter JIRA username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter JIRA token: ")
	token, _ := reader.ReadString('\n')

	fmt.Print("Enter JIRA URL: ")
	url, _ := reader.ReadString('\n')

	f, err := os.Create(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("jira:\n")
	f.WriteString("    username: " + username)
	f.WriteString("    token: " + token)
	f.WriteString("    url: " + url)
	f.WriteString("    maxResults: 100")
	f.Sync()

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	configName := ".goji"

	// Find home directory.
	home, err := os.UserHomeDir()

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		cobra.CheckErr(err)

		// Search config in home directory with name ".goji" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(configName)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Config file not found. Please run config command to create your configuration.")
		buildConfig(home, configName)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintln(os.Stderr, "Error in config file: ", err)
		}
	}

	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
}
