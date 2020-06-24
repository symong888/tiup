// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"errors"
	"io/ioutil"
	"os"

	perrs "github.com/pingcap/errors"
	"github.com/pingcap/tiup/pkg/cliutil"
	"github.com/pingcap/tiup/pkg/cluster/spec"
	"github.com/pingcap/tiup/pkg/meta"
	tiuputils "github.com/pingcap/tiup/pkg/utils"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listCluster()
		},
	}
	return cmd
}

func listCluster() error {
	clusterDir := spec.ProfilePath(spec.TiOpsClusterDir)
	clusterTable := [][]string{
		// Header
		{"Name", "User", "Version", "Path", "PrivateKey"},
	}
	fileInfos, err := ioutil.ReadDir(clusterDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, fi := range fileInfos {
		if tiuputils.IsNotExist(spec.ClusterPath(fi.Name(), spec.MetaFileName)) {
			continue
		}
		metadata, err := spec.ClusterMetadata(fi.Name())
		if err != nil && !errors.Is(perrs.Cause(err), meta.ErrValidate) {
			return perrs.Trace(err)
		}

		clusterTable = append(clusterTable, []string{
			fi.Name(),
			metadata.User,
			metadata.Version,
			spec.ClusterPath(fi.Name()),
			spec.ClusterPath(fi.Name(), "ssh", "id_rsa"),
		})
	}

	cliutil.PrintTable(clusterTable, true)
	return nil
}
