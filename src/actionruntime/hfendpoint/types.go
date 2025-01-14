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

package hfendpoint

import "reflect"

type Resource struct {
	Endpoint string `mapstructure:"endpoint" validate:"required"`
	Token    string `mapstructure:"token" validate:"required"`
}

type Action struct {
	Params Parameters `mapstructure:"params"`
}

type Parameters struct {
	Inputs           Inputs  `mapstructure:"inputs"`
	WithDetailParams bool    `mapstructure:"withDetailParams"`
	DetailParams     []Pairs `mapstructure:"detailParams"`
}

type Inputs struct {
	Type    string      `mapstructure:"type"`
	Content interface{} `mapstructure:"content"`
}

type Pairs struct {
	Key   string      `mapstructure:"key"`
	Value interface{} `mapstructure:"value"`
}

func buildDetailedParams(pairs []Pairs) map[string]interface{} {
	res := make(map[string]interface{})
	for _, pair := range pairs {
		switch pair.Key {
		case "useCache":
			vT := reflect.TypeOf(pair.Value)
			value := reflect.ValueOf(pair.Value)
			if vT.Kind() == reflect.Bool {
				res["use_cache"] = value
			}
		case "waitForModel":
			vT := reflect.TypeOf(pair.Value)
			value := reflect.ValueOf(pair.Value)
			if vT.Kind() == reflect.Bool {
				res["wait_for_model"] = value
			}
		case "minLength":
			vT := reflect.TypeOf(pair.Value)
			value := reflect.ValueOf(pair.Value)
			if vT.Kind() == reflect.Int && !value.IsNil() {
				res["min_length"] = value
			}
		case "maxLength":
			vT := reflect.TypeOf(pair.Value)
			value := reflect.ValueOf(pair.Value)
			if vT.Kind() == reflect.Int && !value.IsNil() {
				res["max_length"] = value
			}
		case "topK":
			vT := reflect.TypeOf(pair.Value)
			value := reflect.ValueOf(pair.Value)
			if vT.Kind() == reflect.Int && !value.IsNil() {
				res["top_k"] = value
			}
		case "topP":
			vT := reflect.TypeOf(pair.Value)
			value := reflect.ValueOf(pair.Value)
			if vT.Kind() == reflect.Float64 && !value.IsNil() {
				res["top_p"] = value
			}
		case "temperature":
			vT := reflect.TypeOf(pair.Value)
			value := reflect.ValueOf(pair.Value)
			if vT.Kind() == reflect.Float64 && !value.IsNil() {
				res["temperature"] = value
			}
		case "repetitionPenalty":
			vT := reflect.TypeOf(pair.Value)
			value := reflect.ValueOf(pair.Value)
			if vT.Kind() == reflect.Float64 && !value.IsNil() {
				res["repetition_penalty"] = value
			}
		case "maxTime":
			vT := reflect.TypeOf(pair.Value)
			value := reflect.ValueOf(pair.Value)
			if vT.Kind() == reflect.Float64 && !value.IsNil() {
				res["max_time"] = value
			}
		default:
			break
		}
	}

	return res
}
