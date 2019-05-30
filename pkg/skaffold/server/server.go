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

package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/constants"
	runcontext "github.com/GoogleContainerTools/skaffold/pkg/skaffold/runner/context"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/server/proto"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/util"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var once sync.Once

type server struct {
	buildTrigger  chan bool
	deployTrigger chan bool
	syncTrigger   chan bool
}

func newGRPCServer(port int, buildTrigger, deployTrigger, syncTrigger chan bool) (func() error, error) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", util.Loopback, port))
	if err != nil {
		return func() error { return nil }, errors.Wrap(err, "creating listener")
	}
	logrus.Infof("starting gRPC server on port %d", port)

	s := grpc.NewServer()
	proto.RegisterSkaffoldServiceServer(s, &server{
		buildTrigger:  buildTrigger,
		deployTrigger: deployTrigger,
		syncTrigger:   syncTrigger,
	})

	go func() {
		if err := s.Serve(l); err != nil {
			logrus.Errorf("failed to start grpc server: %s", err)
		}
	}()
	return func() error {
		s.Stop()
		return l.Close()
	}, nil
}

func newHTTPServer(port, proxyPort int) (func() error, error) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := proto.RegisterSkaffoldServiceHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf("%s:%d", util.Loopback, proxyPort), opts)
	if err != nil {
		return func() error { return nil }, err
	}

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", util.Loopback, port))
	if err != nil {
		return func() error { return nil }, errors.Wrap(err, "creating listener")
	}
	logrus.Infof("starting gRPC HTTP server on port %d", port)

	go http.Serve(l, mux)

	return l.Close, nil
}

// Initialize creates the gRPC and HTTP servers for serving the state and event log.
// It returns a shutdown callback for tearing down the grpc server,
// which the runner is responsible for calling.
func Initialize(runctx *runcontext.RunContext) (func() error, error) {
	var callback func() error
	var err error
	once.Do(func() {
		callback, err = initialize(runctx)
	})
	return callback, err
}

func initialize(runctx *runcontext.RunContext) (func() error, error) {
	originalRPCPort := runctx.Opts.RPCPort
	if originalRPCPort == -1 {
		return func() error { return nil }, nil
	}
	rpcPort := util.GetAvailablePort(originalRPCPort, &sync.Map{})
	if rpcPort != originalRPCPort && originalRPCPort != constants.DefaultRPCPort {
		logrus.Warnf("provided port %d already in use: using %d instead", originalRPCPort, rpcPort)
	}
	grpcCallback, err := newGRPCServer(rpcPort, runctx.BuildTrigger, runctx.DeployTrigger, runctx.SyncTrigger)
	if err != nil {
		return grpcCallback, errors.Wrap(err, "starting gRPC server")
	}
	m := &sync.Map{}
	m.Store(rpcPort, true)

	originalHTTPPort := runctx.Opts.RPCHTTPPort
	httpPort := util.GetAvailablePort(originalHTTPPort, m)
	if httpPort != originalHTTPPort && originalHTTPPort != constants.DefaultRPCHTTPPort {
		logrus.Warnf("provided port %d already in use: using %d instead", originalHTTPPort, httpPort)
	}

	httpCallback, err := newHTTPServer(httpPort, rpcPort)
	callback := func() error {
		httpErr := httpCallback()
		grpcErr := grpcCallback()
		errStr := ""
		if grpcErr != nil {
			errStr += fmt.Sprintf("grpc callback error: %s\n", grpcErr.Error())
		}
		if httpErr != nil {
			errStr += fmt.Sprintf("http callback error: %s\n", httpErr.Error())
		}
		return errors.New(errStr)
	}
	if err != nil {
		return callback, errors.Wrap(err, "starting HTTP server")
	}
	return callback, nil
}
