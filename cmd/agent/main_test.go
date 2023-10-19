package main

import (
	"context"
	rtest "rbackup/agent/test"
	"testing"
)

func TestRunInit(t *testing.T) {
	ctx := context.Background()
	err := runInit(ctx)
	if err != nil {
		t.Errorf("runInit() error = %v", err)
	}

	// clean up the repo dir after the test
	rtest.RemoveAll(t, "/tmp/rbackup-repo-tmp")
}
