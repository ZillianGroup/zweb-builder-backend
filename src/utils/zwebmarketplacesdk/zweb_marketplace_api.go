package zwebmarketplacesdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/zilliangroup/zweb-builder-backend/src/utils/config"
	"github.com/zilliangroup/zweb-builder-backend/src/utils/tokenvalidator"
)

const (
	BASEURL = "http://127.0.0.1:9001/api/v1"
	// api route part
	FORK_COUNTER_API         = "/products/%s/%d/fork"
	RUN_COUNTER_API          = "/products/%s/%d/run"
	DELETE_TEAM_ALL_PRODUCTS = "/products/byTeam/%d"
	DELETE_PRODUCT           = "/products/%s/%d"
	UPDATE_PRODUCTS          = "/products/%s/%d"
)

const (
	PRODUCT_TYPE_AIAGENTS = "aiAgents"
	PRODUCT_TYPE_APPS     = "apps"
	PRODUCT_TYPE_HUBS     = "hubs"
)

type ZWebMarketplaceRestAPI struct {
	Config    *config.Config
	Validator *tokenvalidator.RequestTokenValidator
	Debug     bool `json:"-"`
}

func NewZWebMarketplaceRestAPI() *ZWebMarketplaceRestAPI {
	requestTokenValidator := tokenvalidator.NewRequestTokenValidator()
	return &ZWebMarketplaceRestAPI{
		Config:    config.GetInstance(),
		Validator: requestTokenValidator,
	}
}

func (r *ZWebMarketplaceRestAPI) CloseDebug() {
	r.Debug = false
}

func (r *ZWebMarketplaceRestAPI) OpenDebug() {
	r.Debug = true
}

func (r *ZWebMarketplaceRestAPI) ForkCounter(productType string, productID int) error {
	// self-hist need skip this method.
	if !r.Config.IsCloudMode() {
		return nil
	}
	client := resty.New()
	resp, err := client.R().
		SetHeader("Request-Token", r.Validator.GenerateValidateToken(fmt.Sprintf("%d", productID))).
		Post(r.Config.ZWebMarketplaceInternalRestAPI + fmt.Sprintf(FORK_COUNTER_API, productType, productID))
	if r.Debug {
		log.Printf("[ZWebMarketplaceRestAPI.ForkCounter()]  response: %+v, err: %+v", resp, err)
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		if err != nil {
			return err
		}
		return errors.New(resp.String())
	}
	return nil
}

func (r *ZWebMarketplaceRestAPI) DeleteTeamAllProducts(teamID int) error {
	// self-hist need skip this method.
	if !r.Config.IsCloudMode() {
		return nil
	}
	client := resty.New()
	resp, err := client.R().
		SetHeader("Request-Token", r.Validator.GenerateValidateToken(fmt.Sprintf("%d", teamID))).
		Delete(r.Config.ZWebMarketplaceInternalRestAPI + fmt.Sprintf(DELETE_TEAM_ALL_PRODUCTS, teamID))
	if r.Debug {
		log.Printf("[ZWebMarketplaceRestAPI.DeleteTeamAllProducts()]  response: %+v, err: %+v", resp, err)
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		if err != nil {
			return err
		}
		return errors.New(resp.String())
	}
	return nil
}

func (r *ZWebMarketplaceRestAPI) DeleteProduct(productType string, productID int) error {
	// self-hist need skip this method.
	if !r.Config.IsCloudMode() {
		return nil
	}
	client := resty.New()
	resp, err := client.R().
		SetHeader("Request-Token", r.Validator.GenerateValidateToken(fmt.Sprintf("%d", productID))).
		Delete(r.Config.ZWebMarketplaceInternalRestAPI + fmt.Sprintf(DELETE_PRODUCT, productType, productID))
	if r.Debug {
		log.Printf("[ZWebMarketplaceRestAPI.DeleteProduct()]  response: %+v, err: %+v", resp, err)
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		if err != nil {
			return err
		}
		return errors.New(resp.String())
	}
	return nil
}

func (r *ZWebMarketplaceRestAPI) UpdateProduct(productType string, productID int, product interface{}) error {
	// self-hist need skip this method.
	if !r.Config.IsCloudMode() {
		return nil
	}
	b, err := json.Marshal(product)
	if err != nil {
		return err
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Request-Token", r.Validator.GenerateValidateToken(fmt.Sprintf("%d", productID), string(b))).
		SetBody(product).
		Put(r.Config.ZWebMarketplaceInternalRestAPI + fmt.Sprintf(UPDATE_PRODUCTS, productType, productID))
	log.Printf("[ZWebMarketplaceRestAPI.UpdateProduct()]  response: %+v, err: %+v", resp, err)
	if r.Debug {
		log.Printf("[ZWebMarketplaceRestAPI.UpdateProduct()]  response: %+v, err: %+v", resp, err)
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		if err != nil {
			return err
		}
		return errors.New(resp.String())
	}
	return nil
}
