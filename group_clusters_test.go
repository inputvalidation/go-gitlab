//
// Copyright 2021, Paul Shoemaker
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
//

package gitlab

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestGroupListClusters(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v4/groups/26/clusters", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprintf(w, `[
			{
			  "id":18,
			  "name":"cluster-1",
			  "domain":"example.com",
			  "created_at":"2019-01-02T20:18:12.563Z",
			  "managed": true,
			  "enabled": true,
			  "provider_type":"user",
			  "platform_type":"kubernetes",
			  "environment_scope":"*",
			  "cluster_type":"group_type",
			  "user":
			  {
				"id":1,
				"name":"Administrator",
				"username":"root",
				"state":"active",
				"avatar_url":"https://www.gravatar.com/avatar/4249f4df72b..",
				"web_url":"https://gitlab.example.com/root"
			  },
			  "platform_kubernetes":
			  {
				"api_url":"https://104.197.68.152",
				"authorization_type":"rbac",
				"ca_cert":"-----BEGIN CERTIFICATE-----\r\nhFiK1L61owwDQYJKoZIhvcNAQELBQAw\r\nLzEtMCsGA1UEAxMkZDA1YzQ1YjctNzdiMS00NDY0LThjNmEtMTQ0ZDJkZjM4ZDBj\r\nMB4XDTE4MTIyNzIwMDM1MVoXDTIzMTIyNjIxMDM1MVowLzEtMCsGA1UEAxMkZDA1\r\nYzQ1YjctNzdiMS00NDY0LThjNmEtMTQ0ZDJkZjM.......-----END CERTIFICATE-----"
			  },
			  "management_project":
			  {
				"id":2,
				"description":null,
				"name":"project2",
				"name_with_namespace":"John Doe8 / project2",
				"path":"project2",
				"path_with_namespace":"namespace2/project2",
				"created_at":"2019-10-11T02:55:54.138Z"
			  }
			},
			{
			  "id":19,
			  "name":"cluster-2"
			}
		  ]`)
	})

	clusters, _, err := client.GroupCluster.ListClusters(26)
	if err != nil {
		t.Errorf("GroupCluster.ListClusters returned error: %v", err)
	}

	timeLayout := "2006-01-02T15:04:05Z07:00"
	createdAt, err := time.Parse(timeLayout, "2019-01-02T20:18:12.563Z")
	if err != nil {
		t.Errorf("DeployKeys.ListAllDeployKeys returned an error while parsing time: %v", err)
	}

	createdAt2, err := time.Parse(timeLayout, "2019-10-11T02:55:54.138Z")
	if err != nil {
		t.Errorf("DeployKeys.ListAllDeployKeys returned an error while parsing time: %v", err)
	}

	want := []*GroupCluster{
		{
			ID:               18,
			Name:             "cluster-1",
			Domain:           "example.com",
			CreatedAt:        &createdAt,
			Managed:          true,
			Enabled:          true,
			ProviderType:     "user",
			PlatformType:     "kubernetes",
			EnvironmentScope: "*",
			ClusterType:      "group_type",
			User: &User{
				ID:        1,
				Name:      "Administrator",
				Username:  "root",
				State:     "active",
				AvatarURL: "https://www.gravatar.com/avatar/4249f4df72b..",
				WebURL:    "https://gitlab.example.com/root",
			},
			PlatformKubernetes: &PlatformKubernetes{
				APIURL:            "https://104.197.68.152",
				AuthorizationType: "rbac",
				CaCert:            "-----BEGIN CERTIFICATE-----\r\nhFiK1L61owwDQYJKoZIhvcNAQELBQAw\r\nLzEtMCsGA1UEAxMkZDA1YzQ1YjctNzdiMS00NDY0LThjNmEtMTQ0ZDJkZjM4ZDBj\r\nMB4XDTE4MTIyNzIwMDM1MVoXDTIzMTIyNjIxMDM1MVowLzEtMCsGA1UEAxMkZDA1\r\nYzQ1YjctNzdiMS00NDY0LThjNmEtMTQ0ZDJkZjM.......-----END CERTIFICATE-----",
			},
			ManagementProject: &ManagementProject{
				ID:                2,
				Description:       nil,
				Name:              "project2",
				NameWithNamespace: "John Doe8 / project2",
				Path:              "project2",
				PathWithNamespace: "namespace2/project2",
				CreatedAt:         &createdAt2,
			},
		},
		{
			ID:   19,
			Name: "cluster-2",
		},
	}
	if !reflect.DeepEqual(want, clusters) {
		t.Errorf("GroupCluster.ListClusters returned %+v, want %+v", clusters, want)
	}
}
