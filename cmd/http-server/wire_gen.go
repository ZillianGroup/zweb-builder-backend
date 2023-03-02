// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/illacloud/illa-builder-backend/api/resthandler"
	"github.com/illacloud/illa-builder-backend/api/router"
	"github.com/illacloud/illa-builder-backend/internal/repository"
	"github.com/illacloud/illa-builder-backend/internal/util"
	"github.com/illacloud/illa-builder-backend/internal/accesscontrol"
	"github.com/illacloud/illa-builder-backend/pkg/action"
	"github.com/illacloud/illa-builder-backend/pkg/builder"
	"github.com/illacloud/illa-builder-backend/pkg/app"
	"github.com/illacloud/illa-builder-backend/pkg/db"
	"github.com/illacloud/illa-builder-backend/pkg/resource"
	"github.com/illacloud/illa-builder-backend/pkg/room"
	"github.com/illacloud/illa-builder-backend/pkg/state"
)

// Injectors from wire.go:

func Initialize() (*Server, error) {
	config, err := GetAppConfig()
	if err != nil {
		return nil, err
	}
	engine := gin.New()
	sugaredLogger := util.NewSugardLogger()
	dbConfig, err := db.GetConfig()
	if err != nil {
		return nil, err
	}
	gormDB, err := db.NewDbConnection(dbConfig, sugaredLogger)
	if err != nil {
		return nil, err
	}
	
	// init supervisior
	attrg, err := accesscontrol.NewRawAttributeGroup()
	if err != nil {
		return nil, err
	}
	appRepositoryImpl := repository.NewAppRepositoryImpl(sugaredLogger, gormDB)
	kvStateRepositoryImpl := repository.NewKVStateRepositoryImpl(sugaredLogger, gormDB)
	treeStateRepositoryImpl := repository.NewTreeStateRepositoryImpl(sugaredLogger, gormDB)
	setStateRepositoryImpl := repository.NewSetStateRepositoryImpl(sugaredLogger, gormDB)
	actionRepositoryImpl := repository.NewActionRepositoryImpl(sugaredLogger, gormDB)	
	resourceRepositoryImpl := repository.NewResourceRepositoryImpl(sugaredLogger, gormDB)
	actionServiceImpl := action.NewActionServiceImpl(sugaredLogger, appRepositoryImpl, actionRepositoryImpl, resourceRepositoryImpl)
	appServiceImpl := app.NewAppServiceImpl(sugaredLogger, appRepositoryImpl, kvStateRepositoryImpl, treeStateRepositoryImpl, setStateRepositoryImpl, actionRepositoryImpl)
	treeStateServiceImpl := state.NewTreeStateServiceImpl(sugaredLogger, treeStateRepositoryImpl)
	// App
	appRestHandlerImpl := resthandler.NewAppRestHandlerImpl(sugaredLogger, appServiceImpl, actionServiceImpl, attrg, treeStateServiceImpl)
	appRouterImpl := router.NewAppRouterImpl(appRestHandlerImpl)
	// public App
	publicAppRestHandlerImpl := resthandler.NewPublicAppRestHandlerImpl(sugaredLogger, appServiceImpl, attrg, treeStateServiceImpl)
	publicAppRouterImpl := router.NewPublicAppRouterImpl(publicAppRestHandlerImpl)
	// room
	roomServiceImpl := room.NewRoomServiceImpl(sugaredLogger)
	roomRestHandlerImpl := resthandler.NewRoomRestHandlerImpl(sugaredLogger, roomServiceImpl, attrg)
	roomRouterImpl := router.NewRoomRouterImpl(roomRestHandlerImpl)
	// resource
	resourceServiceImpl := resource.NewResourceServiceImpl(sugaredLogger, resourceRepositoryImpl)
	resourceRestHandlerImpl := resthandler.NewResourceRestHandlerImpl(sugaredLogger, resourceServiceImpl, attrg)
	resourceRouterImpl := router.NewResourceRouterImpl(resourceRestHandlerImpl)
	// actions
	actionRestHandlerImpl := resthandler.NewActionRestHandlerImpl(sugaredLogger, appServiceImpl, actionServiceImpl, attrg)
	actionRouterImpl := router.NewActionRouterImpl(actionRestHandlerImpl)
	// public actions
	publicActionRestHandlerImpl := resthandler.NewPublicActionRestHandlerImpl(sugaredLogger, actionServiceImpl, attrg)
	publicActionRouterImpl := router.NewPublicActionRouterImpl(publicActionRestHandlerImpl)
	// internalActions
	internalActionRestHandlerImpl := resthandler.NewInternalActionRestHandlerImpl(sugaredLogger, resourceServiceImpl, attrg)
	internalActionRouterImpl := router.NewInternalActionRouterImpl(internalActionRestHandlerImpl)
	// builder
	builderServiceImpl := builder.NewBuilderServiceImpl(sugaredLogger, appRepositoryImpl, resourceRepositoryImpl, actionRepositoryImpl)
	builderRestHandlerImpl := resthandler.NewBuilderRestHandlerImpl(sugaredLogger, builderServiceImpl, attrg)
	builderRouterImpl := router.NewBuilderRouterImpl(builderRestHandlerImpl)
	restRouter := router.NewRESTRouter(sugaredLogger, builderRouterImpl, appRouterImpl, publicAppRouterImpl, roomRouterImpl, actionRouterImpl, publicActionRouterImpl, internalActionRouterImpl, resourceRouterImpl)
	server := NewServer(config, engine, restRouter, sugaredLogger)
	return server, nil
}
