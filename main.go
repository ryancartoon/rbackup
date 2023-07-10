package main

import (
	"context"
	"rbackup/internal/errors"
	"rbackup/internal/options"
	"rbackup/internal/repository"
)

func main() {
	ctx := context.Background()
	runInit(ctx)
}

func runInit(ctx context.Context) error {

	repo := "/tmp/rbackup-repo-tmp"
	password := "redhat"

	be, err := create(ctx, repo, options.Options{})
	if err != nil {
		return errors.Fatalf("create repository failed: %v\n", err)
	}

	compressionOff := repository.CompressionMode(1)

	s, err := repository.New(be, repository.Options{
		Compression: compressionOff,
		PackSize:    1 * 1024 * 1024,
	})
	if err != nil {
		return errors.Fatal(err.Error())
	}
	var version = uint(1)

	err = s.Init(ctx, version, password, nil)
	if err != nil {
		return errors.Fatalf("create key in repository failed: %v\n", err)
	}

	return nil
}

// OpenRepository reads the password and opens the repository.
// func OpenRepository(ctx context.Context, repo string) (*repository.Repository, error) {
// 	be, err := open(ctx, repo, opts, opts.extended)
// 	if err != nil {
// 		return nil, err
// 	}

// 	report := func(msg string, err error, d time.Duration) {
// 		fmt.Printf("%v returned error, retrying after %v: %v\n", msg, d, err)
// 	}
// 	success := func(msg string, retries int) {
// 		fmt.Printf("%v operation successful after %d retries\n", msg, retries)
// 	}
// 	be = retry.New(be, 10, report, success)

// 	// wrap backend if a test specified a hook
// 	if opts.backendTestHook != nil {
// 		be, err = opts.backendTestHook(be)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	s, err := repository.New(be, repository.Options{
// 		Compression: opts.Compression,
// 		PackSize:    opts.PackSize * 1024 * 1024,
// 	})
// 	if err != nil {
// 		return nil, errors.Fatal(err.Error())
// 	}

// 	passwordTriesLeft := 1
// 	if stdinIsTerminal() && opts.password == "" {
// 		passwordTriesLeft = 3
// 	}

// 	for ; passwordTriesLeft > 0; passwordTriesLeft-- {
// 		opts.password, err = ReadPassword(opts, "enter password for repository: ")
// 		if err != nil && passwordTriesLeft > 1 {
// 			opts.password = ""
// 			fmt.Printf("%s. Try again\n", err)
// 		}
// 		if err != nil {
// 			continue
// 		}

// 		err = s.SearchKey(ctx, opts.password, maxKeys, opts.KeyHint)
// 		if err != nil && passwordTriesLeft > 1 {
// 			opts.password = ""
// 			fmt.Fprintf(os.Stderr, "%s. Try again\n", err)
// 		}
// 	}
// 	if err != nil {
// 		if errors.IsFatal(err) {
// 			return nil, err
// 		}
// 		return nil, errors.Fatalf("%s", err)
// 	}

// 	if stdoutIsTerminal() && !opts.JSON {
// 		id := s.Config().ID
// 		if len(id) > 8 {
// 			id = id[:8]
// 		}
// 		if !opts.JSON {
// 			extra := ""
// 			if s.Config().Version >= 2 {
// 				extra = ", compression level " + opts.Compression.String()
// 			}
// 			Verbosef("repository %v opened (version %v%s)\n", id, s.Config().Version, extra)
// 		}
// 	}

// 	if opts.NoCache {
// 		return s, nil
// 	}

// 	c, err := cache.New(s.Config().ID, opts.CacheDir)
// 	if err != nil {
// 		Warnf("unable to open cache: %v\n", err)
// 		return s, nil
// 	}

// 	if c.Created && !opts.JSON && stdoutIsTerminal() {
// 		Verbosef("created new cache in %v\n", c.Base)
// 	}

// 	// start using the cache
// 	s.UseCache(c)

// 	oldCacheDirs, err := cache.Old(c.Base)
// 	if err != nil {
// 		Warnf("unable to find old cache directories: %v", err)
// 	}

// 	// nothing more to do if no old cache dirs could be found
// 	if len(oldCacheDirs) == 0 {
// 		return s, nil
// 	}

// 	// cleanup old cache dirs if instructed to do so
// 	if opts.CleanupCache {
// 		if stdoutIsTerminal() && !opts.JSON {
// 			Verbosef("removing %d old cache dirs from %v\n", len(oldCacheDirs), c.Base)
// 		}
// 		for _, item := range oldCacheDirs {
// 			dir := filepath.Join(c.Base, item.Name())
// 			err = fs.RemoveAll(dir)
// 			if err != nil {
// 				Warnf("unable to remove %v: %v\n", dir, err)
// 			}
// 		}
// 	} else {
// 		if stdoutIsTerminal() {
// 			Verbosef("found %d old cache directories in %v, run `restic cache --cleanup` to remove them\n",
// 				len(oldCacheDirs), c.Base)
// 		}
// 	}

// 	return s, nil
// }
