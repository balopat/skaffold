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

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/filewatch"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/kubernetes"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
	"github.com/pkg/errors"
)

// ErrorConfigurationChanged is a special error that's returned when the skaffold configuration was changed.
var ErrorConfigurationChanged = errors.New("configuration changed")

// Dev watches for changes and runs the skaffold build and deploy
// config until interrupted by the user.
func (r *SkaffoldRunner) Dev(ctx context.Context, out io.Writer, artifacts []*latest.Artifact) error {
	logger := r.newLogger(out, artifacts)
	defer logger.Stop()

	portForwarder := kubernetes.NewPortForwarder(out, r.imageList, r.runCtx.Namespaces)
	defer portForwarder.Stop()

	// Create filewatcher and register artifacts to build current state of files.
	r.changeSet = &changes{}

	// Watch artifacts
	for i := range artifacts {
		artifact := artifacts[i]
		if !r.runCtx.Opts.IsTargetImage(artifact) {
			continue
		}

		if err := r.FileWatcher.Register(
			func() ([]string, error) { return r.Builder.DependenciesForArtifact(ctx, artifact) },
			func(e filewatch.Events) { r.changeSet.AddDirtyArtifact(artifact, e) },
		); err != nil {
			return errors.Wrapf(err, "watching files for artifact %s", artifact.ImageName)
		}
	}

	// Watch test configuration
	if err := r.FileWatcher.Register(
		r.TestDependencies,
		func(filewatch.Events) { r.changeSet.needsRedeploy = true },
	); err != nil {
		return errors.Wrap(err, "watching test files")
	}

	// Watch deployment configuration
	if err := r.FileWatcher.Register(
		r.Dependencies,
		func(filewatch.Events) { r.changeSet.needsRedeploy = true },
	); err != nil {
		return errors.Wrap(err, "watching files for deployer")
	}

	// Watch Skaffold configuration
	if err := r.FileWatcher.Register(
		func() ([]string, error) { return []string{r.runCtx.Opts.ConfigurationFile}, nil },
		func(filewatch.Events) { r.changeSet.needsReload = true },
	); err != nil {
		return errors.Wrapf(err, "watching skaffold configuration %s", r.runCtx.Opts.ConfigurationFile)
	}

	// First run
	if err := r.buildTestDeploy(ctx, out, artifacts); err != nil {
		return errors.Wrap(err, "exiting dev mode because first run failed")
	}

	// Start logs
	if r.runCtx.Opts.TailDev {
		if err := logger.Start(ctx); err != nil {
			return errors.Wrap(err, "starting logger")
		}
	}

	if r.runCtx.Opts.PortForward {
		if err := portForwarder.Start(ctx); err != nil {
			return errors.Wrap(err, "starting port-forwarder")
		}
	}

	triggerCtx, triggerCancel := context.WithCancel(ctx)
	defer triggerCancel()

	buildTrigger, err := filewatch.StartTrigger(triggerCtx, r.BuildTrigger)
	if err != nil {
		return errors.Wrap(err, "unable to start build trigger")
	}
	deployTrigger, err := filewatch.StartTrigger(triggerCtx, r.DeployTrigger)
	if err != nil {
		return errors.Wrap(err, "unable to start deploy trigger")
	}
	syncTrigger, err := filewatch.StartTrigger(triggerCtx, r.SyncTrigger)
	if err != nil {
		return errors.Wrap(err, "unable to start sync trigger")
	}

	return r.Listen(ctx, out, buildTrigger, deployTrigger, syncTrigger, logger)
}
