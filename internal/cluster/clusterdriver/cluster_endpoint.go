// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clusterdriver

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/banzaicloud/pipeline/internal/cluster"
)

type deleteClusterRequest struct {
	OrganizationID uint
	ClusterID      uint
	ClusterName    string
	Force          bool
}

// MakeDeleteClusterEndpoint creates an endpoint for a cluster service
func MakeDeleteClusterEndpoint(service cluster.Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(deleteClusterRequest)

		return service.DeleteCluster(
			ctx,
			cluster.Identifier{
				OrganizationID: request.OrganizationID,
				ClusterID:      request.ClusterID,
				ClusterName:    request.ClusterName,
			},
			cluster.DeleteClusterOptions{
				Force: request.Force,
			},
		)
	}
}
