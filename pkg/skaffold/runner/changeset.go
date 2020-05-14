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
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/sync"
)

type devWorkItems struct {
	needsRebuild  []*latest.Artifact
	needsResync   []*sync.Item
	needsRedeploy bool
	needsReload   bool
}

type changeSet struct {
	work           devWorkItems
	rebuildTracker map[string]*latest.Artifact
	resyncTracker  map[string]*sync.Item
}

func (c *changeSet) AddRebuild(a *latest.Artifact) {
	if _, ok := c.rebuildTracker[a.ImageName]; ok {
		return
	}
	c.rebuildTracker[a.ImageName] = a
	c.work.needsRebuild = append(c.work.needsRebuild, a)
	c.work.needsRedeploy = true
}

func (c *changeSet) AddResync(s *sync.Item) {
	if _, ok := c.resyncTracker[s.Image]; ok {
		return
	}
	c.resyncTracker[s.Image] = s
	c.work.needsResync = append(c.work.needsResync, s)
}

func (c *changeSet) resetBuild() {
	c.rebuildTracker = make(map[string]*latest.Artifact)
	c.work.needsRebuild = nil
}

func (c *changeSet) resetSync() {
	c.resyncTracker = make(map[string]*sync.Item)
	c.work.needsResync = nil
}

func (c *changeSet) resetDeploy() {
	c.work.needsRedeploy = false
}
