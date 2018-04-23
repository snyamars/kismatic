package cli

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/apprenda/kismatic/pkg/install"
	"github.com/apprenda/kismatic/pkg/provision"
	"github.com/apprenda/kismatic/pkg/store"

	"github.com/spf13/cobra"
)

// NewCmdProvision creates a new provision command
func NewCmdProvision(in io.Reader, out io.Writer) *cobra.Command {
	provisionOpts := provision.ProvisionOpts{}
	cmd := &cobra.Command{
		Use:   "provision CLUSTER_NAME",
		Short: "provision your Kubernetes cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Usage()
			}
			clusterName := args[0]
			dbPath := filepath.Join(assetsFolder, defaultDBName)
			s, _ := CreateStoreIfNotExists(dbPath)
			defer s.Close()
			if exists, err := CheckClusterExists(clusterName, s); !exists {
				return err
			}
			planPath, _, _ := generateDirsFromName(clusterName)
			fp := &install.FilePlanner{File: planPath}
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

			updatedPlan, provErr := tf.Provision(*plan, provisionOpts)
			spec := updatedPlan.ConvertToSpec(store.Provisioned)
			if provErr != nil {
				spec.Status.CurrentState = store.ProvisionFailed
				if err := s.Put(updatedPlan.Cluster.Name, spec); err != nil {
					return fmt.Errorf("%v: %v", provErr, err)
				}
				return provErr
			}
			if err := fp.Write(updatedPlan); err != nil {
				return fmt.Errorf("error writing updated plan file to %s: %v", planPath, err)
			}
			return s.Put(updatedPlan.Cluster.Name, spec)
		},
	}
	cmd.Flags().BoolVar(&provisionOpts.AllowDestruction, "allow-destruction", false, "Allows possible infrastructure destruction through provisioner planning, required if mutation is scaling down (Use with care)")
	return cmd
}

// NewCmdDestroy creates a new destroy command
func NewCmdDestroy(in io.Reader, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy CLUSTER_NAME",
		Short: "destroy your provisioned cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Usage()
			}
			clusterName := args[0]
			dbPath := filepath.Join(assetsFolder, defaultDBName)
			s, _ := CreateStoreIfNotExists(dbPath)
			defer s.Close()
			if exists, err := CheckClusterExists(clusterName, s); !exists {
				return err
			}
			planPath, _, _ := generateDirsFromName(clusterName)
			fp := &install.FilePlanner{File: planPath}
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
			if err := tf.Destroy(plan.Provisioner.Provider, plan.Cluster.Name); err != nil {
				return err
			}
			// No point in committing a state update to the database, just remove the entry.
			return doRemove(out, clusterName, s)
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
