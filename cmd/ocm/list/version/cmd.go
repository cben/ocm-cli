/*
Copyright (c) 2020 Red Hat, Inc.

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

package version

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/openshift-online/ocm-cli/pkg/ocm"
	"github.com/openshift-online/ocm-cli/pkg/provider"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "versions",
	Aliases: []string{"version"},
	Short:   "List available versions",
	Long:    "List the versions available for provisioning a cluster",
	Example: `  # List all supported cluster versions
  ocm list versions`,
	Args: cobra.NoArgs,
	RunE: run,
}

func run(cmd *cobra.Command, argv []string) error {
	// Create the client for the OCM API:
	connection, err := ocm.NewConnection().Build()
	if err != nil {
		return fmt.Errorf("Failed to create OCM connection: %v", err)
	}
	defer connection.Close()

	client := connection.ClustersMgmt().V1()
	versions, defaultVersion, err := provider.GetEnabledVersions(client)
	if err != nil {
		return fmt.Errorf("Can't retrieve versions: %v", err)
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprint(writer, "DEFAULT?\tVERSION\n")
	for _, v := range versions {
		isDefault := " "
		if v == defaultVersion {
			isDefault = "default"
		}
		fmt.Fprintf(writer, "%s\t%s\n", isDefault, v)
	}
	writer.Flush()

	return nil
}
