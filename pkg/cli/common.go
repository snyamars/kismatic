package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/apprenda/kismatic/pkg/store"
)

const (
	defaultDBName        = "clusterStates.db"
	defaultPlanName      = "kismatic-cluster.yaml"
	defaultClusterName   = "kubernetes"
	assetsFolder         = "clusters"
	defaultGeneratedName = "generated"
	defaultRunsName      = "runs"

	playbooksPath  = "playbooks"
	ansiblePath    = "ansible"
	clustersBucket = "kismatic"
)

type planFileNotFoundErr struct {
	filename string
}

func (e planFileNotFoundErr) Error() string {
	return fmt.Sprintf("Plan file not found at %q. If you don't have a plan file, you may generate one with 'kismatic plan'", e.filename)
}

// Returns a path to a plan file, generated dir, and runs dir according to the clusterName
func generateDirsFromName(clusterName string) (string, string, string) {
	return filepath.Join(assetsFolder, clusterName, defaultPlanName), filepath.Join(assetsFolder, clusterName, defaultGeneratedName), filepath.Join(assetsFolder, clusterName, defaultRunsName)
}

// CheckClusterExists does a simple check to see if the cluster folder+plan file exists in clusters
// returns true even in the cases where an error exists if the scan hasn't completed.
func CheckClusterExists(name string, s store.ClusterStore) (bool, error) {
	if err := os.MkdirAll(assetsFolder, 0700); err != nil {
		return true, err
	}
	if spec, err := s.Get(name); spec != nil {
		if err != nil {
			return true, err
		}
		return true, nil
	}
	files, err := ioutil.ReadDir(assetsFolder)
	if err != nil {
		return true, err
	}
	for _, finfo := range files {
		if finfo.Name() == name {
			possiblePlans, err := ioutil.ReadDir(filepath.Join(assetsFolder, finfo.Name()))
			if err != nil {
				return true, err
			}
			for _, possiblePlan := range possiblePlans {
				if possiblePlan.Name() == defaultPlanName {
					return true, nil
				}
			}
		}
	}
	return false, fmt.Errorf("cluster with name %q not found. If you have a plan file, but your cluster doesn't exist, please run kismatic import PLAN_FILE_PATH GENERATED_ASSETS_DIR", name)
}

// CreateStoreIfNotExists creates a database file at location path. returns ClusterStore that will interact with that file, and a logger for the store.
func CreateStoreIfNotExists(path string) (store.ClusterStore, *log.Logger) {
	parent, _ := filepath.Split(path)
	logger := log.New(os.Stdout, "[kismatic] ", log.LstdFlags|log.Lshortfile)
	if err := os.MkdirAll(parent, 0700); err != nil {
		logger.Fatalf("Error creating store directory structure: %v", err)
	}
	// Create the store
	s, err := store.New(path, 0600, logger)
	if err != nil {
		logger.Fatalf("Error creating store: %v", err)
	}

	err = s.CreateBucket(clustersBucket)
	if err != nil {
		logger.Fatalf("Error creating bucket in store: %v", err)
	}

	clusterStore := store.NewClusterStore(s, clustersBucket)
	return clusterStore, logger
}
