package e2e

import (
	goctx "context"
	"fmt"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql/v1alpha1"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

var (
	retryInterval        = time.Second * 10
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 2
	cleanupTimeout       = time.Second * 10
)

func TestDatabase(t *testing.T) {
	databaseList := &v1alpha1.DatabaseList{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, databaseList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	backupList := &v1alpha1.BackupList{}
	err = framework.AddToFrameworkScheme(apis.AddToScheme, backupList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	// run subtests
	t.Run("database-group", func(t *testing.T) {
		t.Run("Cluster", OperatorCluster)
		t.Run("Cluster2", OperatorCluster)
	})
}

func OperatorCluster(t *testing.T) {
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
	err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "postgresql-operator", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}

	if err = postgresalSQLTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}

func postgresalSQLTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}
	// create memcached custom resource
	db := &v1alpha1.Database{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "database",
			Namespace: namespace,
		},
		Spec: v1alpha1.DatabaseSpec{
			Size: 1,
		},
	}
	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), db, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}
	// wait for database to reach 1 replica
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "database", 1, retryInterval, timeout)
	if err != nil {
		return err
	}

	return nil
}
