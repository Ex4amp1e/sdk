// Copyright (c) 2020 Doc.ai and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package setlogoption implements a chain element to set log options before full chain
package setlogoption

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/sdk/pkg/tools/log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/networkservicemesh/api/pkg/api/registry"
)

type setLogOption struct {
	options map[string]string
	server  registry.NetworkServiceEndpointRegistryServer
}

type setLogOptionFindServer struct {
	registry.NetworkServiceEndpointRegistry_FindServer
	ctx context.Context
}

func (s *setLogOptionFindServer) Send(endpoint *registry.NetworkServiceEndpoint) error {
	return s.NetworkServiceEndpointRegistry_FindServer.Send(endpoint)
}

func (s *setLogOptionFindServer) Context() context.Context {
	return s.ctx
}

func (s *setLogOption) Register(ctx context.Context, endpoint *registry.NetworkServiceEndpoint) (*registry.NetworkServiceEndpoint, error) {
	ctx = s.withFields(ctx)
	return s.server.Register(ctx, endpoint)
}

func (s *setLogOption) Find(query *registry.NetworkServiceEndpointQuery, server registry.NetworkServiceEndpointRegistry_FindServer) error {
	ctx := s.withFields(server.Context())
	return s.server.Find(query, &setLogOptionFindServer{ctx: ctx, NetworkServiceEndpointRegistry_FindServer: server})
}

func (s *setLogOption) Unregister(ctx context.Context, endpoint *registry.NetworkServiceEndpoint) (*empty.Empty, error) {
	ctx = s.withFields(ctx)
	return s.server.Unregister(ctx, endpoint)
}

// NewNetworkServiceEndpointRegistryServer creates new instance of NetworkServiceEndpointRegistryServer which sets the passed options
func NewNetworkServiceEndpointRegistryServer(options map[string]string, server registry.NetworkServiceEndpointRegistryServer) registry.NetworkServiceEndpointRegistryServer {
	return &setLogOption{
		options: options,
		server:  server,
	}
}

func (s *setLogOption) withFields(ctx context.Context) context.Context {
	fields := make(logrus.Fields)
	for k, v := range s.options {
		fields[k] = v
	}
	if len(fields) > 0 {
		ctx = log.WithFields(ctx, fields)
	}
	return ctx
}
