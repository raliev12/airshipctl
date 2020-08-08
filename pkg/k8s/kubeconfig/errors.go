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
	"fmt"
)

// ErrInvalidOptions returned when ClusterName or Namespace aren't properly set
type ErrInvalidOptions struct {
}

func (e ErrInvalidOptions) Error() string {
	return "Options are invalid, ClusterName and Namespace must be set"
}

// ErrNoSecretData returned when there is no data in secret
type ErrNoSecretData struct {
	S string
}

func (e ErrNoSecretData) Error() string {
	return fmt.Sprintf("No data in secret object %s", e.S)
}

// ErrNoKubeConfig returned when kubeconfig is not found in secret data or empty
type ErrNoKubeConfig struct {
	S string
}

func (e ErrNoKubeConfig) Error() string {
	return fmt.Sprintf("No kubeconfig in secret object %s", e.S)
}
