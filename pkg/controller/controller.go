package controller

import (
	"context"
	"io"
	"log"
	"path/filepath"
	"time"

	"github.com/apprenda/kismatic/pkg/install"
	"github.com/apprenda/kismatic/pkg/provision"
	"github.com/apprenda/kismatic/pkg/store"
)

// The ClusterController manages the lifecycle of clusters
type ClusterController interface {
	Run(ctx context.Context)
}

// ExecutorCreator creates executors that can be used for executing actions
// against a specific cluster.
type ExecutorCreator func(clusterName, clusterAssetsDir string, logFile io.Writer) (install.Executor, error)

// ProvisionerCreator creates provisioners that can be used for standing up
// infrastructure for a specific cluster.
type ProvisionerCreator func(output io.Writer) provision.Provisioner

// AssetsDir is the location where the controller will store all file-based
// assets that are generated throughout the management of all clusters.
type AssetsDir string

// ForCluster returns the directory that holds assets for the given cluster
func (ad AssetsDir) ForCluster(clusterName string) string {
	return filepath.Join(string(ad), clusterName)
}

// New returns a cluster controller
func New(
	logger *log.Logger,
	planner Planner,
	execCreator ExecutorCreator,
	provisionerCreator ProvisionerCreator,
	cs store.ClusterStore,
	reconFreq time.Duration,
	assetsDir AssetsDir) ClusterController {
	return &multiClusterController{
		assetsDir:          assetsDir,
		log:                logger,
		planner:            planner,
		newExecutor:        execCreator,
		clusterStore:       cs,
		reconcileFreq:      reconFreq,
		clusterControllers: make(map[string]chan<- struct{}),
		provisionerCreator: provisionerCreator,
	}
}

// DefaultExecutorCreator creates an executor that can be used to run operations
// against a single cluster.
func DefaultExecutorCreator() ExecutorCreator {
	return func(clusterName string, clusterAssetsDir string, logFile io.Writer) (install.Executor, error) {
		executorOpts := install.ExecutorOptions{
			GeneratedAssetsDirectory: filepath.Join(clusterAssetsDir, "assets"),
			RunsDirectory:            filepath.Join(clusterAssetsDir, "runs"),
			OutputFormat:             "simple",
			Verbose:                  true,
		}
		executor, err := install.NewExecutor(logFile, logFile, executorOpts)
		if err != nil {
			return nil, err
		}
		return executor, nil
	}
}
