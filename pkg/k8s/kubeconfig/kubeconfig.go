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

package kubeconfig

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipctl/pkg/k8s/client"
)

// Options is used for get kubeconfig from target cluster
type Options struct {
	*environment.AirshipCTLSettings
	ClientFactory client.Factory

	ClusterName string
	Namespace   string
	Output      string
}

// GetKubeConfig returns kubeconfig using defined in Options cluster name and namespace
func (o *Options) GetKubeConfig() ([]byte, error) {
	if o.Namespace == "" || o.ClusterName == "" {
		return nil, ErrInvalidOptions{}
	}

	err := o.AirshipCTLSettings.Config.EnsureComplete()
	if err != nil {
		return nil, err
	}

	c, err := o.ClientFactory(o.AirshipCTLSettings)
	if err != nil {
		return nil, err
	}

	coreClient := c.ClientSet().CoreV1()
	s, err := coreClient.Secrets(o.Namespace).Get(o.ClusterName+"-kubeconfig", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if s.Data == nil {
		return nil, ErrNoSecretData{s.Name}
	}

	v, found := s.Data["value"]
	if !found || (found && len(v) == 0) {
		return nil, ErrNoKubeConfig{s.Name}
	}

	return v, nil
}
