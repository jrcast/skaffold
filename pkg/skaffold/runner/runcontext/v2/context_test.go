/*
Copyright 2019 The Skaffold Authors

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

package v2

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/config"
	latestV1 "github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest/v1"
	latestV2 "github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest/v2"
	"github.com/GoogleContainerTools/skaffold/testutil"
)

func TestRunContext_UpdateNamespaces(t *testing.T) {
	tests := []struct {
		description   string
		oldNamespaces []string
		newNamespaces []string
		expected      []string
	}{
		{
			description:   "update namespace when not present in runContext",
			oldNamespaces: []string{"test"},
			newNamespaces: []string{"another"},
			expected:      []string{"another", "test"},
		},
		{
			description:   "update namespace with duplicates should not return duplicate",
			oldNamespaces: []string{"test", "foo"},
			newNamespaces: []string{"another", "foo", "another"},
			expected:      []string{"another", "foo", "test"},
		},
		{
			description:   "update namespaces when namespaces is empty",
			oldNamespaces: []string{"test", "foo"},
			newNamespaces: []string{},
			expected:      []string{"test", "foo"},
		},
		{
			description:   "update namespaces when runcontext namespaces is empty",
			oldNamespaces: []string{},
			newNamespaces: []string{"test", "another"},
			expected:      []string{"another", "test"},
		},
		{
			description:   "update namespaces when both namespaces and runcontext namespaces is empty",
			oldNamespaces: []string{},
			newNamespaces: []string{},
			expected:      []string{},
		},
		{
			description:   "update namespace when runcontext namespace has an empty string",
			oldNamespaces: []string{""},
			newNamespaces: []string{"another"},
			expected:      []string{"another"},
		},
		{
			description:   "update namespace when namespace is empty string",
			oldNamespaces: []string{"test"},
			newNamespaces: []string{""},
			expected:      []string{"test"},
		},
		{
			description:   "update namespace when namespace is empty string and runContext is empty",
			oldNamespaces: []string{},
			newNamespaces: []string{""},
		},
	}
	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			runCtx := &RunContext{
				Namespaces: test.oldNamespaces,
			}

			runCtx.UpdateNamespaces(test.newNamespaces)

			t.CheckDeepEqual(test.expected, runCtx.Namespaces)
		})
	}
}

func TestGetRunContextDefaultWorkdir(t *testing.T) {
	testutil.Run(t, "default workdir", func(t *testutil.T) {
		rctx, err := GetRunContext(config.SkaffoldOptions{}, []*latestV2.SkaffoldConfig{})
		pwd, _ := os.Getwd()
		t.CheckDeepEqual(pwd, rctx.WorkingDir)
		t.CheckNoError(err)
	})
}

func TestGetRunContextCustomWorkdir(t *testing.T) {
	testutil.Run(t, "default workdir", func(t *testutil.T) {
		tmpDir := t.NewTempDir()
		tmpDir.Write("skaffold.yaml", fmt.Sprintf("apiVersion: %s\nkind: Config", latestV1.Version)).
			Chdir()
		rctx, err := GetRunContext(config.SkaffoldOptions{
			ConfigurationFile: filepath.Join(tmpDir.Root(), "skaffold.yaml"),
		}, []*latestV2.SkaffoldConfig{})
		t.CheckDeepEqual(tmpDir.Root(), rctx.WorkingDir)
		t.CheckNoError(err)
	})
}