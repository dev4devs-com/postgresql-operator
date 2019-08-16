package backup

import (
	"context"
	"fmt"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//updateAppStatus returns error when status regards the all required resources could not be updated
func (r *ReconcileBackup) updateBackupStatus(cronJob *v1beta1.CronJob, dbSecret, awsSecret *corev1.Secret, dbPod *corev1.Pod, dbService *corev1.Service, request reconcile.Request) error {
	//Get the latest version of the CR
	instance, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	// Check just encSecretPrefix which is Optional
	hasErrorWithEncKey, err := r.hasErrorWithEncKey(instance)
	if err != nil {
		return err
	}

	// Check if ALL required objects are created
	if len(cronJob.Name) < 1 || len(dbSecret.Name) < 1 || len(awsSecret.Name) < 1 || dbPod == nil || len(dbPod.Name) < 1 || dbService == nil || len(dbService.Name) < 1 || hasErrorWithEncKey {
		err := fmt.Errorf("Unable to set OK Status for Backup")
		return err
	}
	status := "OK"

	// Update Backup Status == OK
	if !reflect.DeepEqual(status, instance.Status.BackupStatus) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchBackupCR(request)
		if err != nil {
			return err
		}

		// Set the data
		instance.Status.BackupStatus = status

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return err
		}
	}
	return nil
}

//hasErrorWithEncKey return true when the name or details or informed and was not possible check it
func (r *ReconcileBackup) hasErrorWithEncKey(instance *v1alpha1.Backup) (bool, error) {
	if hasEncryptionKeySecret(instance) {
		//Get the latest version of the CR
		encKey, err := r.fetchSecret(getEncSecretNamespace(instance), getEncSecretName(instance))
		if err != nil {
			return false, err
		}
		if len(encKey.Name) < 1 {
			return true, nil
		}
	}
	return false, nil
}

func (r *ReconcileBackup) updateCronJobStatus(request reconcile.Request) (*v1beta1.CronJob, error) {
	// Get the latest version of the CR
	instance, err := r.fetchBackupCR(request)
	if err != nil {
		return nil, err
	}

	// Get Object
	cronJob, err := r.fetchCronJob(instance)
	if err != nil {
		return cronJob, err
	}

	//Update the CR
	if cronJob.Name != instance.Status.CronJobName || !reflect.DeepEqual(cronJob.Status, instance.Status.CronJobStatus) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchBackupCR(request)
		if err != nil {
			return nil, err
		}

		// Set the data
		instance.Status.CronJobStatus = cronJob.Status
		instance.Status.CronJobName = cronJob.Name

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return cronJob, err
		}
	}

	return cronJob, nil
}

func (r *ReconcileBackup) updateAWSSecretStatus(request reconcile.Request) (*corev1.Secret, error) {
	// Get the latest version of the CR
	instance, err := r.fetchBackupCR(request)
	if err != nil {
		return nil, err
	}

	// Get Object
	aws, err := r.fetchSecret(getAwsSecretNamespace(instance), getAWSSecretName(instance))
	if err != nil {
		return aws, err
	}

	data := make(map[string]string)
	if aws.Data != nil {
		for k, v := range aws.Data {
			value := ""
			if v != nil {
				value = string(v)
			}
			data[k] = value
		}
	}

	//Update the CR
	if aws.Name != instance.Status.AWSSecretName || !reflect.DeepEqual(data, instance.Status.AWSSecretData) || aws.Namespace != instance.Status.AwsCredentialsSecretNamespace {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchBackupCR(request)
		if err != nil {
			return nil, err
		}

		// Set the data
		instance.Status.AWSSecretName = aws.Name
		instance.Status.AWSSecretData = data
		instance.Status.AwsCredentialsSecretNamespace = aws.Namespace

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return aws, err
		}
	}

	return aws, nil
}

func (r *ReconcileBackup) updateEncSecretStatus(request reconcile.Request) error {
	// Get the latest version of the CR
	instance, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	// check if the optinal secret was or not used
	hasEncKeySecret := hasEncryptionKeySecret(instance)

	// If has, update name and data status
	if hasEncKeySecret {
		encSecret, err := r.fetchSecret(getEncSecretNamespace(instance), getEncSecretName(instance))
		if err != nil {
			return err
		}

		data := make(map[string]string)
		if encSecret.Data != nil {
			for k, v := range encSecret.Data {
				value := ""
				if v != nil {
					value = string(v)
				}
				data[k] = value
			}
			for k, v := range encSecret.StringData {
				data[k] = v
			}
		}

		//Update the CR
		if encSecret.Name != instance.Status.AWSSecretName || !reflect.DeepEqual(data, instance.Status.AWSSecretData) {
			// Get the latest version of the CR in order to try to avoid errors when try to update the CR
			instance, err := r.fetchBackupCR(request)
			if err != nil {
				return err
			}

			// Set the data
			instance.Status.EncryptionKeySecretName = encSecret.Name
			instance.Status.EncryptionKeySecretData = data
			instance.Status.EncryptionKeySecretNamespace = encSecret.Namespace

			// Update the CR
			err = r.client.Status().Update(context.TODO(), instance)
			if err != nil {
				return err
			}
		}
	}

	// Update boolean status
	if hasEncKeySecret != instance.Status.HasEncryptionKey {
		// Set the data
		instance.Status.HasEncryptionKey = hasEncKeySecret

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcileBackup) updateDBSecretStatus(request reconcile.Request) (*corev1.Secret, error) {
	// Get the latest version of the CR
	instance, err := r.fetchBackupCR(request)
	if err != nil {
		return nil, err
	}

	// Get Object
	dbSecret, err := r.fetchSecret(instance.Namespace, dbSecretPrefix+instance.Name)
	if err != nil {
		return dbSecret, err
	}

	data := make(map[string]string)
	if dbSecret.Data != nil {
		for k, v := range dbSecret.Data {
			value := ""
			if v != nil {
				value = string(v)
			}
			data[k] = value
		}
	}
	//Update the CR
	if dbSecret.Name != instance.Status.DBSecretName || !reflect.DeepEqual(data, instance.Status.DBSecretData) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchBackupCR(request)
		if err != nil {
			return nil, err
		}

		// Set the data
		instance.Status.DBSecretName = dbSecret.Name
		instance.Status.DBSecretData = data

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return dbSecret, err
		}
	}

	return dbSecret, nil
}

func (r *ReconcileBackup) updatePodDatabaseFoundStatus(request reconcile.Request, dbPod *corev1.Pod) error {
	// Get the latest version of the CR
	instance, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	found := false
	if dbPod != nil && len(dbPod.Name) > 0 {
		found = true
	}

	//Update the CR
	if found != instance.Status.DatabasePodFound {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchBackupCR(request)
		if err != nil {
			return err
		}

		// Set the data
		instance.Status.DatabasePodFound = found

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcileBackup) updateServiceDatabaseFoundStatus(request reconcile.Request, dbService *corev1.Service) error {
	// Get the latest version of the CR
	instance, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	found := false
	if dbService != nil && len(dbService.Name) > 0 {
		found = true
	}

	//Update the CR
	if found != instance.Status.DatabaseServiceFound {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchBackupCR(request)
		if err != nil {
			return err
		}

		// Set the data
		instance.Status.DatabaseServiceFound = found

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return err
		}
	}

	return nil
}
