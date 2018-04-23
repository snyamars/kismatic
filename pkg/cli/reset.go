package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/apprenda/kismatic/pkg/install"
	"github.com/apprenda/kismatic/pkg/util"
	"github.com/spf13/cobra"
)

type resetOpts struct {
	planFilename       string
	generatedAssetsDir string
	verbose            bool
	outputFormat       string
	limit              []string
	force              bool
	removeAssets       bool
}

// NewCmdReset resets nodes
func NewCmdReset(in io.Reader, out io.Writer) *cobra.Command {
	opts := &resetOpts{}
	cmd := &cobra.Command{
		Use:   "reset CLUSTER_NAME",
		Short: "reset any changes made to the hosts by 'apply'",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Usage()
			}
			clusterName := args[0]
			path := filepath.Join(assetsFolder, defaultDBName)

			s, _ := CreateStoreIfNotExists(path)
			defer s.Close()
			if exists, err := CheckClusterExists(clusterName, s); !exists {
				return err
			}
			opts.planFilename, opts.generatedAssetsDir, _ = generateDirsFromName(clusterName)
			if opts.force == false {
				ans, err := util.PromptForString(in, out, "Are you sure you want to reset the cluster? All data will be lost", "N", []string{"N", "y"})
				if err != nil {
					return fmt.Errorf("error getting user response: %v", err)
				}
				if strings.ToLower(ans) != "y" {
					os.Exit(0)
				}
			}
			return doReset(out, opts)
		},
	}

	cmd.Flags().StringSliceVar(&opts.limit, "limit", []string{}, "comma-separated list of hostnames to limit the execution to a subset of nodes")
	cmd.Flags().BoolVar(&opts.verbose, "verbose", false, "enable verbose logging from the installation")
	cmd.Flags().StringVarP(&opts.outputFormat, "output", "o", "simple", "installation output format (options \"simple\"|\"raw\")")
	cmd.Flags().BoolVar(&opts.force, "force", false, `do not prompt`)
	cmd.Flags().BoolVar(&opts.removeAssets, "remove-assets", false, "remove generated-assets-dir")

	return cmd
}

func doReset(out io.Writer, opts *resetOpts) error {
	planner := &install.FilePlanner{File: opts.planFilename}
	if !planner.PlanExists() {
		return planFileNotFoundErr{filename: opts.planFilename}
	}
	plan, err := planner.Read()
	if err != nil {
		return fmt.Errorf("failed to read plan file: %v", err)
	}
	executorOpts := install.ExecutorOptions{
		GeneratedAssetsDirectory: opts.generatedAssetsDir,
		OutputFormat:             opts.outputFormat,
		Verbose:                  opts.verbose,
	}
	executor, err := install.NewExecutor(out, os.Stderr, executorOpts)
	if err != nil {
		return err
	}
	if err := executor.Reset(plan, opts.limit...); err != nil {
		return fmt.Errorf("error running reset: %v", err)
	}

	if opts.removeAssets {
		util.PrintHeader(out, "Removing Assets Directory", '=')
		if _, err := os.Stat(opts.generatedAssetsDir); os.IsNotExist(err) {
			util.PrettyPrintSkipped(out, "Removed %q", opts.generatedAssetsDir)
		} else {
			err := os.RemoveAll(opts.generatedAssetsDir)
			if err != nil {
				return fmt.Errorf("error deleting assets directory: %v", err)
			}
			util.PrettyPrintOk(out, "Remove %q directory", opts.generatedAssetsDir)
		}
	}

	return nil
}
