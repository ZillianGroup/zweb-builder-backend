package zwebresourcemanagersdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/zilliangroup/zweb-builder-backend/src/utils/config"
	"github.com/zilliangroup/zweb-builder-backend/src/utils/resourcelist"
	"github.com/zilliangroup/zweb-builder-backend/src/utils/tokenvalidator"
)

const (
	BASEURL = "http://127.0.0.1:8008/api/v1"
	// api route part
	GET_AI_AGENT_INTERNAL_API             = "/api/v1/aiAgent/%d"
	RUN_AI_AGENT_INTERNAL_API             = "/api/v1/aiAgent/%d/run"
	DELETE_TEAM_ALL_AI_AGENT_INTERNAL_API = "/api/v1/teams/%d/aiAgent/all"
)

const (
	PRODUCT_TYPE_AIAGENTS = "aiAgents"
	PRODUCT_TYPE_APPS     = "apps"
	PRODUCT_TYPE_HUBS     = "hubs"
)

type ZWebResourceManagerRestAPI struct {
	Config *config.Config
	Debug  bool `json:"-"`
}

func NewZWebResourceManagerRestAPI() (*ZWebResourceManagerRestAPI, error) {
	return &ZWebResourceManagerRestAPI{
		Config: config.GetInstance(),
	}, nil
}

func (r *ZWebResourceManagerRestAPI) CloseDebug() {
	r.Debug = false
}

func (r *ZWebResourceManagerRestAPI) OpenDebug() {
	r.Debug = true
}

func (r *ZWebResourceManagerRestAPI) GetResource(resourceType int, resourceID int) (map[string]interface{}, error) {
	// self-hist need skip this method.
	if !r.Config.IsCloudMode() {
		return nil, nil
	}
	switch resourceType {
	case resourcelist.TYPE_AI_AGENT_ID:
		return r.GetAIAgent(resourceID)
	default:
		return nil, errors.New("Invalied resource type: " + resourcelist.GetResourceIDMappedType(resourceType))
	}
}

func (r *ZWebResourceManagerRestAPI) RunResource(resourceType int, resourceID int, req map[string]interface{}) (*RunResourceResult, error) {
	// self-hist need skip this method.
	if !r.Config.IsCloudMode() {
		return nil, nil
	}
	switch resourceType {
	case resourcelist.TYPE_AI_AGENT_ID:
		return r.RunAIAgent(req)
	default:
		return nil, errors.New("Invalied resource type: " + resourcelist.GetResourceIDMappedType(resourceType))
	}
}

func (r *ZWebResourceManagerRestAPI) GetAIAgent(aiAgentID int) (map[string]interface{}, error) {
	// self-hist need skip this method.
	if !r.Config.IsCloudMode() {
		return nil, nil
	}
	client := resty.New()
	tokenValidator := tokenvalidator.NewRequestTokenValidator()
	uri := r.Config.GetZWebResourceManagerInternalRestAPI() + fmt.Sprintf(GET_AI_AGENT_INTERNAL_API, aiAgentID)
	resp, errInPost := client.R().
		SetHeader("Request-Token", tokenValidator.GenerateValidateToken(strconv.Itoa(aiAgentID))).
		Get(uri)
	if r.Debug {
		log.Printf("[ZWebResourceManagerRestAPI.GetAiAgent()]  uri: %+v \n", uri)
		log.Printf("[ZWebResourceManagerRestAPI.GetAiAgent()]  response: %+v, err: %+v \n", resp, errInPost)
		log.Printf("[ZWebResourceManagerRestAPI.GetAiAgent()]  resp.StatusCode(): %+v \n", resp.StatusCode())
	}
	if errInPost != nil {
		return nil, errInPost
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return nil, errors.New(resp.String())
	}

	var aiAgent map[string]interface{}
	errInUnMarshal := json.Unmarshal([]byte(resp.String()), &aiAgent)
	if errInUnMarshal != nil {
		return nil, errInUnMarshal
	}
	return aiAgent, nil
}

func (r *ZWebResourceManagerRestAPI) RunAIAgent(req map[string]interface{}) (*RunResourceResult, error) {
	// self-hist need skip this method.
	if !r.Config.IsCloudMode() {
		return nil, nil
	}
	reqInstance, errInNewReq := NewRunAIAgentRequest(req)
	if errInNewReq != nil {
		return nil, errInNewReq
	}
	client := resty.New()
	uri := r.Config.GetZWebResourceManagerInternalRestAPI() + fmt.Sprintf(RUN_AI_AGENT_INTERNAL_API, reqInstance.ExportAIAgentIDInInt())
	if r.Debug {
		log.Printf("[ZWebResourceManagerRestAPI.RunAiAgent()]  uri: %+v \n", uri)
	}
	requestToken := reqInstance.ExportRequestToken()
	log.Printf("[ZWebResourceManagerRestAPI.RunAiAgent()]  uri: %+v \n", uri)
	log.Printf("[reqInstance]  reqInstance: %+v \n", reqInstance)
	fmt.Printf("[requestToken] %+v\n", requestToken)
	resp, errInPost := client.R().
		SetHeader("Request-Token", requestToken).
		SetHeader("Authorization", reqInstance.ExportAuthorization()).
		SetBody(req).
		Post(uri)
	if r.Debug {
		log.Printf("[ZWebResourceManagerRestAPI.RunAiAgent()]  response: %+v, err: %+v \n", resp, errInPost)
		log.Printf("[ZWebResourceManagerRestAPI.RunAiAgent()]  resp.StatusCode(): %+v \n", resp.StatusCode())
	}
	if errInPost != nil {
		return nil, errInPost
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return nil, errors.New(resp.String())
	}

	runResourceResult := NewRunResourceResult()
	errInUnMarshal := json.Unmarshal([]byte(resp.String()), &runResourceResult)
	if errInUnMarshal != nil {
		return nil, errInUnMarshal
	}
	return runResourceResult, nil
}

func (r *ZWebResourceManagerRestAPI) DeleteTeamAllAIAgent(teamID int) error {
	// self-hist need skip this method.
	if !r.Config.IsCloudMode() {
		return nil
	}
	client := resty.New()
	tokenValidator := tokenvalidator.NewRequestTokenValidator()
	uri := r.Config.GetZWebResourceManagerInternalRestAPI() + fmt.Sprintf(DELETE_TEAM_ALL_AI_AGENT_INTERNAL_API, teamID)
	resp, errInDelete := client.R().
		SetHeader("Request-Token", tokenValidator.GenerateValidateToken(strconv.Itoa(teamID))).
		Delete(uri)
	if r.Debug {
		log.Printf("[ZWebResourceManagerRestAPI.DeleteTeamAllAiAgent()]  uri: %+v \n", uri)
		log.Printf("[ZWebResourceManagerRestAPI.DeleteTeamAllAiAgent()]  response: %+v, err: %+v \n", resp, errInDelete)
		log.Printf("[ZWebResourceManagerRestAPI.DeleteTeamAllAiAgent()]  resp.StatusCode(): %+v \n", resp.StatusCode())
	}
	if errInDelete != nil {
		return errInDelete
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return errors.New(resp.String())
	}

	return nil
}
