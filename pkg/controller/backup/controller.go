package backup

import (
	"context"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	"github.com/go-logr/logr"
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

const (
	cronJob   = "cronJob"
	dbSecret  = "dbSecret"
	swsSecret = "swsSecret"
	encSecret = "encSecret"
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

	// Watch cronJob
	if err := watchCronJob(c); err != nil {
		return err
	}

	// Watch watchSecret
	if err := watchSecret(c); err != nil {
		return err
	}

	// Watch Pod
	if err := watchPod(c); err != nil {
		return err
	}

	// Watch Service
	if err := watchService(c); err != nil {
		return err
	}

	return nil
}

// Update the object and reconcile it
func (r *ReconcileBackup) update(obj runtime.Object, reqLogger logr.Logger) error {
	err := r.client.Update(context.TODO(), obj)
	if err != nil {
		reqLogger.Error(err, "Failed to update Object", "obj:", obj)
		return err
	}
	reqLogger.Info("Object updated", "obj:", obj)
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

// Create the object and reconcile it
func (r *ReconcileBackup) create(bkp *v1alpha1.Backup, db *v1alpha1.Postgresql, kind string, reqLogger logr.Logger) error {
	obj, err := r.buildFactory(bkp,db, kind, reqLogger)
	if err != nil {
		reqLogger.Error(err, "Failed to build object ", "kind", kind, "Namespace", bkp.Namespace)
		return err
	}
	reqLogger.Info("Creating a new ", "kind", kind, "Namespace", bkp.Namespace)
	err = r.client.Create(context.TODO(), obj)
	if err != nil {
		reqLogger.Error(err, "Failed to create new ", "kind", kind, "Namespace", bkp.Namespace)
		return err
	}
	reqLogger.Info("Created successfully", "kind", kind, "Namespace", bkp.Namespace)
	return nil
}

// buildFactory will return the resource according to the kind defined
func (r *ReconcileBackup) buildFactory(bkp *v1alpha1.Backup, db *v1alpha1.Postgresql, kind string, reqLogger logr.Logger) (runtime.Object, error) {
	reqLogger.Info("Check "+kind, "into the namespace", bkp.Namespace)
	switch kind {
	case cronJob:
		return r.buildCronJob(bkp), nil
	case dbSecret:
		// build Database secret data
		secretData, err := r.buildDBSecretData(bkp, db)
		if err != nil {
			reqLogger.Error(err, "Unable to create DB Data secret")
			return nil, err
		}
		return r.buildSecret(bkp, dbSecretPrefix, secretData, nil), nil
	case swsSecret:
		secretData := buildAwsSecretData(bkp)
		return r.buildSecret(bkp, awsSecretPrefix, secretData, nil), nil
	case encSecret:
		secretData, secretStringData := buildEncSecretData(bkp)
		return r.buildSecret(bkp, encSecretPrefix, secretData, secretStringData), nil
	default:
		msg := "Failed to recognize type of object" + kind + " into the Namespace " + bkp.Namespace
		panic(msg)
	}
}

// Reconcile reads that state of the cluster for a Backup object and makes changes based on the state read
// and what is in the Backup.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileBackup) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Backup")

	// Fetch the PostgreSQL DB Backup
	bkp := &v1alpha1.Backup{}
	bkp, err := r.fetchBkpInstance(reqLogger, request)
	if err != nil {
		reqLogger.Error(err, "Failed to get PostgreSQL Backup")
		return reconcile.Result{}, err
	}

	// Add const values for mandatory specs
	addMandatorySpecsDefinitions(bkp)

	// Check if the database instance was created
	db, err := r.fetchDBInstance( bkp,reqLogger);
	if err != nil {
		return reconcile.Result{}, err
	}

	// Get database pod
	dbPod, err := r.fetchBDPod(bkp, db, reqLogger)
	if err != nil || dbPod == nil {
		reqLogger.Error(err, "Unable to find the database pod", "request.namespace", request.Namespace)
		return reconcile.Result{RequeueAfter: time.Second * 10}, err
	}

	// set in the reconcile
	r.dbPod = dbPod

	// Get database service
	dbService, err := r.fetchServiceDB(bkp, db, reqLogger)
	if err != nil || dbService == nil {
		reqLogger.Error(err, "Unable to find the database service", "request.namespace", request.Namespace)
		return reconcile.Result{RequeueAfter: time.Second * 10}, err
	}

	// set in the reconcile
	r.dbService = dbService

	// Check if the secret for the database is created, if not create one
	if _, err := r.fetchSecret(reqLogger, bkp.Namespace, dbSecretPrefix+bkp.Name); err != nil {
		if err := r.create(bkp, db, dbSecret, reqLogger); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Check if the secret for the s3 is created, if not create one
	if _, err := r.fetchSecret(reqLogger, getAwsSecretNamespace(bkp), getAWSSecretName(bkp)); err != nil {
		if bkp.Spec.AwsCredentialsSecretName != "" {
			reqLogger.Error(err, "Unable to find AWS secret informed and will not be created by the operator", "SecretName", bkp.Spec.AwsCredentialsSecretName)
			return reconcile.Result{}, err
		}
		if err := r.create(bkp, db, swsSecret, reqLogger); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Check if the secret for the encryptionKey is created, if not create one just when the data is informed
	if hasEncryptionKeySecret(bkp) {
		if _, err := r.fetchSecret(reqLogger, getEncSecretNamespace(bkp), getEncSecretName(bkp)); err != nil {
			if bkp.Spec.EncryptionKeySecretName != "" {
				reqLogger.Error(err, "Unable to find EncryptionKey secret informed and will not be created by the operator", "SecretName", bkp.Spec.EncryptionKeySecretName)
				return reconcile.Result{}, err
			}
			if err := r.create(bkp, db, encSecret, reqLogger); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	// Check if the cronJob is created, if not create one
	if _, err := r.fetchCronJob(reqLogger, bkp); err != nil {
		if err := r.create(bkp, db, cronJob, reqLogger); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Update status for pod database found
	if err := r.updatePodDatabaseFoundStatus(reqLogger, request, dbPod); err != nil {
		return reconcile.Result{}, err
	}

	// Update status for service database found
	if err := r.updateServiceDatabaseFoundStatus(reqLogger, request, dbService); err != nil {
		return reconcile.Result{}, err
	}

	// Update status for CronJobStatus
	cronJobStatus, err := r.updateCronJobStatus(reqLogger, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update status for database secret
	dbSecretStatus, err := r.updateDBSecretStatus(reqLogger, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update status for aws secret
	awsSecretStatus, err := r.updateAWSSecretStatus(reqLogger, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update status for pod database found
	if err := r.updateEncSecretStatus(reqLogger, request); err != nil {
		return reconcile.Result{}, err
	}

	// Update status for Backup
	if err := r.updateBackupStatus(reqLogger, cronJobStatus, dbSecretStatus, awsSecretStatus, dbPod, dbService, request); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
