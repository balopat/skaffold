package shared

import (
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build/plugin/proto"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build/tag"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/color"
	"golang.org/x/net/context"
	"os"
)


// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl build.Builder
}

func (s *GRPCServer) Labels(ctx context.Context, req *proto.LabelsRequest) (*proto.LabelsResponse, error) {
	return &proto.LabelsResponse{
		Labels: s.Impl.Labels(),
	}, nil
}

func (s *GRPCServer) Build(ctx context.Context, req *proto.BuildRequest) (*proto.BuildResponse, error) {
	//TODO: this is cheating - os.Stderr should be something better maybe from the plugin config itself?
	// tagger should not be here - as soon as taggers are decoupled,  this will go away
	color.Cyan.Fprintf(os.Stderr, "Stderr: %v", os.Stderr)
	artifacts, err := s.Impl.Build(ctx, os.Stderr, &tag.ChecksumTagger{}, fromProtoArtifacts(req.Artifacts))
	errorString := "no error"
	if err != nil {
		errorString = err.Error()
	}
	return &proto.BuildResponse {
		BuildResults: toProtoBuildResults(artifacts),
		Error:        errorString,
	}, err
}


func (s *GRPCServer) Validate(context.Context, *proto.ValidationRequest) (*proto.ValidationResponse, error) {
	panic("implement me")
}


