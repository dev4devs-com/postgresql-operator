package postgresql

import (
	"context"
	"fmt"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//updateDBStatus returns error when status regards the all required resources could not be updated
func (r *ReconcilePostgresql) updateDBStatus(request reconcile.Request) error {
	//Get the latest version of the CR
	instance, err := r.fetchPostgreSQLCR(request)
	if err != nil {
		return err
	}

	if err := r.validateDatabaseRequirements(instance); err != nil {
		return err
	}

	status := "OK"

	// Update Database Status == OK
	if !reflect.DeepEqual(status, instance.Status.DatabaseStatus) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchPostgreSQLCR(request)
		if err != nil {
			return err
		}

		// Set the data
		instance.Status.DatabaseStatus = status

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return err
		}
	}
	return nil
}

//updateDeploymentStatus returns error when status regards the deployment resource could not be updated
func (r *ReconcilePostgresql) updateDeploymentStatus(request reconcile.Request) error {
	// Get the latest version of the instance CR
	instance, err := r.fetchPostgreSQLCR(request)
	if err != nil {
		return err
	}

	deploymentStatus, err := r.fetchDBDeployment(instance)
	if err != nil {
		return err
	}

	// Update the deployment  and Status
	if !reflect.DeepEqual(deploymentStatus.Status, instance.Status.DeploymentStatus) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchPostgreSQLCR(request)
		if err != nil {
			return err
		}

		// Set the data
		instance.Status.DeploymentStatus = deploymentStatus.Status

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return err
		}
	}

	return nil
}

//updateServiceStatus returns error when status regards the service resource could not be updated
func (r *ReconcilePostgresql) updateServiceStatus(request reconcile.Request) error {
	// Get the latest version of the instance CR
	instance, err := r.fetchPostgreSQLCR(request)
	if err != nil {
		return err
	}

	serviceStatus, err := r.fetchDBService(instance)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(serviceStatus.Status, instance.Status.ServiceStatus) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchPostgreSQLCR(request)
		if err != nil {
			return err
		}

		instance.Status.ServiceStatus = serviceStatus.Status

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return err
		}
	}

	return nil
}

//updatePvcStatus returns error when status regards the PersistentVolumeClaim resource could not be updated
func (r *ReconcilePostgresql) updatePvcStatus(request reconcile.Request) error {
	// Get the latest version of the CR
	instance, err := r.fetchPostgreSQLCR(request)
	if err != nil {
		return err
	}

	pvcStatus, err := r.fetchDBPersistentVolumeClaim(instance)
	if err != nil {
		return err
	}

	// Update CR with pvc name
	if !reflect.DeepEqual(pvcStatus.Status, instance.Status.PVCStatus) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchPostgreSQLCR(request)
		if err != nil {
			return err
		}

		// Set the data
		instance.Status.PVCStatus = pvcStatus.Status

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return err
		}
	}
	return nil
}

//validateBackup returns error when some requirement is missing
func (r *ReconcilePostgresql) validateDatabaseRequirements(db *v1alpha1.Postgresql) error {

	_, err := r.fetchDBPersistentVolumeClaim(db)
	if err != nil {
		err = fmt.Errorf("Unable to set OK Status for PostgreSQL Database. The PVC was not found")
	}

	_, err = r.fetchDBDeployment(db)
	if err != nil {
		err = fmt.Errorf("Unable to set OK Status for PostgreSQL Database. The Deployment was not found")
	}

	_, err = r.fetchDBService(db)
	if err != nil {
		err = fmt.Errorf("Unable to set OK Status for PostgreSQL Database. The Service was not found")
	}

	return nil
}
