// Copyright 2022 The ILLA Authors.
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

package model

import (
	"time"
)

type GetBuilderDescResponse struct {
	AppNum         int       `json:"appNum"`
	ResourceNum    int       `json:"resourceNum"`
	ActionNum      int       `json:"actionNum"`
	LastModifiedAt time.Time `json:"lastModifiedAt"`
}

type EmptyBuilderDescResponse struct {
	AppNum         int    `json:"appNum"`
	ResourceNum    int    `json:"resourceNum"`
	ActionNum      int    `json:"actionNum"`
	LastModifiedAt string `json:"lastModifiedAt"` // is "" by first time enter builder.
}

func NewGetBuilderDescResponse(appNum int, resourceNum int, actionNum int, lastModifiedAt time.Time) *GetBuilderDescResponse {
	return &GetBuilderDescResponse{
		AppNum:         appNum,
		ResourceNum:    resourceNum,
		ActionNum:      actionNum,
		LastModifiedAt: lastModifiedAt,
	}
}

func NewEmptyBuilderDescResponse(appNum int, resourceNum int, actionNum int) *EmptyBuilderDescResponse {
	return &EmptyBuilderDescResponse{
		AppNum:      appNum,
		ResourceNum: resourceNum,
		ActionNum:   actionNum,
	}
}
