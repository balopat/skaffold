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

package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/filewatch"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/sync"
	"github.com/GoogleContainerTools/skaffold/testutil"
	"k8s.io/client-go/tools/clientcmd/api"
)

type NoopWatcher struct{}

func (t *NoopWatcher) Register(func() ([]string, error), func(filewatch.Events)) error {
	return nil
}

func (t *NoopWatcher) Run(context.Context, io.Writer, func() error) error {
	return nil
}

func (t *NoopWatcher) ComputeChanges(_ bool) error {
	return nil
}

type FailWatcher struct{}

func (t *FailWatcher) Register(func() ([]string, error), func(filewatch.Events)) error {
	return nil
}

func (t *FailWatcher) Run(context.Context, io.Writer, func() error) error {
	return errors.New("BUG")
}

func (t *FailWatcher) ComputeChanges(_ bool) error {
	return errors.New("UH OH")
}

type TestWatcher struct {
	events    []filewatch.Events
	callbacks []func(filewatch.Events)
	testBench *TestBench
}

func (t *TestWatcher) Register(deps func() ([]string, error), onChange func(filewatch.Events)) error {
	t.callbacks = append(t.callbacks, onChange)
	return nil
}

func (t *TestWatcher) ComputeChanges(_ bool) error {
	for _, evt := range t.events {
		for _, file := range evt.Modified {
			_, changed := t.testBench.changeSet.Load(file)
			// _, changed := t.testBench.changeSet[file]
			if !changed {
				return nil
			}
			switch file {
			case "file1":
				t.callbacks[0](evt) // 1st artifact changed
			case "file2":
				t.callbacks[1](evt) // 2nd artifact changed
			case "manifest.yaml":
				t.callbacks[3](evt) // deployment configuration changed
			}
			t.testBench.MarkProcessed(file)
		}
	}

	return nil
}

func TestDevFailFirstCycle(t *testing.T) {
	var tests = []struct {
		description     string
		testBench       *TestBench
		watcher         filewatch.FileWatcher
		expectedActions Actions
	}{
		{
			description:     "fails to build the first time",
			testBench:       NewTestBench().WithBuildErrors([]error{errors.New("")}),
			watcher:         &NoopWatcher{},
			expectedActions: Actions{},
		},
		{
			description: "fails to test the first time",
			testBench:   NewTestBench().WithTestErrors([]error{errors.New("")}),
			watcher:     &NoopWatcher{},
			expectedActions: Actions{
				Built: []string{"img:1"},
			},
		},
		{
			description: "fails to deploy the first time",
			testBench:   NewTestBench().WithDeployErrors([]error{errors.New("")}),
			watcher:     &NoopWatcher{},
			expectedActions: Actions{
				Built:  []string{"img:1"},
				Tested: []string{"img:1"},
			},
		},
		{
			description: "fails to watch after first cycle",
			testBench:   NewTestBench(),
			watcher:     &FailWatcher{},
			expectedActions: Actions{
				Built:    []string{"img:1"},
				Tested:   []string{"img:1"},
				Deployed: []string{"img:1"},
			},
		},
	}
	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			t.SetupFakeKubernetesContext(api.Config{CurrentContext: "cluster1"})

			runner := createRunner(t, test.testBench)
			runner.FileWatcher = test.watcher

			err := runner.Dev(context.Background(), ioutil.Discard, []*latest.Artifact{{
				ImageName: "img",
			}})

			t.CheckErrorAndDeepEqual(true, err, test.expectedActions, test.testBench.Actions())
		})
	}
}

