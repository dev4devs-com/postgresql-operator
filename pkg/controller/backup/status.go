package backup

import (
	"context"
	"fmt"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//updateAppStatus returns error when status regards the all required resources could not be updated
func (r *ReconcileBackup) updateBackupStatus(request reconcile.Request) error {
	//Get the latest version of the CR
	bkp, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	if err := r.validateBackupRequirements(bkp); err != nil {
		return err
	}

	status := "OK"

	// Update Backup Status == OK
	if !reflect.DeepEqual(status, bkp.Status.BackupStatus) {
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

func (r *ReconcileBackup) updateCronJobStatus(request reconcile.Request) error {
	// Get the latest version of the CR
	instance, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	// Get Object
	cronJob, err := r.fetchCronJob(instance)
	if err != nil {
		return err
	}

	//Update the CR
	if cronJob.Name != instance.Status.CronJobName || !reflect.DeepEqual(cronJob.Status, instance.Status.CronJobStatus) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchBackupCR(request)
		if err != nil {
			return err
		}

		// Set the data
		instance.Status.CronJobStatus = cronJob.Status
		instance.Status.CronJobName = cronJob.Name

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReconcileBackup) updateAWSSecretStatus(request reconcile.Request) error {
	// Get the latest version of the CR
	instance, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	// Get Object
	aws, err := r.fetchSecret(getAwsSecretNamespace(instance), getAWSSecretName(instance))
	if err != nil {
		return err
	}

	data := covertDataSecretToString(aws)

	//Update the CR
	if aws.Name != instance.Status.AWSSecretName || !reflect.DeepEqual(data, instance.Status.AWSSecretData) || aws.Namespace != instance.Status.AwsCredentialsSecretNamespace {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchBackupCR(request)
		if err != nil {
			return err
		}

		// Set the data
		instance.Status.AWSSecretName = aws.Name
		instance.Status.AWSSecretData = data
		instance.Status.AwsCredentialsSecretNamespace = aws.Namespace

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return err
		}
	}

	return nil
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

		data := covertDataSecretToString(encSecret)

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

func covertDataSecretToString(secret *corev1.Secret) map[string]string {
	data := make(map[string]string)
	if secret.Data != nil {
		for k, v := range secret.Data {
			value := ""
			if v != nil {
				value = string(v)
			}
			data[k] = value
		}
		for k, v := range secret.StringData {
			data[k] = v
		}
	}
	return data
}

func (r *ReconcileBackup) updateDBSecretStatus(request reconcile.Request) error {
	// Get the latest version of the CR
	instance, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	// Get Object
	dbSecret, err := r.fetchSecret(instance.Namespace, dbSecretPrefix+instance.Name)
	if err != nil {
		return err
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
			return err
		}

		// Set the data
		instance.Status.DBSecretName = dbSecret.Name
		instance.Status.DBSecretData = data

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcileBackup) updatePodDatabaseFoundStatus(request reconcile.Request, dbPod *corev1.Pod) error {
	// Get the latest version of the CR
	instance, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	//Update the CR
	if r.isDbPodFound() != instance.Status.DatabasePodFound {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchBackupCR(request)
		if err != nil {
			return err
		}

		// Set the data
		instance.Status.DatabasePodFound = r.isDbPodFound()

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcileBackup) updateServiceDatabaseFoundStatus(request reconcile.Request) error {
	// Get the latest version of the CR
	instance, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	//Update the CR
	if r.isDbServiceFound() != instance.Status.DatabaseServiceFound {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchBackupCR(request)
		if err != nil {
			return err
		}

		// Set the data
		instance.Status.DatabaseServiceFound = r.isDbServiceFound()

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return err
		}
	}

	return nil
}

//isDbServiceFound returns false when the database service which should be created by the PostgreSQL controller was not found
func (r *ReconcileBackup) isDbServiceFound() bool {
	return r.dbService != nil && len(r.dbService.Name) > 0
}

//isDbPodFound returns false when the database pod which should be created by the PostgreSQL controller was not found
func (r *ReconcileBackup) isDbPodFound() bool {
	return r.dbService != nil && len(r.dbService.Name) > 0
}

//validateBackupRequirements returns error when some requirement is missing
func (r *ReconcileBackup) validateBackupRequirements(bkp *v1alpha1.Backup) error {

	if !r.isDbPodFound() {
		err := fmt.Errorf("Unable to set OK Status for Backup. The postgresql pod was not found")
		return err
	}

	if !r.isDbServiceFound() {
		err := fmt.Errorf("Unable to set OK Status for Backup. The postgresql database service was not found")
		return err
	}

	_, err := r.fetchSecret(bkp.Namespace, dbSecretPrefix+bkp.Name)
	if err != nil {
		err = fmt.Errorf("Unable to set OK Status for Backup. The DB Secret name %v was not found", dbSecretPrefix+bkp.Name)
		return err
	}

	_, err = r.fetchSecret(getAwsSecretNamespace(bkp), getAWSSecretName(bkp))
	if err != nil {
		err := fmt.Errorf("Unable to set OK Status for Backup. The AWS Secret name %v in the namespace %v was not found", getAWSSecretName(bkp), getAwsSecretNamespace(bkp))
		return err
	}

	if hasEncryptionKeySecret(bkp) {
		//Get the latest version of the CR
		_, err := r.fetchSecret(getEncSecretNamespace(bkp), getEncSecretName(bkp))
		if err != nil {
			err := fmt.Errorf("Unable to set OK Status for Backup. The Encript Key configured was not found")
			return err
		}
	}

	//check if the cronJob was created
	_, err = r.fetchCronJob(bkp)
	if err != nil {
		err := fmt.Errorf("Unable to set OK Status for Backup. The CronJob was not found")
		return err
	}

	return nil
}
