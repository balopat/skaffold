package plugin

import (
	"context"
	"fmt"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build/plugin/shared"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build/tag"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
	"github.com/hashicorp/go-plugin"
	"io"
	"log"
	"os"
	"os/exec"
)

func NewPluginBuilder(pluginName string) build.Builder {
	// We're a host. Start by launching the plugin process.
	log.SetOutput(os.Stdout)

	client := plugin.NewClient(&plugin.ClientConfig{
		//Stderr: os.Stderr,
		SyncStdout: os.Stdout,
		SyncStderr: os.Stderr,
		Managed: true,
		HandshakeConfig:  shared.Handshake,
		Plugins:          shared.PluginMap,
		Cmd:              exec.Command(pluginName),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(pluginName)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	return &BuilderPlugin{
		builder: raw.(build.Builder),
	}
}

type BuilderPlugin struct {
	builder build.Builder
}

// Labels are labels applied to deployed resources.
func (b *BuilderPlugin) Labels() map[string]string {
	fmt.Println("Calling Labels() on plugin!")
	labels := b.builder.Labels()
	fmt.Printf("Result: %s\n", labels)
	return labels
}

func (b *BuilderPlugin) Build(ctx context.Context, out io.Writer, tagger tag.Tagger, artifacts []*latest.Artifact) ([]build.Artifact, error) {
	fmt.Println("Calling Build() on plugin!")
	i, e := b.builder.Build(ctx, out, tagger, artifacts)
	fmt.Printf("Result: %+v, %v\n", i, e)
	return i, e
}
