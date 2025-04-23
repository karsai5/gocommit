package git

import (
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
		fmt.Println("file does not exist: " + preCommitHookFilePath)
		return "", nil
	}

	res, err := exec.Command("bash", preCommitHookFilePath).CombinedOutput()

	return string(res), err
}
