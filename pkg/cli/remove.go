package cli

import (
	"io"
	"os"
	"path/filepath"

	"github.com/apprenda/kismatic/pkg/store"

	"github.com/spf13/cobra"
)

// NewCmdRemove removes all refs to a cluster from the database and filesystem
func NewCmdRemove(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove CLUSTER_NAME",
		Short:   "removes all references to a cluster from the database and filesystem",
		Aliases: []string{"rm"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Usage()
			}
			name := args[0]
			dbPath := filepath.Join(assetsFolder, defaultDBName)
			s, _ := CreateStoreIfNotExists(dbPath)
			defer s.Close()
			exists, err := CheckClusterExists(name, s)
			if !exists || err != nil {
				return err
			}
			return doRemove(out, name, s)
		},
	}
	return cmd
}

func doRemove(out io.Writer, name string, c store.ClusterStore) error {
	if err := c.Delete(name); err != nil {
		return err
	}
	planFile, _, _ := generateDirsFromName(name)
	dir, _ := filepath.Split(planFile)
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return nil
}
