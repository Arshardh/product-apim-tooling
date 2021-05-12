/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package impl

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// GetAPIProductListFromEnv
// @param accessToken : Access Token for the environment
// @param environment : Environment where API Product should be imported to
// @param query : String to be matched against the API Product names
// @param limit : Total number of API Products to return
// @return count (no. of API Products)
// @return array of API Product objects
// @return error
func GetAPIProductListFromEnv(accessToken, environment, query, limit string) (count int32, apiProducts []utils.APIProduct, err error) {
	unifiedSearchEndpoint := utils.GetUnifiedSearchEndpointOfEnv(environment, utils.MainConfigFilePath)
	return GetAPIProductList(accessToken, unifiedSearchEndpoint, query, limit)
}

// GetAPIProductList Get the list of API Products available in a particular environment
// @param accessToken : Access Token for the environment
// @param unifiedSearchEndpoint : Unified Search Endpoint for the environment to retreive API Product list
// @param query : String to be matched against the API Product names
// @return count (no. of API Products)
// @return array of API Product objects
// @return error
func GetAPIProductList(accessToken, unifiedSearchEndpoint, query, limit string) (count int32, apiProducts []utils.APIProduct, err error) {
	// Unified Search endpoint from the config file to search API Products
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	// To filter API Products from unified search
	queryParamString := "query=type:\"" + utils.DefaultApiProductType + "\""

	// Setting up the query parameter and limit parameter
	if query != "" {
		queryParamString += " " + query
	}
	if limit != "" {
		queryParamString += "&limit=" + limit
	}
	utils.Logln(utils.LogPrefixInfo+"URL:", unifiedSearchEndpoint+"?"+queryParamString)
	resp, err := utils.InvokeGETRequestWithQueryParamsString(unifiedSearchEndpoint, queryParamString, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+unifiedSearchEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		apiProductListResponse := &utils.APIProductListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &apiProductListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}
		return apiProductListResponse.Count, apiProductListResponse.List, nil
	} else {
		return 0, nil, errors.New(string(resp.Body()))
	}
}
