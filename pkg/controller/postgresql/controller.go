package postgresql

import (
	"context"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_postgresql")

const (
	deployment = "deployment"
	pvc        = "pvc"
	service    = "service"
)

// Add creates a new PostgreSQL Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcilePostgresql{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("postgresql-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Postgresql
	err = c.Watch(&source.Kind{Type: &v1alpha1.Postgresql{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	/** Watch for changes to secondary resources and create the owner PostgreSQL **/

	// deployment
	if err := watchDeployment(c); err != nil {
		return err
	}

	// service
	if err := watchService(c); err != nil {
		return err
	}

	// PersistenceVolume
	if err := watchPersistenceVolumeClaim(c); err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcilePostgresql implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcilePostgresql{}

// ReconcilePostgresql reconciles a PostgreSQL object
type ReconcilePostgresql struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Update the object and reconcile it
func (r *ReconcilePostgresql) update(obj runtime.Object, reqLogger logr.Logger) error {
	err := r.client.Update(context.TODO(), obj)
	if err != nil {
		reqLogger.Error(err, "Failed to update Object", "obj:", obj)
		return err
	}
	reqLogger.Info("Object updated", "obj:", obj)
	return nil
}

// Create the object and reconcile it
func (r *ReconcilePostgresql) create(db *v1alpha1.Postgresql, kind string, reqLogger logr.Logger) error {
	obj := r.buildFactory(db, kind, reqLogger)
	reqLogger.Info("Creating a new ", "kind", kind, "Namespace", db.Namespace)
	err := r.client.Create(context.TODO(), obj)
	if err != nil {
		reqLogger.Error(err, "Failed to create new ", "kind", kind, "Namespace", db.Namespace)
		return err
	}
	reqLogger.Info("Created successfully", "kind", kind, "Namespace", db.Namespace)
	return nil
}

// buildFactory will return the resource according to the kind defined
func (r *ReconcilePostgresql) buildFactory(db *v1alpha1.Postgresql, kind string, reqLogger logr.Logger) runtime.Object {
	reqLogger.Info("Check "+kind, "into the namespace", db.Namespace)
	switch kind {
	case pvc:
		return r.buildPVCForDB(db)
	case deployment:
		return r.buildDBDeployment(db)
	case service:
		return r.buildDBService(db)
	default:
		msg := "Failed to recognize type of object" + kind + " into the Namespace " + db.Namespace
		panic(msg)
	}
}

// Reconcile reads that state of the cluster for a Postgresql object and makes changes based on the state read
// and what is in the Postgresql.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcilePostgresql) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling PostgreSQL Database")

	// Fetch the PostgreSQL DB
	db := &v1alpha1.Postgresql{}
	db, err := r.fetchDBInstance(reqLogger, request)
	if err != nil {
		reqLogger.Error(err, "Failed to get PostgreSQL Custom Resource")
		return reconcile.Result{}, err
	}

	// Add const values for mandatory specs
	addMandatorySpecsDefinitions(db)

	// Check if deployment for the app exist, if not create one
	dep, err := r.fetchDBDeployment(reqLogger, db)
	if err != nil {
		// Create the deployment
		if err := r.create(db, deployment, reqLogger); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}

	// Check if service for the app exist, if not create one
	if _, err := r.fetchDBService(reqLogger, db); err != nil {
		if err := r.create(db, service, reqLogger); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Check if PersistentVolumeClaim for the app exist, if not create one
	if _, err := r.fetchDBPersistentVolumeClaim(reqLogger, db); err != nil {
		if err := r.create(db, pvc, reqLogger); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Ensure the deployment size is the same as the spec
	reqLogger.Info("Ensuring the PostgreSQL Database deployment size is the same as the spec")
	size := db.Spec.Size
	if *dep.Spec.Replicas != size {
		// Set the number of Replicas spec in the CR
		dep.Spec.Replicas = &size
		// Update
		if err := r.update(dep, reqLogger); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Update status for deployment
	deploymentStatus, err := r.updateDeploymentStatus(reqLogger, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update status for service
	serviceStatus, err := r.updateServiceStatus(reqLogger, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update status for pvc
	pvcStatus, err := r.updatePvcStatus(reqLogger, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update status for DB
	if err := r.updateDBStatus(reqLogger, deploymentStatus, serviceStatus, pvcStatus, request); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
