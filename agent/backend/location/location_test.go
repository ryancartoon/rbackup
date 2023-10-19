package location

import (
	"net/url"
	"reflect"
	"testing"

	"rbackup/agent/backend/local"
)

func parseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	return u
}

var parseTests = []struct {
	s string
	u Location
}{
	{
		"local:/srv/repo",
		Location{Scheme: "local",
			Config: local.Config{
				Path:        "/srv/repo",
				Connections: 2,
			},
		},
	},
	{
		"local:dir1/dir2",
		Location{Scheme: "local",
			Config: local.Config{
				Path:        "dir1/dir2",
				Connections: 2,
			},
		},
	},
	{
		"local:dir1/dir2",
		Location{Scheme: "local",
			Config: local.Config{
				Path:        "dir1/dir2",
				Connections: 2,
			},
		},
	},
	{
		"dir1/dir2",
		Location{Scheme: "local",
			Config: local.Config{
				Path:        "dir1/dir2",
				Connections: 2,
			},
		},
	},
	{
		"/dir1/dir2",
		Location{Scheme: "local",
			Config: local.Config{
				Path:        "/dir1/dir2",
				Connections: 2,
			},
		},
	},
	{
		"local:../dir1/dir2",
		Location{Scheme: "local",
			Config: local.Config{
				Path:        "../dir1/dir2",
				Connections: 2,
			},
		},
	},
	{
		"/dir1/dir2",
		Location{Scheme: "local",
			Config: local.Config{
				Path:        "/dir1/dir2",
				Connections: 2,
			},
		},
	},
	{
		"/dir1:foobar/dir2",
		Location{Scheme: "local",
			Config: local.Config{
				Path:        "/dir1:foobar/dir2",
				Connections: 2,
			},
		},
	},
	{
		`\dir1\foobar\dir2`,
		Location{Scheme: "local",
			Config: local.Config{
				Path:        `\dir1\foobar\dir2`,
				Connections: 2,
			},
		},
	},
	{
		`c:\dir1\foobar\dir2`,
		Location{Scheme: "local",
			Config: local.Config{
				Path:        `c:\dir1\foobar\dir2`,
				Connections: 2,
			},
		},
	},
	{
		`C:\Users\appveyor\AppData\Local\Temp\1\restic-test-879453535\repo`,
		Location{Scheme: "local",
			Config: local.Config{
				Path:        `C:\Users\appveyor\AppData\Local\Temp\1\restic-test-879453535\repo`,
				Connections: 2,
			},
		},
	},
	{
		`c:/dir1/foobar/dir2`,
		Location{Scheme: "local",
			Config: local.Config{
				Path:        `c:/dir1/foobar/dir2`,
				Connections: 2,
			},
		},
	},
}

func TestParse(t *testing.T) {
	for i, test := range parseTests {
		t.Run(test.s, func(t *testing.T) {
			u, err := Parse(test.s)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if test.u.Scheme != u.Scheme {
				t.Errorf("test %d: scheme does not match, want %q, got %q",
					i, test.u.Scheme, u.Scheme)
			}

			if !reflect.DeepEqual(test.u.Config, u.Config) {
				t.Errorf("test %d: cfg map does not match, want:\n  %#v\ngot: \n  %#v",
					i, test.u.Config, u.Config)
			}
		})
	}
}

func TestInvalidScheme(t *testing.T) {
	var invalidSchemes = []string{
		"foobar:xxx",
		"foobar:/dir/dir2",
	}

	for _, s := range invalidSchemes {
		t.Run(s, func(t *testing.T) {
			_, err := Parse(s)
			if err == nil {
				t.Fatalf("error for invalid location %q not found", s)
			}
		})
	}
}
