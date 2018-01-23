package cli

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/apprenda/kismatic/pkg/install"
	"github.com/apprenda/kismatic/pkg/provision"

	"github.com/spf13/cobra"
)

// NewCmdProvision creates a new provision command
func NewCmdProvision(in io.Reader, out io.Writer, opts *installOpts) *cobra.Command {
	provisionOpts := provision.ProvisionOpts{}
	cmd := &cobra.Command{
		Use:   "provision",
		Short: "provision your Kubernetes cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			fp := &install.FilePlanner{File: opts.planFilename}
			plan, err := fp.Read()
			if err != nil {
				return fmt.Errorf("unable to read plan file: %v", err)
			}
			path, err := os.Getwd()
			if err != nil {
				return err
			}
			user, err := user.Current()
			if err != nil {
				return err
			}

			tf := provision.AnyTerraform{
				ClusterOwner:    user.Username,
				Output:          out,
				BinaryPath:      filepath.Join(path, "terraform"),
				KismaticVersion: install.KismaticVersion.String(),
				ProvidersDir:    filepath.Join(path, "providers"),
				StateDir:        filepath.Join(path, assetsFolder),
				SecretsGetter:   environmentSecretsGetter{},
			}

			updatedPlan, err := tf.Provision(*plan, provisionOpts)
			if err != nil {
				return err
			}
			if err := fp.Write(updatedPlan); err != nil {
				return fmt.Errorf("error writing updated plan file to %s: %v", opts.planFilename, err)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&provisionOpts.AllowDestruction, "allow-destruction", false, "Allows possible infrastructure destruction through provisioner planning, required if mutation is scaling down (Use with care)")
	return cmd
}

// NewCmdDestroy creates a new destroy command
func NewCmdDestroy(in io.Reader, out io.Writer, opts *installOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "destroy your provisioned cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			fp := &install.FilePlanner{File: opts.planFilename}
			plan, err := fp.Read()
			if err != nil {
				return fmt.Errorf("unable to read plan file: %v", err)
			}
			path, err := os.Getwd()
			if err != nil {
				return err
			}
			tf := provision.AnyTerraform{
				Output:          out,
				BinaryPath:      filepath.Join(path, "./terraform"),
				KismaticVersion: install.KismaticVersion.String(),
				ProvidersDir:    filepath.Join(path, "providers"),
				StateDir:        filepath.Join(path, assetsFolder),
				SecretsGetter:   environmentSecretsGetter{},
			}

			return tf.Destroy(plan.Provisioner.Provider, plan.Cluster.Name)
		},
	}
	return cmd
}

type environmentSecretsGetter struct{}

// GetAsEnvironmentVariables returns a slice of the expected environment
// variables sourcing them from the current process' environment.
func (environmentSecretsGetter) GetAsEnvironmentVariables(clusterName string, expected map[string]string) ([]string, error) {
	var vars []string
	var missingVars []string
	for _, expectedEnvVar := range expected {
		val := os.Getenv(expectedEnvVar)
		if val == "" {
			missingVars = append(missingVars, expectedEnvVar)
		}
		vars = append(vars, fmt.Sprintf("%s=%s", expectedEnvVar, val))
	}
	if len(missingVars) > 0 {
		return nil, fmt.Errorf("%v", missingVars)
	}
	return vars, nil
}
