/*

 Copyright 2020 The KubeSphere Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

*/

package auth

import (
	"errors"
	"fmt"
	"time"

	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/klog/v2"

	"github.com/kubesphere/ks-devops/pkg/apiserver/authentication/oauth"
	authoptions "github.com/kubesphere/ks-devops/pkg/apiserver/authentication/options"
	"github.com/kubesphere/ks-devops/pkg/client/cache"
	"github.com/kubesphere/ks-devops/pkg/jwt/token"
)

// TokenManagementInterface is the interface for the token management
type TokenManagementInterface interface {
	// Verify verifies a token, and return a User if it's a valid token, otherwise return error
	Verify(token string) (user.Info, error)
	// IssueTo issues a token a User, return error if issuing process failed
	IssueTo(user user.Info) (*oauth.Token, error)
	// RevokeAllUserTokens revoke all user tokens
	RevokeAllUserTokens(username string) error
}

type tokenOperator struct {
	issuer  token.Issuer
	options *authoptions.AuthenticationOptions
	cache   cache.Interface
}

// NewTokenOperator creates an instance of the token operator
func NewTokenOperator(cache cache.Interface, options *authoptions.AuthenticationOptions) TokenManagementInterface {
	operator := &tokenOperator{
		issuer:  token.NewTokenIssuer(options.JwtSecret, options.MaximumClockSkew),
		options: options,
		cache:   cache,
	}
	return operator
}

// Verify verifies the token
func (t tokenOperator) Verify(tokenStr string) (user.Info, error) {
	authenticated, tokenType, err := t.issuer.Verify(tokenStr)
	if err != nil {
		return nil, err
	}
	if t.options.OAuthOptions == nil || t.options.OAuthOptions.AccessTokenMaxAge == 0 ||
		tokenType == token.StaticToken {
		return authenticated, nil
	}
	if err := t.tokenCacheValidate(authenticated.GetName(), tokenStr); err != nil {
		return nil, err
	}
	return authenticated, nil
}

// IssueTo creates a new token for a user
func (t tokenOperator) IssueTo(user user.Info) (*oauth.Token, error) {
	accessTokenExpiresIn := t.options.OAuthOptions.AccessTokenMaxAge
	refreshTokenExpiresIn := accessTokenExpiresIn + t.options.OAuthOptions.AccessTokenInactivityTimeout

	accessToken, err := t.issuer.IssueTo(user, token.AccessToken, accessTokenExpiresIn)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	refreshToken, err := t.issuer.IssueTo(user, token.RefreshToken, refreshTokenExpiresIn)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	result := &oauth.Token{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		RefreshToken: refreshToken,
		ExpiresIn:    int(accessTokenExpiresIn.Seconds()),
	}

	if !t.options.MultipleLogin {
		if err = t.RevokeAllUserTokens(user.GetName()); err != nil {
			klog.Error(err)
			return nil, err
		}
	}

	if accessTokenExpiresIn > 0 {
		if err = t.cacheToken(user.GetName(), accessToken, accessTokenExpiresIn); err != nil {
			klog.Error(err)
			return nil, err
		}
		if err = t.cacheToken(user.GetName(), refreshToken, refreshTokenExpiresIn); err != nil {
			klog.Error(err)
			return nil, err
		}
	}

	return result, nil
}

// RevokeAllUserTokens revokes all tokens of a user
func (t tokenOperator) RevokeAllUserTokens(username string) error {
	pattern := fmt.Sprintf("kubesphere:user:%s:token:*", username)
	if keys, err := t.cache.Keys(pattern); err != nil {
		klog.Error(err)
		return err
	} else if len(keys) > 0 {
		if err := t.cache.Del(keys...); err != nil {
			klog.Error(err)
			return err
		}
	}
	return nil
}

func (t tokenOperator) tokenCacheValidate(username, token string) error {
	key := fmt.Sprintf("kubesphere:user:%s:token:%s", username, token)
	if exist, err := t.cache.Exists(key); err != nil {
		return err
	} else if !exist {
		err = errors.New("token not found in cache")
		klog.V(4).Info(fmt.Errorf("%s: %s", err, token))
		return err
	}
	return nil
}

func (t tokenOperator) cacheToken(username, token string, duration time.Duration) error {
	key := fmt.Sprintf("kubesphere:user:%s:token:%s", username, token)
	if err := t.cache.Set(key, token, duration); err != nil {
		klog.Error(err)
		return err
	}
	return nil
}
