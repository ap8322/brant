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
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/ap8322/brant/config"
	"github.com/ap8322/brant/ticket"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print ticket information linked to the current branch",
	Long:  `Print ticket information linked to the current branch`,
	RunE:  info,
}

func info(cmd *cobra.Command, args []string) error {
	var tickets ticket.Tickets
	command := "git symbolic-ref --short HEAD"

	if err := tickets.Load(); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := run(command, os.Stdin, &buf); err != nil {
		return err
	}

	branch := buf.String()

	var target *ticket.Ticket

	for _, t := range tickets.Tickets {
		if strings.Contains(branch, t.ID) {
			target = &t
			break
		}
	}

	if target != nil {
		link := config.Conf.Jira.Host + "/browse/" + target.ID
		fmt.Println(target.Title)
		fmt.Println(link)
	} else {
		fmt.Println("no ticket linked to the current branch.")
	}

	return nil
}

func init() {
	RootCmd.AddCommand(infoCmd)
}
