// Copyright 2023 ZWeb Soft, Inc.
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

package appwrite

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/zilliangroup/appwrite-sdk-go/appwrite"
)

func (a *Connector) getClientWithOpts(Opts map[string]interface{}) (*appwrite.Databases, error) {
	// format resource options
	if err := mapstructure.Decode(Opts, &a.Resource); err != nil {
		return nil, err
	}

	// create appwrite client
	client := appwrite.NewClient()
	client.SetEndpoint(a.Resource.Host)
	client.SetProject(a.Resource.ProjectID)
	client.SetKey(a.Resource.APIKey)

	// create appwrite database service
	database := appwrite.NewDatabases(client)

	return database, nil
}

func modifyMapKeysWithPattern(in map[string]interface{}, pattern string, replacement string) {
	for k, v := range in {
		if strings.Contains(k, pattern) {
			newKey := strings.ReplaceAll(k, pattern, replacement)
			in[newKey] = v
			delete(in, k)
		}
	}
}
