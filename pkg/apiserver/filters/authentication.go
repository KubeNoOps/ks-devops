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

package filters

import (
	"errors"
	"fmt"
	"net/http"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/klog/v2"

	"github.com/kubesphere/ks-devops/pkg/apiserver/request"
)

// WithAuthentication installs authentication handler to handler chain.
// The following part is a little bit ugly, WithAuthentication also logs user failed login attempt
// if using basic auth. But only treats request with requestURI `/oauth/authorize` as login attempt
func WithAuthentication(handler http.Handler, authRequest authenticator.Request) http.Handler {
	if authRequest == nil {
		klog.Warningf("Authentication is disabled")
		return handler
	}
	s := serializer.NewCodecFactory(runtime.NewScheme()).WithoutConversion()

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		resp, ok, err := authRequest.AuthenticateRequest(req)
		_, _, usingBasicAuth := req.BasicAuth()

		defer func() {
			// if we authenticated successfully, go ahead and remove the bearer token so that no one
			// is ever tempted to use it inside of the API server
			if usingBasicAuth && ok {
				req.Header.Del("Authorization")
			}
		}()

		if err != nil || !ok {
			ctx := req.Context()
			requestInfo, found := request.RequestInfoFrom(ctx)
			if !found {
				responsewriters.InternalError(w, req, errors.New("no RequestInfo found in the context"))
				return
			}
			gv := schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}
			responsewriters.ErrorNegotiated(apierrors.NewUnauthorized(fmt.Sprintf("Unauthorized: %s", err)), s, gv, w, req)
			return
		}

		req = req.WithContext(request.WithUser(req.Context(), resp.User))
		handler.ServeHTTP(w, req)
	})
}
