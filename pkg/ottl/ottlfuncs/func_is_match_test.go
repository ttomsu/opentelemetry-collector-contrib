// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ottlfuncs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
)

func Test_isMatch(t *testing.T) {
	tests := []struct {
		name     string
		target   ottl.Getter[interface{}]
		pattern  string
		expected bool
	}{
		{
			name: "replace match true",
			target: &ottl.StandardGetSetter[interface{}]{
				Getter: func(ctx context.Context, tCtx interface{}) (interface{}, error) {
					return "hello world", nil
				},
			},
			pattern:  "hello.*",
			expected: true,
		},
		{
			name: "replace match false",
			target: &ottl.StandardGetSetter[interface{}]{
				Getter: func(ctx context.Context, tCtx interface{}) (interface{}, error) {
					return "goodbye world", nil
				},
			},
			pattern:  "hello.*",
			expected: false,
		},
		{
			name: "replace match complex",
			target: &ottl.StandardGetSetter[interface{}]{
				Getter: func(ctx context.Context, tCtx interface{}) (interface{}, error) {
					return "-12.001", nil
				},
			},
			pattern:  "[-+]?\\d*\\.\\d+([eE][-+]?\\d+)?",
			expected: true,
		},
		{
			name: "target bool",
			target: &ottl.StandardGetSetter[interface{}]{
				Getter: func(ctx context.Context, tCtx interface{}) (interface{}, error) {
					return true, nil
				},
			},
			pattern:  "true",
			expected: true,
		},
		{
			name: "target int",
			target: &ottl.StandardGetSetter[interface{}]{
				Getter: func(ctx context.Context, tCtx interface{}) (interface{}, error) {
					return int64(1), nil
				},
			},
			pattern:  `\d`,
			expected: true,
		},
		{
			name: "target float",
			target: &ottl.StandardGetSetter[interface{}]{
				Getter: func(ctx context.Context, tCtx interface{}) (interface{}, error) {
					return 1.1, nil
				},
			},
			pattern:  `\d\.\d`,
			expected: true,
		},
		{
			name: "target pcommon.Value",
			target: &ottl.StandardGetSetter[interface{}]{
				Getter: func(ctx context.Context, tCtx interface{}) (interface{}, error) {
					v := pcommon.NewValueEmpty()
					v.SetStr("test")
					return v, nil
				},
			},
			pattern:  `test`,
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exprFunc, err := IsMatch(tt.target, tt.pattern)
			assert.NoError(t, err)
			result, err := exprFunc(context.Background(), nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_isMatch_validation(t *testing.T) {
	target := &ottl.StandardGetSetter[interface{}]{
		Getter: func(ctx context.Context, tCtx interface{}) (interface{}, error) {
			return "anything", nil
		},
	}
	_, err := IsMatch[interface{}](target, "\\K")
	require.Error(t, err)
}

func Test_isMatch_error(t *testing.T) {
	target := &ottl.StandardGetSetter[interface{}]{
		Getter: func(ctx context.Context, tCtx interface{}) (interface{}, error) {
			v := ottl.Path{}
			return v, nil
		},
	}
	exprFunc, err := IsMatch[interface{}](target, "test")
	assert.NoError(t, err)
	_, err = exprFunc(context.Background(), nil)
	require.Error(t, err)
}
