package diff

import (
	"reflect"
	"testing"

	"github.com/tdiprima/driftwatch/internal/scanner"
	"github.com/tdiprima/driftwatch/internal/snapshot"
)

func TestCompareDetectsHostChanges(t *testing.T) {
	baseline := snapshot.Snapshot{
		Host: scanner.HostInfo{
			Hostname: "vulcan",
			OS:       "Ubuntu 24.04.2 LTS",
			Kernel:   "6.8.0-49-generic",
		},
	}
	current := snapshot.Snapshot{
		Host: scanner.HostInfo{
			Hostname: "apollo",
			OS:       "Ubuntu 24.04.3 LTS",
			Kernel:   "6.8.0-50-generic",
		},
	}

	got := Compare(baseline, current).Changes
	want := []Change{
		{
			Category: "Host",
			Type:     "changed",
			Detail:   "hostname: vulcan -> apollo",
		},
		{
			Category: "Host",
			Type:     "changed",
			Detail:   "OS: Ubuntu 24.04.2 LTS -> Ubuntu 24.04.3 LTS",
		},
		{
			Category: "Host",
			Type:     "changed",
			Detail:   "kernel: 6.8.0-49-generic -> 6.8.0-50-generic",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Compare() changes = %#v, want %#v", got, want)
	}
}

func TestCompareDetectsUserUIDChanges(t *testing.T) {
	baseline := snapshot.Snapshot{
		Users: []scanner.UserInfo{
			{Username: "deploy", UID: 1001},
			{Username: "root", UID: 0},
		},
	}
	current := snapshot.Snapshot{
		Users: []scanner.UserInfo{
			{Username: "deploy", UID: 0},
			{Username: "root", UID: 0},
		},
	}

	got := Compare(baseline, current).Changes
	want := []Change{
		{
			Category: "User",
			Type:     "changed",
			Detail:   "deploy uid: 1001 -> 0",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Compare() changes = %#v, want %#v", got, want)
	}
}
