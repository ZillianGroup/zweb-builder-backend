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

type ZwebMarketplaceRestAPI struct {
	Config    *config.Config
	Validator *tokenvalidator.RequestTokenValidator
	Debug     bool `json:"-"`
}

func NewZwebMarketplaceRestAPI() *ZwebMarketplaceRestAPI {
	requestTokenValidator := tokenvalidator.NewRequestTokenValidator()
	return &ZwebMarketplaceRestAPI{
		Config:    config.GetInstance(),
		Validator: requestTokenValidator,
	}
}

func (r *ZwebMarketplaceRestAPI) CloseDebug() {
	r.Debug = false
}

func (r *ZwebMarketplaceRestAPI) OpenDebug() {
	r.Debug = true
}

func (r *ZwebMarketplaceRestAPI) ForkCounter(productType string, productID int) error {
	// self-hist need skip this method.
	if !r.Config.IsCloudMode() {
		return nil
	}
	client := resty.New()
	resp, err := client.R().
		SetHeader("Request-Token", r.Validator.GenerateValidateToken(fmt.Sprintf("%d", productID))).
		Post(r.Config.ZwebMarketplaceInternalRestAPI + fmt.Sprintf(FORK_COUNTER_API, productType, productID))
	if r.Debug {
		log.Printf("[ZwebMarketplaceRestAPI.ForkCounter()]  response: %+v, err: %+v", resp, err)
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		if err != nil {
			return err
		}
		return errors.New(resp.String())
	}
	return nil
}

func (r *ZwebMarketplaceRestAPI) DeleteTeamAllProducts(teamID int) error {
	// self-hist need skip this method.
	if !r.Config.IsCloudMode() {
		return nil
	}
	client := resty.New()
	resp, err := client.R().
		SetHeader("Request-Token", r.Validator.GenerateValidateToken(fmt.Sprintf("%d", teamID))).
		Delete(r.Config.ZwebMarketplaceInternalRestAPI + fmt.Sprintf(DELETE_TEAM_ALL_PRODUCTS, teamID))
	if r.Debug {
		log.Printf("[ZwebMarketplaceRestAPI.DeleteTeamAllProducts()]  response: %+v, err: %+v", resp, err)
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		if err != nil {
			return err
		}
		return errors.New(resp.String())
	}
	return nil
}

func (r *ZwebMarketplaceRestAPI) DeleteProduct(productType string, productID int) error {
	// self-hist need skip this method.
	if !r.Config.IsCloudMode() {
		return nil
	}
	client := resty.New()
	resp, err := client.R().
		SetHeader("Request-Token", r.Validator.GenerateValidateToken(fmt.Sprintf("%d", productID))).
		Delete(r.Config.ZwebMarketplaceInternalRestAPI + fmt.Sprintf(DELETE_PRODUCT, productType, productID))
	if r.Debug {
		log.Printf("[ZwebMarketplaceRestAPI.DeleteProduct()]  response: %+v, err: %+v", resp, err)
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		if err != nil {
			return err
		}
		return errors.New(resp.String())
	}
	return nil
}

func (r *ZwebMarketplaceRestAPI) UpdateProduct(productType string, productID int, product interface{}) error {
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
		Put(r.Config.ZwebMarketplaceInternalRestAPI + fmt.Sprintf(UPDATE_PRODUCTS, productType, productID))
	log.Printf("[ZwebMarketplaceRestAPI.UpdateProduct()]  response: %+v, err: %+v", resp, err)
	if r.Debug {
		log.Printf("[ZwebMarketplaceRestAPI.UpdateProduct()]  response: %+v, err: %+v", resp, err)
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		if err != nil {
			return err
		}
		return errors.New(resp.String())
	}
	return nil
}
