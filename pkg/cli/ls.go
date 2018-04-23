package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"text/tabwriter"

	"github.com/apprenda/kismatic/pkg/store"

	"github.com/spf13/cobra"
)

type clustersOpts struct {
}

// NewCmdList creates a new list command
func NewCmdList(out io.Writer) *cobra.Command {
	opts := &clustersOpts{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list the names of the clusters currently being managed",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := filepath.Join(assetsFolder, defaultDBName)
			s, _ := CreateStoreIfNotExists(path)
			defer s.Close()
			return doList(out, s, *opts)
		},
	}
	return cmd
}

func doList(out io.Writer, s store.ClusterStore, opts clustersOpts) error {
	clustersFromFSInfo, err := ioutil.ReadDir(assetsFolder)
	if err != nil {
		return err
	}
	clustersFromDB, err := s.GetAll()
	if err != nil {
		return err
	}
	header := "Cluster Name\tCurrent State\tDesired State\tLast modified\t"

	w := tabwriter.NewWriter(out, 10, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, header)

	for _, file := range clustersFromFSInfo {
		name := file.Name()
		// do not print db file
		if name == defaultDBName {
			continue
		}
		current := "not found"
		desired := "not found"
		if value, ok := clustersFromDB[name]; ok {
			current = value.Status.CurrentState
			desired = value.Spec.DesiredState
		}
		line := fmt.Sprintf("%s\t%s\t%s\t%s\t", name, current, desired, file.ModTime().Format("2006-01-02 15:04:05"))
		fmt.Fprintln(w, line)
	}
	w.Flush()
	return nil
}
