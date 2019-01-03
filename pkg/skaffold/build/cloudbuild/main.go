package main

import (
	"context"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build/gcb"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build/plugin/shared"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build/tag"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/color"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/constants"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
	"github.com/hashicorp/go-plugin"
	"io"
	"os"
)

var OldErr *os.File

func main() {
	OldErr = os.Stderr
	color.IsTerminal = func(w io.Writer) bool {
		return true
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"cloudbuild": &shared.BuilderGRPCPlugin{Impl: newBuilder()},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

type StdErrLoggingBuilder struct {
	Impl build.Builder
}

func (b *StdErrLoggingBuilder) Labels() map[string]string {
	return b.Impl.Labels()
}

func (b *StdErrLoggingBuilder) Build(ctx context.Context, out io.Writer, tagger tag.Tagger, artifacts []*latest.Artifact) ([]build.Artifact, error) {
	return b.Impl.Build(ctx, OldErr, tagger, artifacts)
}

func newBuilder() build.Builder {
	return &StdErrLoggingBuilder{
		Impl: gcb.NewBuilder(&latest.GoogleCloudBuild{
			ProjectID:   "balintp-gcp-lab",
			DockerImage: constants.DefaultCloudBuildDockerImage,
		}),
	}
}
