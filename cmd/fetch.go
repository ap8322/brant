// Copyright © 2017 Yuki Haneda
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/ap8322/brant/config"
	"github.com/ap8322/brant/ticket"
	"github.com/briandowns/spinner"
	"gopkg.in/kyokomi/emoji.v1"

	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch tiket list and cache them.",
	Long:  "fetch tiket list and cache them.",
	RunE:  fetch,
}

func fetch(cmd *cobra.Command, args []string) (err error) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = " fetch ticket list..."
	s.Start()

	jiraClient, err := jira.NewClient(nil, config.Conf.Jira.Host)
	if err != nil {
		return err
	}

	res, err := jiraClient.Authentication.AcquireSessionCookie(
		config.Conf.Jira.UserName,
		config.Conf.Jira.Password,
	)

	if err != nil || res == false {
		fmt.Printf("Result: %v\n", res)
		return err
	}

	issues, _, err := jiraClient.Issue.Search(config.Conf.Jira.Jql, nil)
	if err != nil {
		return err
	}

	s.Stop()

	var tikets ticket.Tickets

	for _, issue := range issues {
		ticket := ticket.Ticket{
			ID:    issue.Key,
			Title: issue.Fields.Summary,
		}

		tikets.Tickets = append(tikets.Tickets, ticket)
	}

	if err := tikets.Save(); err != nil {
		return nil
	}

	emoji.Println(":white_check_mark: saved ticket list!")

	return nil
}

func init() {
	RootCmd.AddCommand(fetchCmd)
}
