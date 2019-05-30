/*
Copyright 2019 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package filewatch

import (
	"bytes"
	"testing"
	"time"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/config"
	runcontext "github.com/GoogleContainerTools/skaffold/pkg/skaffold/runner/context"
	"github.com/GoogleContainerTools/skaffold/testutil"
)

func TestNewBuildTrigger(t *testing.T) {
	var tests = []struct {
		description string
		opts        *config.SkaffoldOptions
		expected    Trigger
		shouldErr   bool
	}{
		{
			description: "polling trigger",
			opts:        &config.SkaffoldOptions{BuildTrigger: "polling", WatchPollInterval: 1},
			expected: &pollTrigger{
				Interval: time.Duration(1) * time.Millisecond,
			},
		},
		{
			description: "notify trigger",
			opts:        &config.SkaffoldOptions{BuildTrigger: "notify", WatchPollInterval: 1},
			expected: &fsNotifyTrigger{
				Interval: time.Duration(1) * time.Millisecond,
			},
		},
		{
			description: "manual trigger",
			opts:        &config.SkaffoldOptions{BuildTrigger: "manual"},
			expected:    &manualTrigger{},
		},
		{
			description: "api trigger",
			opts:        &config.SkaffoldOptions{BuildTrigger: "api"},
			expected:    &apiTrigger{},
		},
		{
			description: "unknown trigger",
			opts:        &config.SkaffoldOptions{BuildTrigger: "unknown"},
			shouldErr:   true,
		},
	}
	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			runCtx := &runcontext.RunContext{
				Opts: test.opts,
			}

			got, err := newTrigger(runCtx.Opts.BuildTrigger, runCtx.BuildTrigger, runCtx.Opts.WatchPollInterval)
			t.CheckErrorAndDeepEqual(test.shouldErr, err, test.expected, got)
		})
	}
}

func TestNewTriggers(t *testing.T) {
	var tests = []struct {
		name           string
		opts           *config.SkaffoldOptions
		expectedBuild  Trigger
		expectedDeploy Trigger
		expectedSync   Trigger
		shouldErr      bool
	}{
		{
			name: "all api",
			opts: &config.SkaffoldOptions{
				BuildTrigger:  "api",
				DeployTrigger: "api",
				SyncTrigger:   "api",
			},
			expectedBuild:  &apiTrigger{},
			expectedDeploy: &apiTrigger{},
			expectedSync:   &apiTrigger{},
		},
		{
			name: "one of each",
			opts: &config.SkaffoldOptions{
				BuildTrigger:  "polling",
				DeployTrigger: "api",
				SyncTrigger:   "manual",
			},
			expectedBuild:  &pollTrigger{},
			expectedDeploy: &apiTrigger{},
			expectedSync:   &manualTrigger{},
		},
		{
			name:      "unspecified triggers",
			shouldErr: true,
			opts:      &config.SkaffoldOptions{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runCtx := &runcontext.RunContext{
				Opts: tt.opts,
			}
			b, d, s, err := NewTriggers(runCtx)
			testutil.CheckError(t, tt.shouldErr, err)
			testutil.CheckDeepEqual(t, b, tt.expectedBuild)
			testutil.CheckDeepEqual(t, d, tt.expectedDeploy)
			testutil.CheckDeepEqual(t, s, tt.expectedSync)
		})
	}
}

func TestPollTrigger_Debounce(t *testing.T) {
	trigger := &pollTrigger{}
	got, want := trigger.Debounce(), true
	testutil.CheckDeepEqual(t, want, got)
}

func TestPollTrigger_WatchForChanges(t *testing.T) {
	out := new(bytes.Buffer)

	trigger := &pollTrigger{Interval: 10}
	trigger.WatchForChanges(out)

	got, want := out.String(), "Watching for changes every 10ns...\n"
	testutil.CheckDeepEqual(t, want, got)
}

func TestNotifyTrigger_Debounce(t *testing.T) {
	trigger := &fsNotifyTrigger{}
	got, want := trigger.Debounce(), false
	testutil.CheckDeepEqual(t, want, got)
}

func TestNotifyTrigger_WatchForChanges(t *testing.T) {
	out := new(bytes.Buffer)

	trigger := &fsNotifyTrigger{Interval: 10}
	trigger.WatchForChanges(out)

	got, want := out.String(), "Watching for changes...\n"
	testutil.CheckDeepEqual(t, want, got)
}

func TestManualTrigger_Debounce(t *testing.T) {
	trigger := &manualTrigger{}
	got, want := trigger.Debounce(), false
	testutil.CheckDeepEqual(t, want, got)
}

func TestManualTrigger_WatchForChanges(t *testing.T) {
	out := new(bytes.Buffer)

	trigger := &manualTrigger{}
	trigger.WatchForChanges(out)

	got, want := out.String(), "Press any key to rebuild/redeploy the changes\n"
	testutil.CheckDeepEqual(t, want, got)
}
