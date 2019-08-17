package backup

import (
	"context"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"
)

var log = logf.Log.WithName("controller_backup")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Backup Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileBackup{client: mgr.GetClient(), scheme: mgr.GetScheme(), config: mgr.GetConfig()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("backup-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Backup
	err = c.Watch(&source.Kind{Type: &v1alpha1.Backup{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	if err := watchCronJob(c); err != nil {
		return err
	}

	if err := watchSecret(c); err != nil {
		return err
	}

	if err := watchPod(c); err != nil {
		return err
	}

	if err := watchService(c); err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileBackup implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileBackup{}

// ReconcileBackup reconciles a Backup object
type ReconcileBackup struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client    client.Client
	config    *rest.Config
	scheme    *runtime.Scheme
	dbPod     *v1.Pod
	dbService *v1.Service
}

// Reconcile reads that state of the cluster for a Backup object and makes changes based on the state read
// and what is in the Backup.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileBackup) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Backup ...")

	bkp, err := r.fetchBackupCR(request)
	if err != nil {
		reqLogger.Error(err, "Failed to get the Backup Custom Resource. Check if the Backup CR is applied in the cluster")
		return reconcile.Result{}, err
	}

	// Add const values for mandatory specs
	addMandatorySpecsDefinitions(bkp)

	// Create mandatory objects for the Backup
	if err := r.createSecondaryResources(bkp, request); err != nil {
		reqLogger.Error(err, "Failed to create and update the secondary resources required for the Backup CR")
		return reconcile.Result{}, err
	}

	// Update the CR status for the primary resource
	if err := r.createUpdateCRStatus(request); err != nil {
		reqLogger.Error(err, "Failed to create and update the status in the Backup CR")
		return reconcile.Result{}, err
	}

	// stop reconciliation
	return reconcile.Result{}, nil
}

//createUpdateCRStatus will create and update the status in the CR applied in the cluster
func (r *ReconcileBackup) createUpdateCRStatus(request reconcile.Request) error {
	if err := r.updatePodDatabaseFoundStatus(request, r.dbPod); err != nil {
		return err
	}

	if err := r.updateServiceDatabaseFoundStatus(request); err != nil {
		return err
	}

	if err := r.updateCronJobStatus(request); err != nil {
		return err
	}

	if err := r.updateDBSecretStatus(request); err != nil {
		return err
	}

	if err := r.updateAWSSecretStatus(request); err != nil {
		return err
	}

	if err := r.updateEncSecretStatus(request); err != nil {
		return err
	}

	// Update status for Backup
	if err := r.updateBackupStatus(request); err != nil {
		return err
	}
	return nil
}

//createSecondaryResources will create and update the secondary resources which are required in order to make works successfully the primary resource(CR)
func (r *ReconcileBackup) createSecondaryResources(bkp *v1alpha1.Backup, request reconcile.Request) error {
	// Check if the database instance was created
	db, err := r.fetchPostgreSQLInstance(bkp)
	if err != nil {
		return err
	}

	// Get database pod
	dbPod, err := r.fetchPostgreSQLPod(bkp, db)
	if err != nil || dbPod == nil {
		time.Sleep(2 * time.Second)
		return err
	}

	// set in the reconcile
	r.dbPod = dbPod

	// Get database service
	dbService, err := r.fetchPostgreSQLService(bkp, db)
	if err != nil || dbService == nil {
		return err
	}

	// set in the reconcile
	r.dbService = dbService

	// Check if the secret for the database is created, if not create one
	if _, err := r.fetchSecret(bkp.Namespace, dbSecretPrefix+bkp.Name); err != nil {
		secretData, err := r.buildDBSecretData(bkp, db)
		if err != nil {
			return err
		}
		dbSecret := buildSecret(bkp, dbSecretPrefix, secretData, nil, r.scheme)
		if err := r.client.Create(context.TODO(), dbSecret); err != nil {
			return err
		}
	}

	// Check if the secret for the s3 is created, if not create one
	if _, err := r.fetchSecret(getAwsSecretNamespace(bkp), getAWSSecretName(bkp)); err != nil {
		if bkp.Spec.AwsCredentialsSecretName != "" {
			return err
		}
		secretData := buildAwsSecretData(bkp)
		awsSecret := buildSecret(bkp, awsSecretPrefix, secretData, nil, r.scheme)
		if err := r.client.Create(context.TODO(), awsSecret); err != nil {
			return err
		}
	}

	// Check if the secret for the encryptionKey is created, if not create one just when the data is informed
	if hasEncryptionKeySecret(bkp) {
		if _, err := r.fetchSecret(getEncSecretNamespace(bkp), getEncSecretName(bkp)); err != nil {
			if bkp.Spec.EncryptionKeySecretName != "" {
				return err
			}
			secretData, secretStringData := buildEncSecretData(bkp)
			encSecret := buildSecret(bkp, encSecretPrefix, secretData, secretStringData, r.scheme)
			if err := r.client.Create(context.TODO(), encSecret); err != nil {
				return err
			}
		}
	}

	// Check if the cronJob is created, if not create one
	if _, err := r.fetchCronJob(bkp); err != nil {
		if err := r.client.Create(context.TODO(), buildCronJob(bkp, r.scheme)); err != nil {
			return err
		}
	}
	return nil
}
