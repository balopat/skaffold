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
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/color"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/filewatch/util"
	runcontext "github.com/GoogleContainerTools/skaffold/pkg/skaffold/runner/context"
	"github.com/pkg/errors"
	"github.com/rjeczalik/notify"
	"github.com/sirupsen/logrus"
)

// Trigger describes a mechanism that triggers the watch.
type Trigger interface {
	Start(context.Context) (<-chan bool, error)
	WatchForChanges(io.Writer)
	Debounce() bool
}

func NewTriggers(runctx *runcontext.RunContext) (Trigger, Trigger, Trigger, error) {
	buildTrigger, err := newTrigger(runctx.Opts.BuildTrigger, runctx.BuildTrigger, runctx.Opts.WatchPollInterval)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "creating build trigger")
	}
	deployTrigger, err := newTrigger(runctx.Opts.DeployTrigger, runctx.DeployTrigger, runctx.Opts.WatchPollInterval)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "creating deploy trigger")
	}
	syncTrigger, err := newTrigger(runctx.Opts.SyncTrigger, runctx.SyncTrigger, runctx.Opts.WatchPollInterval)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "creating sync trigger")
	}
	return buildTrigger, deployTrigger, syncTrigger, nil
}

// NewTrigger creates a new trigger.
func newTrigger(triggerType string, trigger chan bool, interval int) (Trigger, error) {
	switch strings.ToLower(triggerType) {
	case util.Polling:
		return &pollTrigger{
			Interval: time.Duration(interval) * time.Millisecond,
		}, nil
	case util.Notify:
		return &fsNotifyTrigger{
			Interval: time.Duration(interval) * time.Millisecond,
		}, nil
	case util.Manual:
		return &manualTrigger{}, nil
	case util.API:
		return &apiTrigger{
			Trigger: trigger,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported trigger: %s", triggerType)
	}
}

// pollTrigger watches for changes on a given interval of time.
type pollTrigger struct {
	Interval time.Duration
}

// Debounce tells the watcher to debounce rapid sequence of changes.
func (t *pollTrigger) Debounce() bool {
	return true
}

func (t *pollTrigger) WatchForChanges(out io.Writer) {
	color.Yellow.Fprintf(out, "Watching for changes every %v...\n", t.Interval)
}

// Start starts a timer.
func (t *pollTrigger) Start(ctx context.Context) (<-chan bool, error) {
	trigger := make(chan bool)

	ticker := time.NewTicker(t.Interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				trigger <- true
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	return trigger, nil
}

// manualTrigger watches for changes when the user presses a key.
type manualTrigger struct{}

// Debounce tells the watcher to not debounce rapid sequence of changes.
func (t *manualTrigger) Debounce() bool {
	return false
}

func (t *manualTrigger) WatchForChanges(out io.Writer) {
	color.Yellow.Fprintln(out, "Press any key to rebuild/redeploy the changes")
}

// Start starts listening to pressed keys.
func (t *manualTrigger) Start(ctx context.Context) (<-chan bool, error) {
	trigger := make(chan bool)

	var stopped int32
	go func() {
		<-ctx.Done()
		atomic.StoreInt32(&stopped, 1)
	}()

	reader := bufio.NewReader(os.Stdin)
	go func() {
		for {
			_, _, err := reader.ReadRune()
			if err != nil {
				logrus.Debugf("manual trigger error: %s", err)
			}

			// Wait until the context is cancelled.
			if atomic.LoadInt32(&stopped) == 1 {
				return
			}
			trigger <- true
		}
	}()

	return trigger, nil
}

// notifyTrigger watches for changes with fsnotify
type fsNotifyTrigger struct {
	Interval time.Duration
}

// Debounce tells the watcher to not debounce rapid sequence of changes.
func (t *fsNotifyTrigger) Debounce() bool {
	// This trigger has built-in debouncing.
	return false
}

func (t *fsNotifyTrigger) WatchForChanges(out io.Writer) {
	color.Yellow.Fprintln(out, "Watching for changes...")
}

// Start listening for file system changes
func (t *fsNotifyTrigger) Start(ctx context.Context) (<-chan bool, error) {
	c := make(chan notify.EventInfo, 100)

	// Watch current directory recursively
	if err := notify.Watch("./...", c, notify.All); err != nil {
		return nil, err
	}

	trigger := make(chan bool)
	go func() {
		timer := time.NewTimer(1<<63 - 1) // Forever

		for {
			select {
			case e := <-c:
				logrus.Debugln("Change detected", e)

				// Wait t.interval before triggering.
				// This way, rapid stream of events will be grouped.
				timer.Reset(t.Interval)
			case <-timer.C:
				trigger <- true
			case <-ctx.Done():
				timer.Stop()
				return
			}
		}
	}()

	return trigger, nil
}

type apiTrigger struct {
	Trigger chan bool
}

// Start receives triggers from gRPC/HTTP and triggers a rebuild.
func (t *apiTrigger) Start(ctx context.Context) (<-chan bool, error) {
	trigger := make(chan bool)

	go func() {
		for {
			select {
			case <-t.Trigger:
				logrus.Debugln("build request received")
				trigger <- true
			case <-ctx.Done():
				break
			}
		}
	}()

	return trigger, nil
}

func (t *apiTrigger) Debounce() bool {
	return false
}

func (t *apiTrigger) WatchForChanges(out io.Writer) {
	color.Yellow.Fprintln(out, "Watching on designated port for build requests...")
}

func StartTrigger(ctx context.Context, trigger Trigger) (<-chan bool, error) {
	t, err := trigger.Start(ctx)
	if err != nil {
		if notifyTrigger, ok := trigger.(*fsNotifyTrigger); ok {
			trigger = &pollTrigger{
				Interval: notifyTrigger.Interval,
			}

			logrus.Debugln("Couldn't start notify trigger. Falling back to a polling trigger")
			t, err = trigger.Start(ctx)
		}
	}
	return t, err
}
