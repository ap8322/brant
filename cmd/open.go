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
	"strings"
	"github.com/ap8322/brant/config"
	"github.com/ap8322/brant/ticket"
	"github.com/spf13/cobra"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open ticket url you selected",
	Long: `Open ticket url you selected`,
	RunE: open,
}

func open(cmd *cobra.Command, args []string) error {
	// TODO refactoring
	var tickets ticket.Tickets
	sep := " : "

	if err := tickets.Load(); err != nil {
		return err
	}

	var list string

	for _, t := range tickets.Tickets {
		list += t.ID + sep + t.Title + "\n"
	}

	var buf bytes.Buffer
	if err := run(config.Conf.Core.SelectCmd, strings.NewReader(list), &buf); err != nil {
		return nil
	}

	line := buf.String()

	var target *ticket.Ticket

	for _, t := range tickets.Tickets {
		if strings.Contains(line, t.ID) {
			target = &t
			break
		}
	}

	err := openbrowser(config.Conf.Jira.Host + "/browse/" + target.ID)

	if err != nil {
		return err
	}

	return nil

}

func init() {
	RootCmd.AddCommand(openCmd)
}
