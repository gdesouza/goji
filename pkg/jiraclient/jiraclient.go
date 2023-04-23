package jiraclient

import (
	"log"
	"os"

	"github.com/jedib0t/go-pretty/table"
	jira "gopkg.in/andygrunwald/go-jira.v1"
)

var client *jira.Client
var pageSize int
var last int
var filter string

func Print(issues []jira.Issue) {
	tbl := table.NewWriter()
	tbl.SetOutputMirror(os.Stdout)
	tbl.AppendHeader(table.Row{"Key", "Description", "Status"})

	tbl.SetStyle(table.StyleLight)
	tbl.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:     "Key",
			WidthMin: 10,
			WidthMax: 10,
		},
		{
			Name:     "Description",
			WidthMin: 120,
			WidthMax: 120,
		},
		{
			Name:     "Status",
			WidthMin: 20,
			WidthMax: 20,
		},
	})

	for _, e := range issues {
		tbl.AppendRow([]interface{}{e.Key, e.Fields.Summary, e.Fields.Status.Name})
	}
	tbl.Render()

}

func Init(username string, token string, url string) {
	// Create a BasicAuth Transport object
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: token,
	}

	// Create a new Jira Client
	var err error
	client, err = jira.NewClient(tp.Client(), url)
	if err != nil {
		log.Fatal(err)
	}
}

func GetPage(jql string, startAt int, maxResults int) ([]jira.Issue, error) {
	var issues []jira.Issue
	opt := &jira.SearchOptions{
		MaxResults: maxResults,
		StartAt:    startAt,
	}

	issues, resp, err := client.Issue.Search(jql, opt)
	if err != nil {
		return nil, err
	}

	// update last and page size for subsequent NextPage calls
	filter = jql
	last = resp.StartAt + len(issues)
	pageSize = maxResults

	return issues, nil

}

func FirstPage(jql string, size int) ([]jira.Issue, error) {
	return GetPage(jql, 0, size)
}

func NextPage() ([]jira.Issue, error) {
	return GetPage(filter, last, pageSize)
}

func All(jql string) ([]jira.Issue, error) {
	var issues []jira.Issue
	for {
		opt := &jira.SearchOptions{
			MaxResults: pageSize, // Max results can go up to 1000
			StartAt:    last,
		}

		chunk, resp, err := client.Issue.Search(jql, opt)
		if err != nil {
			return nil, err
		}

		total := resp.Total
		if issues == nil {
			issues = make([]jira.Issue, 0, total)
		}
		issues = append(issues, chunk...)
		last = resp.StartAt + len(chunk)
		if last >= total {
			return issues, nil
		}
	}
}
