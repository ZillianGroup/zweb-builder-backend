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

package snowflake

type Resource struct {
	AccountName    string `validate:"required"`
	Warehouse      string `validate:"required"`
	Database       string `validate:"required"`
	Schema         string
	Role           string
	Authentication string            `validate:"oneof=basic key"`
	AuthContent    map[string]string `validate:"required"`
}

type Action struct {
	Mode  string `validate:"oneof=gui sql"`
	Query string
}
