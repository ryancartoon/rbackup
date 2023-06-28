package main

import (
	"context"
	"encoding/json"
	"rbackup/internal/backend/location"
	"rbackup/internal/errors"
	"rbackup/internal/repository"

	"github.com/restic/chunker"
)

func main() {
	ctx := context.Background()
	runIint(ctx)
}

func runInit(ctx context.Context, args []string) error {
	if len(args) > 0 {
		return errors.Fatal("the init command expects no arguments, only options - please see `restic help init` for usage and flags")
	}

	chunkerPolynomial, err := maybeReadChunkerPolynomial(ctx, opts, gopts)
	if err != nil {
		return err
	}

	repo, err := ReadRepo(gopts)
	if err != nil {
		return err
	}

	password := "redhat"

	be, err := create(ctx, repo, gopts.extended)
	if err != nil {
		return errors.Fatalf("create repository at %s failed: %v\n", location.StripPassword(gopts.Repo), err)
	}

	s, err := repository.New(be, repository.Options{
		Compression: gopts.Compression,
		PackSize:    gopts.PackSize * 1024 * 1024,
	})
	if err != nil {
		return errors.Fatal(err.Error())
	}

	err = s.Init(ctx, version, password, chunkerPolynomial)
	if err != nil {
		return errors.Fatalf("create key in repository at %s failed: %v\n", location.StripPassword(gopts.Repo), err)
	}

	if !gopts.JSON {
		Verbosef("created restic repository %v at %s", s.Config().ID[:10], location.StripPassword(gopts.Repo))
		if opts.CopyChunkerParameters && chunkerPolynomial != nil {
			Verbosef(" with chunker parameters copied from secondary repository\n")
		} else {
			Verbosef("\n")
		}
		Verbosef("\n")
		Verbosef("Please note that knowledge of your password is required to access\n")
		Verbosef("the repository. Losing your password means that your data is\n")
		Verbosef("irrecoverably lost.\n")

	} else {
		status := initSuccess{
			MessageType: "initialized",
			ID:          s.Config().ID,
			Repository:  location.StripPassword(gopts.Repo),
		}
		return json.NewEncoder(globalOptions.stdout).Encode(status)
	}

	return nil
}

func maybeReadChunkerPolynomial(ctx context.Context, opts InitOptions, gopts GlobalOptions) (*chunker.Pol, error) {
	if opts.CopyChunkerParameters {
		otherGopts, _, err := fillSecondaryGlobalOpts(opts.secondaryRepoOptions, gopts, "secondary")
		if err != nil {
			return nil, err
		}

		otherRepo, err := OpenRepository(ctx, otherGopts)
		if err != nil {
			return nil, err
		}

		pol := otherRepo.Config().ChunkerPolynomial
		return &pol, nil
	}

	if opts.Repo != "" || opts.RepositoryFile != "" || opts.LegacyRepo != "" || opts.LegacyRepositoryFile != "" {
		return nil, errors.Fatal("Secondary repository must only be specified when copying the chunker parameters")
	}
	return nil, nil
}
