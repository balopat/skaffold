package shared

import (
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build/plugin/proto"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/util"
)

func fromProtoBuildResults(buildResults []*proto.BuildResult) []build.Artifact {
	var results []build.Artifact

	for _, buildResult := range buildResults {
		a := build.Artifact{
			ImageName: buildResult.ImageName,
			Tag: buildResult.Tag,
		}

		results = append(results, a)
	}

	return results
}

func toProtoBuildResults(artifacts []build.Artifact) []*proto.BuildResult {
	var results []*proto.BuildResult

	for _, artifact := range artifacts {
		br := &proto.BuildResult{
			ImageName: artifact.ImageName,
			Tag:       artifact.Tag,
		}

		results = append(results, br)
	}

	return results
}

func fromProtoArtifacts(protos []*proto.Artifact) []*latest.Artifact {
	var results []*latest.Artifact
	for _, p := range protos {

		a := &latest.Artifact{
			ImageName: p.ImageName,
			Workspace: p.Workspace,
			Sync:      p.Sync,
		}

		if p.GetDocker() != nil {
			d := p.GetDocker()
			a.DockerArtifact = &latest.DockerArtifact{
				DockerfilePath: d.DockerfilePath,
				BuildArgs: pointerize(d.BuildArgs),
				CacheFrom: d.CacheFrom,
				Target: d.Target,
			}
		} else {
			panic("only docker artifact is supported in this spike")
		}

		results = append(results, a)
	}
	return results
}

//TODO: this is stupid.
func pointerize(strings map[string]string) map[string]*string {
	result := make(map[string]*string)
	for k, v := range strings {
		if v == "" {
			result[k] = nil
		} else {
			result [k] = util.StringPtr(v)
		}
	}
	return result
}

func toProtoArtifacts(artifacts []*latest.Artifact) []*proto.Artifact {
	var protos []*proto.Artifact
	for _, a := range artifacts {

		p := &proto.Artifact{
			ImageName: a.ImageName,
			Workspace: a.Workspace,
			Sync:      a.Sync,
		}

		if a.DockerArtifact != nil {
			p.Type = &proto.Artifact_Docker{
				Docker: &proto.DockerArtifact{
					DockerfilePath: a.DockerArtifact.DockerfilePath,
					BuildArgs:      dePointerize(a.DockerArtifact.BuildArgs),
					CacheFrom:      a.DockerArtifact.CacheFrom,
					Target:         a.DockerArtifact.Target,
				},
			}
		} else {
			panic("only docker artifact is supported in this spike")
		}

		protos = append(protos, p)
	}
	return protos
}

//TODO: this is stupid.
func dePointerize(strings map[string]*string) map[string]string {
	result := make(map[string]string)
	for k, v := range strings {
		if v == nil {
			result[k] = ""
		} else {
			result [k] = *v
		}
	}

	return result
}
