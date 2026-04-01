package main

import (
	"context"
	"testing"
)

func TestGreetCheck(t *testing.T) {
	check := &GreetCheck{}
	if check.ID() != "hello-world-check" {
		t.Errorf("expected ID hello-world-check, got %s", check.ID())
	}

	finding, err := check.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if finding != nil {
		t.Error("expected nil finding for GreetCheck")
	}
}

func TestFailingCheck(t *testing.T) {
	check := &FailingCheck{}
	if check.ID() != "failing-check" {
		t.Errorf("expected ID failing-check, got %s", check.ID())
	}

	finding, err := check.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if finding == nil {
		t.Fatal("expected finding for FailingCheck")
	}
	if finding.Severity != "info" {
		t.Errorf("expected severity info, got %s", finding.Severity)
	}
}
