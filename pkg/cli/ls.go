package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"text/tabwriter"

	"github.com/apprenda/kismatic/pkg/store"

	"github.com/spf13/cobra"
)

type clustersOpts struct {
	verbose bool
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
	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "print the verbose details")
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
	header := "Cluster Name:\tCurrent State:\tDesired State:\t"
	if opts.verbose {
		header = fmt.Sprintf("%sLast modified:\tIs dir:\t", header)
	}

	w := tabwriter.NewWriter(out, 10, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintln(out, "Clusters currently being managed")
	fmt.Fprintln(w, header)

	for _, file := range clustersFromFSInfo {
		name := file.Name()
		var s1, s2 string
		if value, ok := clustersFromDB[name]; ok {
			s1 = value.Status.CurrentState
			s2 = value.Spec.DesiredState
		} else if opts.verbose {
			s1 = "not found"
			s2 = "not found"
		} else {
			continue
		}
		line := fmt.Sprintf("%s\t%s\t%s\t", name, s1, s2)
		if opts.verbose {
			line = fmt.Sprintf("%s%s\t%s\t", line, file.ModTime().Format("2006-01-02 15:04:05"), strconv.FormatBool(file.IsDir()))
		}
		fmt.Fprintln(w, line)
	}
	w.Flush()
	return nil
}
