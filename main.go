package main

import (
	"context"
	"fmt"
	"os"
	"rbackup/internal/archiver"
	"rbackup/internal/backend/local"
	"rbackup/internal/backend/location"
	"rbackup/internal/backend/logger"
	"rbackup/internal/backend/retry"
	"rbackup/internal/backend/sema"
	"rbackup/internal/debug"
	"rbackup/internal/errors"
	"rbackup/internal/fs"
	"rbackup/internal/repository"
	"rbackup/internal/restic"
	"time"

	"golang.org/x/sync/errgroup"
)

var REPO = "/tmp/rbackup-repo-tmp"

func main() {
	ctx := context.Background()
	runInit(ctx)
	runBackup(ctx, "/tmp/source")
}

func runInit(ctx context.Context) error {

	password := "redhat"

	be, err := create(ctx, REPO)
	if err != nil {
		return errors.Fatalf("create repository failed: %v\n", err)
	}

	compressionOff := repository.CompressionMode(1)

	s, err := repository.New(be, repository.Options{
		Compression: compressionOff,
		PackSize:    4 * 1024 * 1024,
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
func OpenRepository(ctx context.Context, repo string) (*repository.Repository, error) {
	be, err := open(ctx, repo)
	if err != nil {
		return nil, err
	}

	report := func(msg string, err error, d time.Duration) {
		fmt.Printf("%v returned error, retrying after %v: %v\n", msg, d, err)
	}
	success := func(msg string, retries int) {
		fmt.Printf("%v operation successful after %d retries\n", msg, retries)
	}
	be = retry.New(be, 10, report, success)

	compressionOff := repository.CompressionMode(1)
	s, err := repository.New(be, repository.Options{
		Compression: compressionOff,
		PackSize:    4 * 1024 * 1024,
	})
	if err != nil {
		return nil, errors.Fatal(err.Error())
	}

	return s, nil
}

// Open the backend specified by a location config.
func open(ctx context.Context, s string) (restic.Backend, error) {
	debug.Log("parsing location %v", location.StripPassword(s))
	loc, err := location.Parse(s)
	if err != nil {
		return nil, errors.Fatalf("parsing repository location failed: %v", err)
	}

	var be restic.Backend

	cfg, err := parseConfig(loc)
	if err != nil {
		return nil, err
	}

	switch loc.Scheme {
	case "local":
		be, err = local.Open(ctx, cfg.(local.Config))
	default:
		return nil, errors.Fatalf("invalid backend: %q", loc.Scheme)
	}

	if err != nil {
		return nil, errors.Fatalf("unable to open repository at %v: %v", location.StripPassword(s), err)
	}

	// wrap with debug logging and connection limiting
	be = logger.New(sema.NewBackend(be))

	// check if config is there
	fi, err := be.Stat(ctx, restic.Handle{Type: restic.ConfigFile})
	if err != nil {
		return nil, errors.Fatalf("unable to open config file: %v\nIs there a repository at the following location?\n%v", err, location.StripPassword(s))
	}

	if fi.Size == 0 {
		return nil, errors.New("config file has zero size, invalid repository?")
	}

	return be, nil
}

func runBackup(ctx context.Context, target string) error {

	timeStamp := time.Now()
	hostname := "localhost"
	selectByNameFilter := func(item string) bool {
		return true
	}

	selectFilter := func(item string, fi os.FileInfo) bool {
		return true
	}

	repo, err := OpenRepository(ctx, REPO)
	if err != nil {
		return err
	}

	var targetFS fs.FS = fs.Local{}
	var concurrency uint = 1

	wg, wgCtx := errgroup.WithContext(ctx)
	_, cancel := context.WithCancel(wgCtx)
	defer cancel()

	arch := archiver.New(repo, targetFS, archiver.Options{ReadConcurrency: concurrency})
	arch.SelectByName = selectByNameFilter
	arch.Select = selectFilter

	snapshotOpts := archiver.SnapshotOptions{
		Time:     timeStamp,
		Hostname: hostname,
	}

	_, id, err := arch.Snapshot(ctx, []string{target}, snapshotOpts)

	fmt.Printf("snapshot %v saved\n", id.Str())

	// cleanly shutdown all running goroutines
	cancel()

	// let's see if one returned an error
	werr := wg.Wait()

	// return original error
	if err != nil {
		return errors.Fatalf("unable to save snapshot: %v", err)
	}

	// Report finished execution
	return werr
}
