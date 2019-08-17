package e2e

import (
	goctx "context"
	"fmt"
	apis "github.com/dev4devs-com/postgresql-operator/pkg/apis"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	"testing"
	"time"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestPostgreSQL(t *testing.T) {
	postgresqlList := &v1alpha1.PostgresqlList{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, postgresqlList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("postgresql-group", func(t *testing.T) {
		t.Run("Cluster", PostgreSQLCluster)
		t.Run("Cluster2", PostgreSQLCluster)
	})
}

func postgresalScaleTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}
	// create memcached custom resource
	examplePostgresql := &v1alpha1.Postgresql{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-postgresql",
			Namespace: namespace,
		},
		Spec: v1alpha1.PostgresqlSpec{
			Size: 1,
		},
	}
	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), examplePostgresql, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}
	// wait for example-memcached to reach 3 replicas
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-postgresql", 3, retryInterval, timeout)
	if err != nil {
		return err
	}

	err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "example-postgresql", Namespace: namespace}, examplePostgresql)
	if err != nil {
		return err
	}
	examplePostgresql.Spec.Size = 2
	err = f.Client.Update(goctx.TODO(), examplePostgresql)
	if err != nil {
		return err
	}

	// wait for example-memcached to reach 4 replicas
	return e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-postgresql", 2, retryInterval, timeout)
}

func PostgreSQLCluster(t *testing.T) {
	t.Parallel()
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()
	err := ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}
	t.Log("Initialized cluster resources")
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

	if err = postgresalScaleTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}