func TestDev(t *testing.T) {
	var tests = []struct {
		description     string
		testBench       *TestBench
		watchEvents     []filewatch.Events
		expectedActions Actions
	}{
		{
			description: "ignore subsequent build errors",
			testBench:   NewTestBench().WithBuildErrors([]error{nil, errors.New("")}),
			watchEvents: []filewatch.Events{
				{Modified: []string{"file1", "file2"}},
			},
			expectedActions: Actions{
				Built:    []string{"img1:1", "img2:1"},
				Tested:   []string{"img1:1", "img2:1"},
				Deployed: []string{"img1:1", "img2:1"},
			},
		},
		{
			description: "ignore subsequent test errors",
			testBench:   NewTestBench().WithTestErrors([]error{nil, errors.New("")}),
			watchEvents: []filewatch.Events{
				{Modified: []string{"file1", "file2"}},
			},
			expectedActions: Actions{
				Built:    []string{"img1:1", "img2:1", "img1:2", "img2:2"},
				Tested:   []string{"img1:1", "img2:1"},
				Deployed: []string{"img1:1", "img2:1"},
			},
		},
		{
			description: "ignore subsequent deploy errors",
			testBench:   NewTestBench().WithDeployErrors([]error{nil, errors.New("")}),
			watchEvents: []filewatch.Events{
				{Modified: []string{"file1", "file2"}},
			},
			expectedActions: Actions{
				Built:    []string{"img1:1", "img2:1", "img1:2", "img2:2"},
				Tested:   []string{"img1:1", "img2:1", "img1:2", "img2:2"},
				Deployed: []string{"img1:1", "img2:1"},
			},
		},
		{
			description: "full cycle twice",
			testBench:   NewTestBench(),
			watchEvents: []filewatch.Events{
				{Modified: []string{"file1", "file2"}},
			},
			expectedActions: Actions{
				Built:    []string{"img1:1", "img2:1", "img1:2", "img2:2"},
				Tested:   []string{"img1:1", "img2:1", "img1:2", "img2:2"},
				Deployed: []string{"img1:1", "img2:1", "img1:2", "img2:2"},
			},
		},
		{
			description: "only change second artifact",
			testBench:   NewTestBench(),
			watchEvents: []filewatch.Events{
				{Modified: []string{"file2"}},
			},
			expectedActions: Actions{
				Built:    []string{"img1:1", "img2:1", "img2:2"},
				Tested:   []string{"img1:1", "img2:1", "img2:2"},
				Deployed: []string{"img1:1", "img2:1", "img2:2", "img1:1"},
			},
		},
		{
			description: "redeploy",
			testBench:   NewTestBench(),
			watchEvents: []filewatch.Events{
				{Modified: []string{"manifest.yaml"}},
			},
			expectedActions: Actions{
				Built:    []string{"img1:1", "img2:1"},
				Tested:   []string{"img1:1", "img2:1"},
				Deployed: []string{"img1:1", "img2:1", "img1:1", "img2:1"},
			},
		},
	}
	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			t.SetupFakeKubernetesContext(api.Config{CurrentContext: "cluster1"})
			test.testBench.currentActions = Actions{}

			runner := createRunner(t, test.testBench)
			runner.FileWatcher = &TestWatcher{
				events:    test.watchEvents,
				testBench: test.testBench,
			}

			var err error
			errChan := make(chan error, 1)
			ctx, cancel := context.WithCancel(context.Background())
			go func(errChan chan error) {
				errChan <- runner.Dev(ctx, ioutil.Discard, []*latest.Artifact{
					{ImageName: "img1"},
					{ImageName: "img2"},
				})
			}(errChan)

			for _, e := range test.watchEvents {
				for _, m := range e.Modified {
					test.testBench.MarkChanged(m)
				}
			}

			timer := time.NewTimer(2 * time.Second)
			<-timer.C
			cancel()
			err = <-errChan

			fmt.Fprintf(os.Stdout, "after dev loop, actions: %+v\n", test.testBench.Actions())

			t.CheckNoError(err)
			t.CheckDeepEqual(test.expectedActions, test.testBench.Actions())
		})
	}
}

func TestDevSync(t *testing.T) {

	var tests = []struct {
		description     string
		testBench       *TestBench
		watchEvents     []filewatch.Events
		expectedActions Actions
	}{
		{
			description: "sync",
			testBench:   NewTestBench(),
			watchEvents: []filewatch.Events{
				{Modified: []string{"file1"}},
			},
			expectedActions: Actions{
				Built:    []string{"img1:1", "img2:1"},
				Tested:   []string{"img1:1", "img2:1"},
				Deployed: []string{"img1:1", "img2:1"},
				Synced:   []string{"img1:1"},
			},
		},
		{
			description: "sync twice",
			testBench:   NewTestBench(),
			watchEvents: []filewatch.Events{
				{Modified: []string{"file1"}},
				{Modified: []string{"file1"}},
			},
			expectedActions: Actions{
				Built:    []string{"img1:1", "img2:1"},
				Tested:   []string{"img1:1", "img2:1"},
				Deployed: []string{"img1:1", "img2:1"},
				Synced:   []string{"img1:1", "img1:1"},
			},
		},
	}

	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			t.SetupFakeKubernetesContext(api.Config{CurrentContext: "cluster1"})

			t.Override(&sync.WorkingDir, func(string, map[string]bool) (string, error) { return "/", nil })

			runner := createRunner(t, test.testBench)
			runner.FileWatcher = &TestWatcher{
				events:    test.watchEvents,
				testBench: test.testBench,
			}

			var err error
			errChan := make(chan error, 1)
			ctx, cancel := context.WithCancel(context.Background())
			go func(errChan chan error) {
				errChan <- runner.Dev(ctx, ioutil.Discard, []*latest.Artifact{
					{
						ImageName: "img1",
						Sync: &latest.Sync{
							Manual: []*latest.SyncRule{{Src: "file1", Dest: "file1"}},
						},
					},
					{
						ImageName: "img2",
					},
				})
			}(errChan)

			debounceTimer := time.NewTimer(500 * time.Millisecond)
			for _, e := range test.watchEvents {
				for _, m := range e.Modified {
					// give repeating sync events time to be processed individually
					<-debounceTimer.C
					debounceTimer.Reset(500 * time.Millisecond)
					test.testBench.MarkChanged(m)
				}
			}

			timer := time.NewTimer(2 * time.Second)
			<-timer.C
			cancel()
			err = <-errChan

			fmt.Fprintf(os.Stdout, "after dev loop, actions: %+v\n", test.testBench.Actions())

			t.CheckNoError(err)
			t.CheckDeepEqual(test.expectedActions, test.testBench.Actions())
		})
	}
}
