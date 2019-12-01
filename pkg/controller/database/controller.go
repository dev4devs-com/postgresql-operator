package database

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/service"
	"github.com/dev4devs-com/postgresql-operator/pkg/utils"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new Database Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDatabase{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(utils.DatabaseControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Database
	if err := c.Watch(&source.Kind{Type: &v1alpha1.Database{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	/** Watch for changes to secondary resource and create the owner Database **/

	// Watch Deployment resource controlled and created by it
	if err := service.Watch(c, &v1.Deployment{}, true, &v1alpha1.Database{}); err != nil {
		return err
	}

	// Watch PersistenceVolumeClaim resource controlled and created by it
	if err := service.Watch(c, &corev1.PersistentVolumeClaim{}, true, &v1alpha1.Database{}); err != nil {
		return err
	}

	// Watch Service resource controlled and created by it
	if err := service.Watch(c, &corev1.Service{}, true, &v1alpha1.Database{}); err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileDatabase implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileDatabase{}

// ReconcileDatabase reconciles a Database object
type ReconcileDatabase struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Database object and makes changes based on the state read
// and what is in the Database.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileDatabase) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := utils.GetLoggerByRequestAndController(request, utils.DatabaseControllerName)
	reqLogger.Info("Reconciling Database ...")

	db, err := service.FetchDatabaseCR(request.Name, request.Namespace, r.client)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("Database resource not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get Database.")
		return reconcile.Result{}, err
	}

	// Add const values for mandatory specs
	utils.AddDatabaseMandatorySpecs(db)

	if err := r.createResources(db, request); err != nil {
		reqLogger.Error(err, "Failed to create the secondary resource required for the Database CR")
		return reconcile.Result{}, err
	}

	if err := r.manageResources(db); err != nil {
		reqLogger.Error(err, "Failed to manage resource required for the Database CR")
		return reconcile.Result{}, err
	}

	if err := r.createUpdateCRStatus(request); err != nil {
		reqLogger.Error(err, "Failed to create and update the status in the Database CR")
		return reconcile.Result{}, err
	}

	reqLogger.Info("Stop Reconciling Database ...")
	return reconcile.Result{}, nil
}

//createResources will create the secondary resource which are required in order to make works successfully the primary resource(CR)
func (r *ReconcileDatabase) createResources(db *v1alpha1.Database, request reconcile.Request) error {
	reqLogger := utils.GetLoggerByRequestAndController(request, utils.DatabaseControllerName)
	reqLogger.Info("Creating secondary Database resources ...")

	// Check if deployment for the app exist, if not create one
	if err := r.createDeployment(db); err != nil {
		reqLogger.Error(err, "Failed to create Deployment")
		return err
	}

	// Check if service for the app exist, if not create one
	if err := r.createService(db); err != nil {
		reqLogger.Error(err, "Failed to create Service")
		return err
	}

	// Check if PersistentVolumeClaim for the app exist, if not create one
	if err := r.createPvc(db); err != nil {
		reqLogger.Error(err, "Failed to create PVC")
		return err
	}

	return nil
}

//createUpdateCRStatus will create and update the status in the CR applied in the cluster
func (r *ReconcileDatabase) createUpdateCRStatus(request reconcile.Request) error {
	reqLogger := utils.GetLoggerByRequestAndController(request, utils.DatabaseControllerName)
	reqLogger.Info("Create/Update Database status ...")

	if err := r.updateDeploymentStatus(request); err != nil {
		reqLogger.Error(err, "Failed to create Deployment Status")
		return err
	}

	if err := r.updateServiceStatus(request); err != nil {
		reqLogger.Error(err, "Failed to create Service Status")
		return err
	}

	if err := r.updatePvcStatus(request); err != nil {
		reqLogger.Error(err, "Failed to create PVC Status")
		return err
	}

	if err := r.updateDBStatus(request); err != nil {
		reqLogger.Error(err, "Failed to create DB Status")
		return err
	}
	return nil
}
