package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/apprenda/kismatic/pkg/store"

	"github.com/apprenda/kismatic/pkg/install"
	"github.com/apprenda/kismatic/pkg/plan"
	"github.com/apprenda/kismatic/pkg/util"
	"github.com/spf13/cobra"
)

// NewCmdPlan creates a new plan command
func NewCmdPlan(in io.Reader, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan",
		Short: "plan your Kubernetes cluster and generate a plan file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("Unexpected args: %v", args)
			}
			clusterStore, _ := CreateStoreIfNotExists(filepath.Join(assetsFolder, defaultDBName))
			defer clusterStore.Close()
			planner := install.FilePlanner{}

			plan, err := doPlan(in, out, planner, clusterStore)
			state := store.Planned
			if err != nil {
				return err
			}
			if plan.Provisioner.Provider == "" { // consider the cluster unmanaged IFF the user doesn't want a cloud provider, emulating the 1.0 workflow
				state = store.Unmanaged
			}
			spec := plan.ConvertToSpec(state)
			return clusterStore.Put(plan.Cluster.Name, spec)
		},
	}

	return cmd
}

func doPlan(in io.Reader, out io.Writer, planner install.FilePlanner, clusterStore store.ClusterStore) (*install.Plan, error) {

	providersDir := "./providers"
	fmt.Fprintln(out, "Plan your Kubernetes cluster:")
	name, err := util.PromptForAnyString(in, out, "Cluster name (must be unique)", defaultClusterName)
	if err != nil {
		return nil, fmt.Errorf("Error setting infrastructure provisioner: %v", err)
	}
	clusterFromStore, err := clusterStore.Get(name)
	if clusterFromStore != nil {
		return nil, fmt.Errorf("cluster with name %s already exists", name)
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving cluster information: %v", err)
	}
	exists, err := CheckClusterExists(name, clusterStore)
	if exists {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("cannot plan cluster with name %s, please use another name", name)
	}

	planner.File, _, _ = generateDirsFromName(name)
	var provider string
	availProviders, err := availableInfraProviders(providersDir)
	if err != nil {
		return nil, err
	}
	if len(availProviders) > 0 {
		provider, err = util.PromptForString(in, out, "Infrastructure provider (optional, leave blank if nodes are already provisioned)", "", availProviders)
		if err != nil {
			return nil, fmt.Errorf("Error setting infrastructure provider: %v", err)
		}
	}

	etcdNodes, err := util.PromptForInt(in, out, "Number of etcd nodes", 3)
	if err != nil {
		return nil, fmt.Errorf("Error reading number of etcd nodes: %v", err)
	}
	if etcdNodes <= 0 {
		return nil, fmt.Errorf("The number of etcd nodes must be greater than zero")
	}

	masterNodes, err := util.PromptForInt(in, out, "Number of master nodes", 2)
	if err != nil {
		return nil, fmt.Errorf("Error reading number of master nodes: %v", err)
	}
	if masterNodes <= 0 {
		return nil, fmt.Errorf("The number of master nodes must be greater than zero")
	}

	workerNodes, err := util.PromptForInt(in, out, "Number of worker nodes", 3)
	if err != nil {
		return nil, fmt.Errorf("Error reading number of worker nodes: %v", err)
	}
	if workerNodes <= 0 {
		return nil, fmt.Errorf("The number of worker nodes must be greater than zero")
	}

	ingressNodes, err := util.PromptForInt(in, out, "Number of ingress nodes (optional, set to 0 if not required)", 2)
	if err != nil {
		return nil, fmt.Errorf("Error reading number of ingress nodes: %v", err)
	}
	if ingressNodes < 0 {
		return nil, fmt.Errorf("The number of ingress nodes must be greater than or equal to zero")
	}

	storageNodes, err := util.PromptForInt(in, out, "Number of storage nodes (optional, set to 0 if not required)", 0)
	if err != nil {
		return nil, fmt.Errorf("Error reading number of storage nodes: %v", err)
	}
	if storageNodes < 0 {
		return nil, fmt.Errorf("The number of storage nodes must be greater than or equal to zero")
	}

	nfsVolumes, err := util.PromptForInt(in, out, "Number of existing NFS volumes to be attached", 0)
	if err != nil {
		return nil, fmt.Errorf("Error reading number of nfs volumes: %v", err)
	}
	if nfsVolumes < 0 {
		return nil, fmt.Errorf("The number of nfs volumes must be greater than or equal to zero")
	}

	fmt.Fprintln(out)
	fmt.Fprintf(out, "Generating installation plan file template with: \n")
	fmt.Fprintf(out, "- %s cluster name\n", name)
	if provider != "" {
		fmt.Fprintf(out, "- %s infrastructure provider\n", provider)
	} else {
		fmt.Fprintf(out, "- no infrastructure provider\n")
	}
	fmt.Fprintf(out, "- %d etcd nodes\n", etcdNodes)
	fmt.Fprintf(out, "- %d master nodes\n", masterNodes)
	fmt.Fprintf(out, "- %d worker nodes\n", workerNodes)
	fmt.Fprintf(out, "- %d ingress nodes\n", ingressNodes)
	fmt.Fprintf(out, "- %d storage nodes\n", storageNodes)
	fmt.Fprintf(out, "- %d nfs volumes\n", nfsVolumes)
	fmt.Fprintln(out)

	// If we are using KET to provision infrastructure, use the template file
	// defined by the infrastructure provider. Otherwise, generate the template
	// as we always have.
	if provider != "" {
		templater := plan.ProviderTemplatePlanner{ProvidersDir: providersDir}
		planTemplate, err := templater.GetPlanTemplate(provider)
		if err != nil {
			return nil, err
		}
		planTemplate.Cluster.Name = name
		planTemplate.Provisioner.Provider = provider
		planTemplate.Etcd.ExpectedCount = etcdNodes
		planTemplate.Master.ExpectedCount = masterNodes
		planTemplate.Worker.ExpectedCount = workerNodes
		planTemplate.Ingress.ExpectedCount = ingressNodes
		planTemplate.Storage.ExpectedCount = storageNodes

		return nil, planner.Write(planTemplate)
	}

	planTemplate := install.PlanTemplateOptions{
		ClusterName:               name,
		InfrastructureProvisioner: provider,
		EtcdNodes:                 etcdNodes,
		MasterNodes:               masterNodes,
		WorkerNodes:               workerNodes,
		IngressNodes:              ingressNodes,
		StorageNodes:              storageNodes,
		NFSVolumes:                nfsVolumes,
	}
	if err = install.WritePlanTemplate(planTemplate, &planner); err != nil {
		return nil, fmt.Errorf("error planning installation: %v", err)
	}

	fmt.Fprintf(out, "Wrote plan file template to %q\n", planner.File)
	fmt.Fprintf(out, "Edit the plan file to further describe your cluster. Once ready, execute the \"validate\" command to proceed.\n")
	plan, err := planner.Read()
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func availableInfraProviders(providersDir string) ([]string, error) {
	files, err := ioutil.ReadDir(providersDir)
	if err != nil {
		return nil, err
	}
	var p []string
	for _, f := range files {
		if f.IsDir() {
			p = append(p, f.Name())
		}
	}
	return p, nil
}
