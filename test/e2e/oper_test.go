package e2e

import (
	goctx "context"
	"fmt"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/service"
	"github.com/dev4devs-com/postgresql-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

var (
	retryInterval        = time.Second * 30
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 4
	cleanupTimeout       = time.Second * 60
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

	t.Run("FullTest", FullTest)
}

func FullTest(t *testing.T) {
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

	// wait for postgresql-operator to be ready
	err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "postgresql-operator", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}

	// create database custom resource
	db := &v1alpha1.Database{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "database",
			Namespace: namespace,
		},
	}

	t.Log("Add database mandatory specs")
	utils.AddDatabaseMandatorySpecs(db)

	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), db, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("wait for database to reach 1 replica")
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "database", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}

	err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "database", Namespace: namespace}, db)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("wait for database status == OK")
	err = wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		cr, err := service.FetchDatabaseCR(db.Name, db.Namespace, f.Client.Client)
		if err != nil {
			return false, err
		}

		if cr.Status.DatabaseStatus == "OK" {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		t.Fatal(fmt.Errorf("could not get Database Status == OK: %v", err))
	}

	// create database custom resource
	bkp := &v1alpha1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "backup",
			Namespace: namespace,
		},
	}

	t.Log("Add bkp mandatory specs")
	utils.AddBackupMandatorySpecs(bkp)

	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), bkp, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("wait for backup status == OK")
	err = wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		cr, err := service.FetchBackupCR(bkp.Name, bkp.Namespace, f.Client.Client)
		if err != nil {
			return false, err
		}

		if cr.Status.BackupStatus == "OK" {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		t.Fatal(fmt.Errorf("could not get Backup Status == OK: %v", err))
	}
}
