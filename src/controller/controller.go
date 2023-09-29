package controller

import (
	"github.com/zilliangroup/builder-backend/src/drive"
	"github.com/zilliangroup/builder-backend/src/storage"
	"github.com/zilliangroup/builder-backend/src/utils/accesscontrol"
	"github.com/zilliangroup/builder-backend/src/utils/tokenvalidator"
)

type Controller struct {
	Storage               *storage.Storage
	Drive                 *drive.Drive
	RequestTokenValidator *tokenvalidator.RequestTokenValidator
	AttributeGroup        *accesscontrol.AttributeGroup
}

func NewControllerForBackend(storage *storage.Storage, drive *drive.Drive, validator *tokenvalidator.RequestTokenValidator, attrg *accesscontrol.AttributeGroup) *Controller {
	return &Controller{
		Storage:               storage,
		Drive:                 drive,
		RequestTokenValidator: validator,
		AttributeGroup:        attrg,
	}
}

func NewControllerForBackendInternal(storage *storage.Storage, drive *drive.Drive, validator *tokenvalidator.RequestTokenValidator, attrg *accesscontrol.AttributeGroup) *Controller {
	return &Controller{
		Storage:               storage,
		Drive:                 drive,
		RequestTokenValidator: validator,
		AttributeGroup:        attrg,
	}
}
