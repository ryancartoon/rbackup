package main

import (
	"context"
	"rbackup/agent/backend/local"
	"rbackup/agent/backend/location"
	"rbackup/agent/backend/logger"
	"rbackup/agent/backend/sema"
	"rbackup/agent/debug"
	"rbackup/agent/errors"
	"rbackup/agent/restic"
)

func create(ctx context.Context, s string) (restic.Backend, error) {
	debug.Log("parsing location %v", s)
	loc, err := location.Parse(s)
	if err != nil {
		return nil, err
	}

	cfg, err := parseConfig(loc)
	if err != nil {
		return nil, err
	}

	var be restic.Backend
	switch loc.Scheme {
	case "local":
		be, err = local.Create(ctx, cfg.(local.Config))
	default:
		debug.Log("invalid repository scheme: %v", s)
		return nil, errors.Fatalf("invalid scheme %q", loc.Scheme)
	}

	if err != nil {
		return nil, err
	}

	return logger.New(sema.NewBackend(be)), nil
}

func parseConfig(loc location.Location) (interface{}, error) {

	switch loc.Scheme {
	case "local":
		cfg := loc.Config.(local.Config)
		return cfg, nil
	}

	return nil, errors.Fatalf("invalid backend: %q", loc.Scheme)
}
