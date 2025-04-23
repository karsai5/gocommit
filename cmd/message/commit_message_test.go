package message

import "testing"

func TestMessageWithOneLine(t *testing.T) {
	cm, err := NewCommitMessage(WithOneLineDescription("hello world"))
	if err != nil {
		t.Error(err)
	}

	expected := "hello world"
	if msg := cm.Message(); msg != expected {
		t.Errorf(`Commit with only one line: wanted "%s" got "%s"`, expected, msg)
	}
}

func TestMessageWithTicket(t *testing.T) {
	cm, err := NewCommitMessage(
		WithOneLineDescription("hello world"),
		WithTicket("TICKET-123"),
	)
	if err != nil {
		t.Error(err)
	}

	expected := "[TICKET-123] hello world"
	if msg := cm.Message(); msg != expected {
		t.Errorf(`Commit with only one line: wanted "%s" got "%s"`, expected, msg)
	}
}

func TestMessageWithType(t *testing.T) {
	cm, err := NewCommitMessage(
		WithOneLineDescription("hello world"),
		WithType("feat"),
	)
	if err != nil {
		t.Error(err)
	}

	expected := "feat: hello world"
	if msg := cm.Message(); msg != expected {
		t.Errorf(`Commit with only one line: wanted "%s" got "%s"`, expected, msg)
	}
}

func TestMessageWithTypeAndTicket(t *testing.T) {
	cm, err := NewCommitMessage(
		WithOneLineDescription("hello world"),
		WithTicket("TICKET-123"),
		WithType("feat"),
	)
	if err != nil {
		t.Error(err)
	}

	expected := "feat: [TICKET-123] hello world"
	if msg := cm.Message(); msg != expected {
		t.Errorf(`Commit with only one line: wanted "%s" got "%s"`, expected, msg)
	}
}
