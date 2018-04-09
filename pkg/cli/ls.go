package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"text/tabwriter"

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
			return doList(out, *opts)
		},
	}
	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "print the verbose details")
	return cmd
}

func doList(out io.Writer, opts clustersOpts) error {
	clusters, err := ioutil.ReadDir("clusters")
	if err != nil {
		return err
	}
	header := "Cluster Name:\t"
	if opts.verbose {
		header = fmt.Sprintf("%sLast modified:\tIs dir:\t", header)
	}

	w := tabwriter.NewWriter(out, 10, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintln(out, "Clusters currently being managed")
	fmt.Fprintln(w, header)
	for _, file := range clusters {
		line := fmt.Sprintf("%s\t", file.Name())
		// Might be worth adding if this is being tracked by the database?
		if opts.verbose {
			line = fmt.Sprintf("%s%s\t%s\t", line, file.ModTime().Format("2006-01-02 15:04:05"), strconv.FormatBool(file.IsDir()))
		}
		fmt.Fprintln(w, line)
	}
	w.Flush()
	return nil
}
