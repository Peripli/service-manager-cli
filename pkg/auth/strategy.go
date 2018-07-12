/*
 * Copyright 2018 The Service Manager Authors
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package auth

import (
	"net/http"
	"time"
)

// Token contains the structure of a typical UAA response token
type Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    time.Time `json:"expires_in"`
	Scope        string    `json:"scope"`
}

// AuthenticationStrategy should be implemented for different authentication strategies
//go:generate counterfeiter . AuthenticationStrategy
type AuthenticationStrategy interface {
	Authenticate(user, password string) (*Token, error)
}

// TokenRefresher should be implemented for different token refresh strategies
//go:generate counterfeiter . TokenRefresher
type TokenRefresher interface {
	Refresh(Token) (*Token, error)
	Client(*Token) *http.Client
}
