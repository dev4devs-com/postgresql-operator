package postgresql

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//updateDBStatus returns error when status regards the all required resources could not be updated
func (r *ReconcilePostgresql) updateDBStatus( deploymentStatus *appsv1.Deployment, serviceStatus *corev1.Service, pvcStatus *corev1.PersistentVolumeClaim, request reconcile.Request) error {
	//Get the latest version of the CR
	instance, err := r.fetchPostgreSQLCR(request)
	if err != nil {
		return err
	}

	// Check if ALL required objects are created
	if len(deploymentStatus.Name) < 1 && len(serviceStatus.Name) < 1 && len(pvcStatus.Name) < 1 {
		err := fmt.Errorf("Failed to get OK Status for PostgreSQL Database")
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
func (r *ReconcilePostgresql) updateDeploymentStatus( request reconcile.Request) (*appsv1.Deployment, error) {
	// Get the latest version of the instance CR
	instance, err := r.fetchPostgreSQLCR(request)
	if err != nil {
		return nil, err
	}
	// Get the deployment Object
	deploymentStatus, err := r.fetchDBDeployment(instance)
	if err != nil {
		return deploymentStatus, err
	}
	// Update the deployment  and Status
	if !reflect.DeepEqual(deploymentStatus.Status, instance.Status.DeploymentStatus) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchPostgreSQLCR(request)
		if err != nil {
			return nil, err
		}

		// Set the data
		instance.Status.DeploymentStatus = deploymentStatus.Status

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return deploymentStatus, err
		}
	}

	return deploymentStatus, nil
}

//updateServiceStatus returns error when status regards the service resource could not be updated
func (r *ReconcilePostgresql) updateServiceStatus( request reconcile.Request) (*corev1.Service, error) {
	// Get the latest version of the instance CR
	instance, err := r.fetchPostgreSQLCR(request)
	if err != nil {
		return nil, err
	}
	// Get the service Object
	serviceStatus, err := r.fetchDBService(instance)
	if err != nil {
		return serviceStatus, err
	}

	if !reflect.DeepEqual(serviceStatus.Status, instance.Status.ServiceStatus) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchPostgreSQLCR(request)
		if err != nil {
			return nil, err
		}

		instance.Status.ServiceStatus = serviceStatus.Status

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return serviceStatus, err
		}
	}

	return serviceStatus, nil
}

//updatePvcStatus returns error when status regards the PersistentVolumeClaim resource could not be updated
func (r *ReconcilePostgresql) updatePvcStatus( request reconcile.Request) (*corev1.PersistentVolumeClaim, error) {
	// Get the latest version of the CR
	instance, err := r.fetchPostgreSQLCR(request)
	if err != nil {
		return nil, err
	}

	// Get pvc Object
	pvcStatus, err := r.fetchDBPersistentVolumeClaim(instance)
	if err != nil {
		return pvcStatus, err
	}

	// Update CR with pvc name
	if !reflect.DeepEqual(pvcStatus.Status, instance.Status.PVCStatus) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchPostgreSQLCR(request)
		if err != nil {
			return nil, err
		}

		// Set the data
		instance.Status.PVCStatus = pvcStatus.Status

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return pvcStatus, err
		}
	}
	return pvcStatus, nil
}
