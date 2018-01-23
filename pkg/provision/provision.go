package provision

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/apprenda/kismatic/pkg/install"
	yaml "gopkg.in/yaml.v2"
)

const providerDescriptorFilename = "provider.yaml"

type ProvisionOpts struct {
	AllowDestruction bool
}

// Provisioner is responsible for creating and destroying infrastructure for
// a given cluster.
type Provisioner interface {
	Provision(install.Plan, ProvisionOpts) (*install.Plan, error)
	Destroy(provider, clusterName string) error
}

// The SecretsGetter provides secrets required when interacting with cloud provider APIs.
type SecretsGetter interface {
	GetAsEnvironmentVariables(clusterName string, expectedEnvVars map[string]string) ([]string, error)
}

type provider struct {
	Version              string            `yaml:"string"`
	Description          string            `yaml:"description"`
	EnvironmentVariables map[string]string `yaml:"environmentVariables"`
}

type ketVars struct {
	KismaticVersion   string `json:"kismatic_version"`
	ClusterOwner      string `json:"cluster_owner"`
	PrivateSSHKeyPath string `json:"private_ssh_key_path"`
	PublicSSHKeyPath  string `json:"public_ssh_key_path"`
	ClusterName       string `json:"cluster_name"`
	MasterCount       int    `json:"master_count"`
	EtcdCount         int    `json:"etcd_count"`
	WorkerCount       int    `json:"worker_count"`
	IngressCount      int    `json:"ingress_count"`
	StorageCount      int    `json:"storage_count"`
}

// reads the descriptor for a specific provider found in the given directory
func readProviderDescriptor(providerDir string) (*provider, error) {
	var p provider
	providerDescriptorFile := filepath.Join(providerDir, providerDescriptorFilename)
	b, err := ioutil.ReadFile(providerDescriptorFile)
	if err != nil {
		return nil, fmt.Errorf("could not read provider descriptor: %v", err)
	}
	if err := yaml.Unmarshal(b, &p); err != nil {
		return nil, fmt.Errorf("could not unmarshal provider descriptor from %q: %v", providerDescriptorFile, err)
	}
	return &p, nil
}

func nodeGroupFromSlices(ips, internalIPs, hosts []string) install.NodeGroup {
	ng := install.NodeGroup{}
	ng.ExpectedCount = len(ips)
	ng.Nodes = []install.Node{}
	for i := range ips {
		n := install.Node{
			IP:   ips[i],
			Host: hosts[i],
		}
		if len(internalIPs) != 0 {
			n.InternalIP = internalIPs[i]
		}
		ng.Nodes = append(ng.Nodes, n)
	}
	return ng
}
