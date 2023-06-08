module rbackup

go 1.19

require (
	github.com/cenkalti/backoff/v4 v4.2.1
	github.com/cespare/xxhash/v2 v2.2.0
	github.com/elithrar/simple-scrypt v1.3.0
	github.com/go-ole/go-ole v1.2.6
	github.com/google/go-cmp v0.5.9
	github.com/minio/sha256-simd v1.0.1
	github.com/pkg/errors v0.9.1
	github.com/pkg/xattr v0.4.9
	github.com/restic/chunker v0.4.0
	golang.org/x/crypto v0.9.0
	golang.org/x/sync v0.2.0
	golang.org/x/sys v0.8.0
)

require github.com/klauspost/cpuid/v2 v2.2.3 // indirect
