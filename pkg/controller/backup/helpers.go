package backup

import (
	"fmt"
	 "github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
)

func getBkpLabels(name string) map[string]string {
	return map[string]string{"app": "postgresql", "backup_cr": name}
}

func (r *ReconcileBackup) buildDBSecretData(bkp *v1alpha1.Backup) (map[string][]byte, error) {
	database := ""
	user := ""
	pwd := ""
	host := r.dbService.Name + "." + bkp.Namespace + ".svc"
	superuser := "false"

	for i := 0; i < len(r.dbPod.Spec.Containers[0].Env); i++ {
		value := r.dbPod.Spec.Containers[0].Env[i].Value
		switch r.dbPod.Spec.Containers[0].Env[i].Name {
		case "POSTGRESQL_DATABASE":
			// Get value from ENV VAR
			database = value

			// Get value from ConfigMap used in the ENV VAR
			if database == "" {
				// get configMap name and key
				cfgName := r.dbPod.Spec.Containers[0].Env[i].ValueFrom.ConfigMapKeyRef.Name
				cfgKey := r.dbPod.Spec.Containers[0].Env[i].ValueFrom.ConfigMapKeyRef.Key
				value, err := r.getValueFromConfigMap(cfgName, cfgKey, bkp)

				// validations
				if err != nil {
					return nil, err
				}
				if value == "" {
					err := fmt.Errorf("Unable to get the database name to add in the secret")
					return nil, err
				}

				// Set ENV value from ConfigMap
				database = value
			}
		case "POSTGRESQL_USER":
			// Get value from ENV VAR
			user = value

			// Get value from ConfigMap used in the ENV VAR
			if user == "" {
				// get configMap name and key
				cfgName := r.dbPod.Spec.Containers[0].Env[i].ValueFrom.ConfigMapKeyRef.Name
				cfgKey := r.dbPod.Spec.Containers[0].Env[i].ValueFrom.ConfigMapKeyRef.Key
				value, err := r.getValueFromConfigMap(cfgName, cfgKey, bkp)

				// validations
				if err != nil {
					return nil, err
				}
				if value == "" {
					err := fmt.Errorf("Unable to get the database user to add in the secret")
					return nil, err
				}

				// Set ENV value from ConfigMap
				user = value
			}
		case "POSTGRESQL_PASSWORD":
			// Get value from ENV VAR
			pwd = value

			// Get value from ConfigMap used in the ENV VAR
			if pwd == "" {
				// get configMap name and key
				cfgName := r.dbPod.Spec.Containers[0].Env[i].ValueFrom.ConfigMapKeyRef.Name
				cfgKey := r.dbPod.Spec.Containers[0].Env[i].ValueFrom.ConfigMapKeyRef.Key
				value, err := r.getValueFromConfigMap(cfgName, cfgKey, bkp)

				// validations
				if err != nil {
					return nil, err
				}
				if value == "" {
					err := fmt.Errorf("Unable to get the database pwd to add in the secret")
					return nil, err
				}

				// Set ENV value from ConfigMap
				pwd = value
			}
		}
	}

	dataByte := map[string][]byte{
		"POSTGRES_USERNAME":  []byte(user),
		"POSTGRES_PASSWORD":  []byte(pwd),
		"POSTGRES_DATABASE":  []byte(database),
		"POSTGRES_HOST":      []byte(host),
		"POSTGRES_SUPERUSER": []byte(superuser),
		"VERSION":            []byte(bkp.Spec.DatabaseVersion),
	}

	return dataByte, nil
}

func (r *ReconcileBackup) getValueFromConfigMap(configMapName, configMapKey string, bkp *v1alpha1.Backup) (string, error) {
	// search for ConfigMap
	cfg, err := r.fetchConfigMap(bkp, configMapName)
	if err != nil {
		return "", err
	}
	// Get ENV value
	return cfg.Data[configMapKey], nil
}

func buildAwsSecretData(bkp *v1alpha1.Backup) map[string][]byte {
	dataByte := map[string][]byte{
		"AWS_S3_BUCKET_NAME":    []byte(bkp.Spec.AwsS3BucketName),
		"AWS_ACCESS_KEY_ID":     []byte(bkp.Spec.AwsAccessKeyId),
		"AWS_SECRET_ACCESS_KEY": []byte(bkp.Spec.AwsSecretAccessKey),
	}
	return dataByte
}

func buildEncSecretData(bkp *v1alpha1.Backup) (map[string][]byte, map[string]string) {
	dataByte := map[string][]byte{
		"GPG_PUBLIC_KEY": []byte(bkp.Spec.GpgPublicKey),
	}

	dataString := map[string]string{
		"GPG_RECIPIENT":   bkp.Spec.GpgEmail,
		"GPG_TRUST_MODEL": bkp.Spec.GpgTrustModel,
	}
	return dataByte, dataString
}

func getAWSSecretName(bkp *v1alpha1.Backup) string {
	awsSecretName := awsSecretPrefix + bkp.Name
	if bkp.Spec.AwsCredentialsSecretName != "" {
		awsSecretName = bkp.Spec.AwsCredentialsSecretName
	}
	return awsSecretName
}

func getAwsSecretNamespace(bkp *v1alpha1.Backup) string {
	if bkp.Spec.AwsCredentialsSecretName != "" && bkp.Spec.AwsCredentialsSecretNamespace != "" {
		return bkp.Spec.AwsCredentialsSecretNamespace
	}
	return bkp.Namespace
}

func getEncSecretNamespace(bkp *v1alpha1.Backup) string {
	if hasEncryptionKeySecret(bkp) {
		if bkp.Spec.EncryptionKeySecretName != "" && bkp.Spec.EncryptionKeySecretNamespace != "" {
			return bkp.Spec.EncryptionKeySecretNamespace
		}
		return bkp.Namespace
	}
	return ""
}

func getEncSecretName(bkp *v1alpha1.Backup) string {
	awsSecretName := ""
	if hasEncryptionKeySecret(bkp) {
		awsSecretName = encryptionKeySecret + bkp.Name
	}
	if bkp.Spec.AwsCredentialsSecretName != "" {
		awsSecretName = bkp.Spec.EncryptionKeySecretName
	}
	return awsSecretName
}

func hasEncryptionKeySecret(bkp *v1alpha1.Backup) bool {
	return bkp.Spec.AwsCredentialsSecretName != "" ||
		(bkp.Spec.GpgTrustModel != "" && bkp.Spec.GpgEmail != "" && bkp.Spec.GpgPublicKey != "")
}
