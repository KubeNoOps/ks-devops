// Copyright 2022 KubeSphere Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package template

import (
	"github.com/kubesphere/ks-devops/pkg/kapis/devops/v1alpha3/common"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type handler struct {
	client.Client
}

func newHandler(options *common.Options) *handler {
	return &handler{
		Client: options.GenericClient,
	}
}
