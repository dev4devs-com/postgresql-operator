package postgresql

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//updateDBStatus returns error when status regards the all required resources could not be updated
func (r *ReconcilePostgresql) updateDBStatus(reqLogger logr.Logger, deploymentStatus *appsv1.Deployment, serviceStatus *corev1.Service, pvcStatus *corev1.PersistentVolumeClaim, request reconcile.Request) error {
	reqLogger.Info("Updating App Status for the PostgreSQL")

	//Get the latest version of the CR
	instance, err := r.fetchDBInstance(reqLogger, request)
	if err != nil {
		return err
	}

	// Check if ALL required objects are created
	if len(deploymentStatus.Name) < 1 && len(serviceStatus.Name) < 1 && len(pvcStatus.Name) < 1 {
		err := fmt.Errorf("Failed to get OK Status for PostgreSQL Database")
		reqLogger.Error(err, "One of the resources are not created", "PostgreSQL.Namespace", instance.Namespace, "PostgreSQL.Name", instance.Name)
		return err
	}
	status := "OK"

	// Update Database Status == OK
	if !reflect.DeepEqual(status, instance.Status.DatabaseStatus) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchDBInstance(reqLogger, request)
		if err != nil {
			return err
		}

		// Set the data
		instance.Status.DatabaseStatus = status

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update deployment Status for the PostgreSQL Database")
			return err
		}
	}
	return nil
}

//updateDeploymentStatus returns error when status regards the deployment resource could not be updated
func (r *ReconcilePostgresql) updateDeploymentStatus(reqLogger logr.Logger, request reconcile.Request) (*appsv1.Deployment, error) {
	reqLogger.Info("Updating deployment Status for the PostgreSQL")
	// Get the latest version of the instance CR
	instance, err := r.fetchDBInstance(reqLogger, request)
	if err != nil {
		return nil, err
	}
	// Get the deployment Object
	deploymentStatus, err := r.fetchDBDeployment(reqLogger, instance)
	if err != nil {
		reqLogger.Error(err, "Failed to get deployment for Status", "PostgreSQL.Namespace", instance.Namespace, "PostgreSQL.Name", instance.Name)
		return deploymentStatus, err
	}
	// Update the deployment  and Status
	if !reflect.DeepEqual(deploymentStatus.Status, instance.Status.DeploymentStatus) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchDBInstance(reqLogger, request)
		if err != nil {
			return nil, err
		}

		// Set the data
		instance.Status.DeploymentStatus = deploymentStatus.Status

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update deployment Status for the PostgreSQL")
			return deploymentStatus, err
		}
	}

	return deploymentStatus, nil
}

//updateServiceStatus returns error when status regards the service resource could not be updated
func (r *ReconcilePostgresql) updateServiceStatus(reqLogger logr.Logger, request reconcile.Request) (*corev1.Service, error) {
	reqLogger.Info("Updating service Status for the PostgreSQL")
	// Get the latest version of the instance CR
	instance, err := r.fetchDBInstance(reqLogger, request)
	if err != nil {
		return nil, err
	}
	// Get the service Object
	serviceStatus, err := r.fetchDBService(reqLogger, instance)
	if err != nil {
		reqLogger.Error(err, "Failed to get service for Status", "PostgreSQL.Namespace", instance.Namespace, "PostgreSQL.Name", instance.Name)
		return serviceStatus, err
	}

	if !reflect.DeepEqual(serviceStatus.Status, instance.Status.ServiceStatus) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchDBInstance(reqLogger, request)
		if err != nil {
			return nil, err
		}

		instance.Status.ServiceStatus = serviceStatus.Status

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update service Status for the PostgreSQL")
			return serviceStatus, err
		}
	}

	return serviceStatus, nil
}

//updatePvcStatus returns error when status regards the PersistentVolumeClaim resource could not be updated
func (r *ReconcilePostgresql) updatePvcStatus(reqLogger logr.Logger, request reconcile.Request) (*corev1.PersistentVolumeClaim, error) {
	reqLogger.Info("Updating PersistentVolumeClaim Status for the PostgreSQL")
	// Get the latest version of the CR
	instance, err := r.fetchDBInstance(reqLogger, request)
	if err != nil {
		return nil, err
	}

	// Get pvc Object
	pvcStatus, err := r.fetchDBPersistentVolumeClaim(reqLogger, instance)
	if err != nil {
		reqLogger.Error(err, "Failed to get PersistentVolumeClaim for Status", "PostgreSQL.Namespace", instance.Namespace, "PostgreSQL.Name", instance.Name)
		return pvcStatus, err
	}

	// Update CR with pvc name
	if !reflect.DeepEqual(pvcStatus.Status, instance.Status.PVCStatus) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchDBInstance(reqLogger, request)
		if err != nil {
			return nil, err
		}

		// Set the data
		instance.Status.PVCStatus = pvcStatus.Status

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update PersistentVolumeClaim Status for the PostgreSQL")
			return pvcStatus, err
		}
	}
	return pvcStatus, nil
}
