/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/gdssouza/goji/pkg/jiraclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Jira issues",
	Long:  `Use this to list Jira issues. Use the -f (filter) flag to specify a filter for your query.`,
	Run: func(cmd *cobra.Command, args []string) {
		username := viper.GetViper().Sub("jira").GetString("username")
		token := viper.GetViper().Sub("jira").GetString("token")
		url := viper.GetViper().Sub("jira").GetString("url")
		maxResults := viper.GetViper().Sub("jira").GetInt("maxResults")

		jiraclient.Init(username, token, url)

		issues, err := jiraclient.FirstPage(filterVar, maxResults)
		if err != nil {
			panic(err)
		}

		jiraclient.Print(issues)

		for len(issues) > 0 {
			fmt.Println("--- press <ENTER> to continue ---")
			fmt.Scanln()

			issues, err = jiraclient.NextPage()
			if err != nil {
				panic(err)
			}

			jiraclient.Print(issues)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Command specific flags
	listCmd.PersistentFlags().StringVarP(&filterVar, "filter", "f", "*", "Filter string for the query")
}
