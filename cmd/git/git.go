package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func TicketNumberFromBranchName() (string, error) {
	branchNameBytes, err := exec.Command("git", "symbolic-ref", "--short", "HEAD").Output()
	if err != nil {
		return "", err
	}

	re, err := regexp.Compile(`[mM][kK][pP]-\d\d\d\d`)
	if err != nil {
		return "", err
	}

	matched := re.FindString(string(branchNameBytes))
	if matched == "" {
		return "", nil
	}
	return strings.ToUpper(matched), nil
}

func NumberOfStagedFiles() (int, error) {
	output, err := exec.Command("git", "diff", "--cached", "--name-only").CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("failed to get git status: %s: %w", output, err)
	}
	lines := removeEmptyStrings(strings.Split(string(output), "\n"))
	return len(lines), nil
}

func removeEmptyStrings(input []string) []string {
	var result []string
	for _, str := range input {
		if strings.TrimSpace(str) != "" {
			result = append(result, str)
		}
	}
	return result
}

func RepoPath() (string, error) {
	path, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(path)), nil
}

func RunPreCommitHook() (output string, err error) {
	repoPath, err := RepoPath()
	if err != nil {
		return "", err
	}
	preCommitHookFilePath := fmt.Sprintf("%s/.git/hooks/pre-commit", repoPath)

	if _, err := os.Stat(preCommitHookFilePath); err != nil {
		return "", nil
	}

	res, err := exec.Command("bash", preCommitHookFilePath).CombinedOutput()

	return string(res), err
}

type commitOptions struct {
	message  string
	noVerify bool
}

func (co *commitOptions) Cmd() *exec.Cmd {
	args := []string{"commit"}
	if co.noVerify == true {
		args = append(args, "--no-verify")
	}

	args = append(args, "-m", co.message)

	return exec.Command("git", args...)
}

type CommitOptionsFunc func(o *commitOptions) error

func WithNoVerify() CommitOptionsFunc {
	return func(o *commitOptions) error {
		o.noVerify = true
		return nil
	}
}

func WithMessage(msg string) CommitOptionsFunc {
	return func(o *commitOptions) error {
		if msg == "" {
			return errors.New("Message cannot be empty.")
		}
		o.message = msg
		return nil
	}
}

func NewCommit(msg string, opts ...CommitOptionsFunc) (*commitOptions, error) {
	opts = append(opts, WithMessage(msg))

	var co commitOptions
	for _, opt := range opts {
		err := opt(&co)
		if err != nil {
			return nil, err
		}
	}
	return &co, nil
}

// func Commit(msg string, optFuncs ...CommitOptionsFunc) (output string, err error) {
// 	co, err := NewCommit(msg, optFuncs...)
// 	if err != nil {
// 		return "", err
// 	}

// 	arguments := []string{
// 		"commit",
// 	}
// 	arguments = append(arguments, co.Arguments()...)

// 	res, err := exec.Command("git", arguments...).CombinedOutput()
// 	return string(res), err
// }
