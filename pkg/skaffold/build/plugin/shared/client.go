package shared

import (
	"context"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build/plugin/proto"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build/tag"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
	"github.com/sirupsen/logrus"
	"io"
)

// GRPCClient is an implementation of proto.BuilderClient that talks over RPC.
type GRPCClient struct {
	Client proto.BuilderClient
}

func (c *GRPCClient) Labels() map[string]string {
	resp, err := c.Client.Labels(context.Background(), &proto.LabelsRequest{})
	if err != nil {
		logrus.Errorf("querying labels plugin: %s", err)
		return nil
	}
	return resp.Labels
}

func (c *GRPCClient) Build(ctx context.Context, out io.Writer, tagger tag.Tagger, artifacts []*latest.Artifact) ([]build.Artifact, error) {
	resp, err := c.Client.Build(ctx, &proto.BuildRequest{
		Artifacts: toProtoArtifacts(artifacts),
	})
	if err != nil {
		logrus.Errorf("build plugin: %s", err)
		return nil, err
	}
	return fromProtoBuildResults(resp.BuildResults), err
}
