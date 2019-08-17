package e2e

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/operator-framework/operator-sdk/test/test-framework/pkg/apis"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 300
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestPostgreSQL(t *testing.T) {
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()
	f := framework.Global

	postgresqlList := &v1alpha1.PostgresqlList{}
	if err := framework.AddToFrameworkScheme(apis.AddToScheme, postgresqlList); err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	backupList := &v1alpha1.BackupList{}
	if err := framework.AddToFrameworkScheme(apis.AddToScheme, backupList); err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatalf("failed to get namespace: %v", err)
	}

	postgresql := &v1alpha1.Postgresql{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-postgresql",
			Namespace: namespace,
		},
	}

	err = f.Client.Create(context.TODO(), postgresql, &framework.CleanupOptions{
		TestContext:   ctx,
		Timeout:       timeout,
		RetryInterval: retryInterval,
	})

	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// wait for postgresql-operator to be ready
	err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "postgresql-operator", 1, time.Second*5, time.Second*30)
	if err != nil {
		t.Fatal(err)
	}

}
