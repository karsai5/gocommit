/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/cqroot/prompt"
	"github.com/karsai5/gocommit/cmd/git"
	"github.com/karsai5/gocommit/cmd/message"
	"github.com/spf13/cobra"
)

var noVerifyFlag bool

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		numStagedFiles, err := git.NumberOfStagedFiles()
		if err != nil {
			panic(err)
		}

		if numStagedFiles == 0 {
			fmt.Println("No staged files to commit.")
			os.Exit(1)
		}

		if !noVerifyFlag {
			output, err := git.RunPreCommitHook()
			println(output)
			if err != nil {
				panic(err)
			}
		}

		ticketNumber, err := git.TicketNumberFromBranchName()
		if err != nil {
			panic(err)
		}

		cm, err := message.NewCommitMessage(
			message.WithTicket(ticketNumber),
			message.WithType(promptForCommitType()),
		)
		if err != nil {
			panic(err)
		}

		oneLineDescription := promptForCommitMessage(cm.Message())
		err = cm.ApplyOption(message.WithOneLineDescription(oneLineDescription))
		if err != nil {
			panic(err)
		}

		if err = cm.Valid(); err != nil {
			panic(err)
		}

		commit, err := git.NewCommit(cm.Message(), git.WithNoVerify())
		if err != nil {
			panic(err)
		}

		cmdOutput, err := commit.Cmd().CombinedOutput()
		println(string(cmdOutput))
		if err != nil {
			panic(err)
		}
	},
}

func promptForCommitMessage(s string) string {
	val, err := prompt.New().Ask(s).Input("One line message...")
	CheckErr(err)
	return val
}

func promptForCommitType() string {
	noneType := "None"
	options := append([]string{"None"}, CommitTypes...)
	val1, err := prompt.New().Ask("Choose:").Choose(options)
	CheckErr(err)
	if val1 == noneType {
		return ""
	}
	return val1
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.PersistentFlags().BoolVarP(&noVerifyFlag, "no-verify", "n", false, "Skip pre-commit hook")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func CheckErr(err error) {
	if err != nil {
		if errors.Is(err, prompt.ErrUserQuit) {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		} else {
			panic(err)
		}
	}
}
