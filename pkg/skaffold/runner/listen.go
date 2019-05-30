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
	"io"
	"time"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/color"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/kubernetes"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/sync"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// TriggerDebounceWindow is the time after a deploy trigger is received
// in which we allow build triggers to take precedence.
const TriggerDebounceWindow = 100 * time.Millisecond

// Listen selects from the component triggers, and performs their actions when received.
func (r *SkaffoldRunner) Listen(ctx context.Context, out io.Writer, buildTrigger, deployTrigger, syncTrigger <-chan bool, logger *kubernetes.LogAggregator) error {
	if err := r.computeChanges(false); err != nil {
		// fail if we can't compute changes successfully the first time
		return errors.Wrap(err, "computing first change set")
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-syncTrigger:
			if err := r.doSync(ctx, out); err != nil {
				logrus.Error(err.Error())
			}
		case <-buildTrigger:
			if err := r.doBuild(ctx, out, logger); err != nil {
				logrus.Error(err.Error())
			}
		case <-deployTrigger:
		debounceLoop:
			for {
				select {
				case <-buildTrigger:
					// give build triggers a chance to be processed before we deploy.
					// this ensures that if we get build and deploy triggers sent at the same time,
					// the build will always happen first.
					if err := r.doBuild(ctx, out, logger); err != nil {
						logrus.Errorf("skipping deploy due to build failure: %s", err.Error())
					}
					break debounceLoop
				case <-ctx.Done():
					return nil
				case <-time.After(TriggerDebounceWindow):
					break debounceLoop
				}
			}

			if err := r.doDeploy(ctx, out, logger); err != nil {
				logrus.Error(err.Error())
			}
		}
	}
}

func (r *SkaffoldRunner) doSync(ctx context.Context, out io.Writer) error {
	if err := r.computeChanges(false); err != nil {
		return errors.Wrap(err, "computing changes after sync trigger")
	}
	if err := r.processDirtyArtifacts(); err != nil {
		return errors.Wrap(err, "processing dirty artifacts")
	}
	for _, s := range r.changeSet.needsResync {
		color.Default.Fprintf(out, "Syncing %d files for %s\n", len(s.Copy)+len(s.Delete), s.Image)

		if err := r.Syncer.Sync(ctx, s); err != nil {
			return errors.Wrap(err, "performing sync")
		}
	}

	r.changeSet.resetSync()
	return nil
}

func (r *SkaffoldRunner) doBuild(ctx context.Context, out io.Writer, logger *kubernetes.LogAggregator) error {
	debounce := r.BuildTrigger.Debounce() // debounce rapid file changes to prevent multiple rebuilds
	if err := r.computeChanges(debounce); err != nil {
		return errors.Wrap(err, "computing changes after build trigger")
	}
	if err := r.processDirtyArtifacts(); err != nil {
		return errors.Wrap(err, "processing dirty artifacts")
	}
	if len(r.changeSet.needsRebuild) == 0 {
		return nil
	}
	logger.Mute()
	defer logger.Unmute()
	defer r.changeSet.resetBuild()
	if err := r.buildAndTest(ctx, out, r.changeSet.needsRebuild); err != nil {
		return errors.Wrap(err, "performing build")
	}
	// after we build, we always need a redeploy
	r.changeSet.needsRedeploy = true
	return nil
}

func (r *SkaffoldRunner) doDeploy(ctx context.Context, out io.Writer, logger *kubernetes.LogAggregator) error {
	if err := r.computeChanges(false); err != nil {
		return errors.Wrap(err, "computing changes after deploy trigger")
	}
	if !r.changeSet.needsRedeploy {
		return nil
	}
	logger.Mute()
	r.changeSet.resetDeploy()
	if err := r.deploy(ctx, out, r.builds); err != nil {
		return errors.Wrap(err, "performing deploy")
	}
	logger.Unmute() // don't defer this, because we only want to unmute after successful deploy
	return nil
}

func (r *SkaffoldRunner) processDirtyArtifacts() error {
	for _, a := range r.changeSet.dirtyArtifacts {
		s, err := sync.NewItem(a.artifact, a.events, r.builds, r.runCtx.InsecureRegistries)
		if err != nil {
			return errors.Wrap(err, "processing sync entry")
		}
		if s != nil {
			r.changeSet.AddResync(s)
		} else {
			r.changeSet.AddRebuild(a.artifact)
		}
	}
	r.changeSet.dirtyArtifacts = nil
	return nil
}

func (r *SkaffoldRunner) computeChanges(debounce bool) error {
	if err := r.FileWatcher.ComputeChanges(debounce); err != nil {
		return errors.Wrap(err, "computing changes")
	}
	if r.changeSet.needsReload {
		return ErrorConfigurationChanged
	}
	return nil
}
