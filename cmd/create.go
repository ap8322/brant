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
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/ap8322/brant/config"
	"github.com/ap8322/brant/ticket"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create ticket ID prefix branch",
	Long:  `Create ticket ID prefix branch (default filtering tool: fzf)`,
	RunE:  create,
}

func create(cmd *cobra.Command, args []string) error {
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

	fmt.Print(line)

	selectedId := strings.Split(line, sep)[0]

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(color.GreenString("branch name? > "))
	branchName, _ := reader.ReadString('\n')

	command := "git checkout -b " + selectedId + "_" + strings.TrimSpace(branchName)

	if len(args) != 0 {
		baseBranch := strings.TrimSpace(args[0])
		command = command + " " + baseBranch
	}

	if err := run(command, os.Stdin, os.Stdout); err != nil {
		fmt.Println("Fatal: Git or anything error. Please read above error message.")
		return nil
	}

	return nil
}

func init() {
	RootCmd.AddCommand(createCmd)
}
