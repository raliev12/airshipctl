/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package kubeconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipctl/pkg/k8s/client"
	"opendev.org/airship/airshipctl/pkg/k8s/client/fake"
	"opendev.org/airship/airshipctl/pkg/k8s/kubeconfig"
	"opendev.org/airship/airshipctl/testutil"
)

func TestGetKubeConfig(t *testing.T) {
	clusterName := "dummy_target_cluster"
	namespace := "default"
	secretData := []uint8("secret")

	tests := []struct {
		testName      string
		resAcc        fake.ResourceAccumulator
		expectedError error
	}{
		{
			testName: "getkubeconfig-no-err",
			resAcc: fake.WithTypedObjects(&coreV1.Secret{
				TypeMeta: metaV1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metaV1.ObjectMeta{
					Name:      clusterName + "-kubeconfig",
					Namespace: namespace,
				},
				Data: map[string][]uint8{
					"value": secretData,
				},
			}),
			expectedError: nil,
		},
		{
			testName: "getkubeconfig-no-data",
			resAcc: fake.WithTypedObjects(&coreV1.Secret{
				TypeMeta: metaV1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metaV1.ObjectMeta{
					Name:      clusterName + "-kubeconfig",
					Namespace: namespace,
				},
			}),
			expectedError: kubeconfig.ErrNoSecretData{S: clusterName + "-kubeconfig"},
		},
		{
			testName: "get-kubeconfig-no-valid-data",
			resAcc: fake.WithTypedObjects(&coreV1.Secret{
				TypeMeta: metaV1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metaV1.ObjectMeta{
					Name:      clusterName + "-kubeconfig",
					Namespace: namespace,
				},
				Data: map[string][]byte{},
			}),
			expectedError: kubeconfig.ErrNoKubeConfig{S: clusterName + "-kubeconfig"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.testName, func(t *testing.T) {
			opts := &kubeconfig.Options{
				AirshipCTLSettings: &environment.AirshipCTLSettings{Config: testutil.DummyConfig()},
				ClusterName:        clusterName,
				Namespace:          namespace,
				ClientFactory: func(_ *environment.AirshipCTLSettings) (client.Interface, error) {
					return fake.NewClient(tt.resAcc), nil
				},
			}
			actual, err := opts.GetKubeConfig()
			assert.Equal(t, tt.expectedError, err)
			if tt.expectedError == nil {
				assert.Equal(t, actual, secretData)
			}
		})
	}
}
