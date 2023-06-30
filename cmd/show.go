/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/gdssouza/goji/pkg/jiraclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows the details of a specified ticket",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("show called with arg %v\n", args)
		if len(args) < 1 {
			fmt.Println("Show command called with no arguments")
			os.Exit(1)
		}

		username := viper.GetViper().Sub("jira").GetString("username")
		token := viper.GetViper().Sub("jira").GetString("token")
		url := viper.GetViper().Sub("jira").GetString("url")

		jiraclient.Init(username, token, url)

		issue, err := jiraclient.GetIssue(args[0])
		if err != nil {
			panic(err)
		}

		jiraclient.PrintDetails(issue)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
