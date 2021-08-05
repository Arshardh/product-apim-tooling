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

package apim

// Application : Application DTO
type Application struct {
	ApplicationID      string           `json:"applicationId"`
	Name               string           `json:"name"`
	ThrottlingPolicy   string           `json:"throttlingPolicy"`
	Description        string           `json:"description"`
	TokenType          string           `json:"tokenType"`
	Status             string           `json:"status"`
	Groups             []string         `json:"groups"`
	SubscriptionCount  int              `json:"subscriptionCount"`
	Keys               []ApplicationKey `json:"keys"`
	SubscriptionScopes []string         `json:"subscriptionScopes"`
	Owner              string           `json:"owner"`
	HashEnabled        bool             `json:"hashEnabled"`
}

// ApplicationKey : Application Key Details
type ApplicationKey struct {
	ConsumerKey          string                 `json:"consumerKey"`
	ConsumerSecret       string                 `json:"consumerSecret"`
	SupportedGrantTypes  []string               `json:"supportedGrantTypes"`
	CallbackURL          string                 `json:"callbackUrl"`
	KeyState             string                 `json:"keyState"`
	KeyType              string                 `json:"keyType"`
	GroupID              string                 `json:"groupId"`
	Token                ApplicationToken       `json:"token"`
	AdditionalProperties map[string]interface{} `json:"additionalProperties"`
}

// ApplicationToken : Application Token Details
type ApplicationToken struct {
	AccessToken  string   `json:"accessToken"`
	TokenScopes  []string `json:"tokenScopes"`
	ValidityTime int64    `json:"validityTime"`
}

// ApplicationKeysList : Applications list
type ApplicationKeysList struct {
	Count int              `json:"count"`
	List  []ApplicationKey `json:"list"`
}
