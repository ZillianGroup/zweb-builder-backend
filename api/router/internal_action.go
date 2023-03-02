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

package router

import (
	"github.com/illacloud/illa-builder-backend/api/resthandler"

	"github.com/gin-gonic/gin"
)

type InternalActionRouter interface {
	InitInternalActionRouter(InternalActionRouter *gin.RouterGroup)
}

type InternalActionRouterImpl struct {
	InternalActionRestHandler resthandler.InternalActionRestHandler
}

func NewInternalActionRouterImpl(InternalActionRestHandler resthandler.InternalActionRestHandler) *InternalActionRouterImpl {
	return &InternalActionRouterImpl{InternalActionRestHandler: InternalActionRestHandler}
}

func (impl InternalActionRouterImpl) InitInternalActionRouter(InternalActionRouter *gin.RouterGroup) {
	InternalActionRouter.POST("/generateSQL", impl.InternalActionRestHandler.GenerateSQL)
}
