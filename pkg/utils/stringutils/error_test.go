/*
Copyright 2022 The KubeSphere Authors.

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

package stringutils

import (
	"github.com/kubesphere/ks-devops/pkg/server/errors"
	"testing"
)

func TestErrorOverride(t *testing.T) {
	type args struct {
		err    error
		format string
		a      []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{{
		name: "error is nil",
		args: args{
			err: nil,
		},
		wantErr: false,
	}, {
		name: "error is not nil",
		args: args{
			err:    errors.New("an error"),
			format: "msg",
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ErrorOverride(tt.args.err, tt.args.format, tt.args.a...); (err != nil) != tt.wantErr {
				t.Errorf("ErrorOverride() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
