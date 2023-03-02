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

package builder

import (
	"github.com/illacloud/illa-builder-backend/internal/repository"
	"time"
	"go.uber.org/zap"
)

type BuilderService interface {
	GetTeamBuilderDesc(teamID int) (interface{}, error)
}

type BuilderServiceImpl struct {
	logger             *zap.SugaredLogger
	appRepository      repository.AppRepository
	resourceRepository repository.ResourceRepository
	actionRepository   repository.ActionRepository
}

func NewBuilderServiceImpl(logger *zap.SugaredLogger, appRepository repository.AppRepository, resourceRepository repository.ResourceRepository,
	actionRepository repository.ActionRepository) *BuilderServiceImpl {
	return &BuilderServiceImpl{
		logger:             logger,
		appRepository:      appRepository,
		resourceRepository: resourceRepository,
		actionRepository:   actionRepository,
	}
}

func (impl *BuilderServiceImpl) GetTeamBuilderDesc(teamID int) (interface{}, error) {
	appNum, errInFetchAppNum := impl.appRepository.CountAPPByTeamID(teamID)
	if errInFetchAppNum != nil {
		return nil, errInFetchAppNum
	}
	resourceNum, errInFetchResourceNum := impl.resourceRepository.CountResourceByTeamID(teamID)
	if errInFetchResourceNum != nil {
		return nil, errInFetchResourceNum
	}
	actionNum, errInFetchAactionNum := impl.actionRepository.CountActionByTeamID(teamID)
	if errInFetchAactionNum != nil {
		return nil, errInFetchAactionNum
	}
	appLastModifyedAt, errInFetchAppModifyTime := impl.appRepository.RetrieveAppLastModifiedTime(teamID)
	resourceLastModifyedAt, errInFetchResourceModifyTime := impl.resourceRepository.RetrieveResourceLastModifiedTime(teamID)

	// compare time
	var lastModifiedAt time.Time 
	if errInFetchAppModifyTime == nil && errInFetchResourceModifyTime == nil {
		if appLastModifyedAt.Before(resourceLastModifyedAt) {
			lastModifiedAt = resourceLastModifyedAt
		} else {
			lastModifiedAt = appLastModifyedAt
		}
	} else if errInFetchResourceModifyTime != nil {
			lastModifiedAt = appLastModifyedAt
	} else if errInFetchAppModifyTime != nil {
			lastModifiedAt = resourceLastModifyedAt
	}

	if errInFetchAppModifyTime != nil && errInFetchResourceModifyTime != nil {
		return NewEmptyBuilderDescResponse(resourceNum, resourceNum, actionNum), nil
	}

	ret := NewGetBuilderDescResponse(appNum, resourceNum, actionNum, lastModifiedAt)
	return ret, nil
}
