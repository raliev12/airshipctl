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

package cluster

import (
	"fmt"

	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipctl/pkg/k8s/client"
	"opendev.org/airship/airshipctl/pkg/k8s/kubeconfig"
	"opendev.org/airship/airshipctl/pkg/k8s/utils"
)

const (
	kubeconfigLong = `
Retrieve cluster kubeconfig and save it to file or stdout.
`
	kubeconfigExample = `
# Retrieve target cluster kubeconfig and print it to stdout
airshipctl cluster kubeconfig target-cluster
`
)

// NewKubeConfigCommand creates a command which retrieves cluster kubeconfig.
func NewKubeConfigCommand(rootSettings *environment.AirshipCTLSettings, factory client.Factory) *cobra.Command {
	o := &kubeconfig.Options{AirshipCTLSettings: rootSettings}
	cmd := &cobra.Command{
		Use:     "kubeconfig [cluster_name]",
		Short:   "Retrieve kubeconfig from target cluster",
		Long:    kubeconfigLong[1:],
		Example: kubeconfigExample[1:],
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.ClientFactory = factory
			o.ClusterName = args[0]

			data, err := o.GetKubeConfig()
			if err != nil {
				return err
			}

			if o.Output != "" {
				err = utils.SaveKubeConfig(data, o.Output)
				if err != nil {
					return err
				}
			} else {
				_, err := fmt.Fprint(cmd.OutOrStdout(), string(data))
				return err
			}
			return nil
		},
	}

	initFlags(o, cmd)
	return cmd
}

func initFlags(o *kubeconfig.Options, cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVarP(
		&o.Namespace,
		"namespace",
		"n",
		"default",
		"namespace where cluster is located, if not specified default one will be used")

	flags.StringVarP(
		&o.Output,
		"file",
		"f",
		"",
		"local file to save kubeconfig data to, otherwise print to stdout")
}
