package e2e

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
	"testing"
	"time"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
)

var (
	retryInterval        = time.Second * 10
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 2
	cleanupTimeout       = time.Second * 10
)

func TestPostgreSQL(t *testing.T) {
	postgresqlList := &v1alpha1.PostgresqlList{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, postgresqlList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	backupList := &v1alpha1.BackupList{}
	err = framework.AddToFrameworkScheme(apis.AddToScheme, backupList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	// run subtests
	t.Run("postgresql-group", func(t *testing.T) {
		t.Run("Cluster", PostgreSQLCluster)
		t.Run("Cluster2", PostgreSQLCluster)
	})
}

func PostgreSQLCluster(t *testing.T) {
	t.Parallel()
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()
	err := ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialize cluster resource: %v", err)
	}
	t.Log("Initialized cluster resource")
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}
	// get global framework variables
	f := framework.Global
	// wait for memcached-operator to be ready
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "postgresql-operator", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}
}
