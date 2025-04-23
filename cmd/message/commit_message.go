package message

import (
	"errors"
	"fmt"
	"strings"
)

type commitMessage struct {
	ticket             *string
	commitType         *string
	descriptionOneLine string
	descriptionLong    *string
}

func (cm *commitMessage) Valid() error {
	if cm.descriptionOneLine == "" {
		return errors.New("One line description cannot be empty")
	}

	return nil
}

func (cm *commitMessage) ApplyOption(opt CommitMessageOption) error {
	return opt(cm)
}

func (cm *commitMessage) Message() string {
	lineOneBlocks := []string{}
	if cm.commitType != nil {
		lineOneBlocks = append(lineOneBlocks, fmt.Sprintf("%s:", *cm.commitType))
	}
	if cm.ticket != nil {
		lineOneBlocks = append(lineOneBlocks, fmt.Sprintf("[%s]", *cm.ticket))
	}
	lineOneBlocks = append(lineOneBlocks, cm.descriptionOneLine)

	lineOne := strings.Join(lineOneBlocks, " ")
	return lineOne
}

type CommitMessageOption func(cm *commitMessage) error

func NewCommitMessage(opts ...CommitMessageOption) (*commitMessage, error) {
	var cm commitMessage

	for _, opt := range opts {
		err := cm.ApplyOption(opt)
		if err != nil {
			return nil, err
		}
	}

	return &cm, nil
}

func WithTicket(ticket string) CommitMessageOption {
	return func(cm *commitMessage) error {
		if ticket != "" {
			cm.ticket = &ticket
		}
		return nil
	}
}

func WithType(commitType string) CommitMessageOption {
	return func(cm *commitMessage) error {
		if commitType != "" {
			cm.commitType = &commitType
		}
		return nil
	}
}

func WithOneLineDescription(msg string) CommitMessageOption {
	return func(cm *commitMessage) error {
		if msg == "" {
			return errors.New("One line description cannot be empty")
		}
		cm.descriptionOneLine = msg
		return nil
	}
}

func WithLongDescription(msg string) CommitMessageOption {
	return func(cm *commitMessage) error {
		if msg != "" {
			cm.descriptionLong = &msg
		}
		return nil
	}
}
