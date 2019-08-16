package backup

import (
	"fmt"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/utils"
)

const (
	awsSecretPrefix = "aws-"
	dbSecretPrefix  = "db-"
	encSecretPrefix = "encryption-"
)

func getBkpLabels(name string) map[string]string {
	return map[string]string{"app": "postgresql", "backup_cr": name}
}

func (r *ReconcileBackup) buildDBSecretData(bkp *v1alpha1.Backup, db *v1alpha1.Postgresql) (map[string][]byte, error) {
	database := ""
	user := ""
	pwd := ""
	host := r.dbService.Name + "." + bkp.Namespace + ".svc"
	superuser := "false"

	for i := 0; i < len(r.dbPod.Spec.Containers[0].Env); i++ {

		envVarName := r.dbPod.Spec.Containers[0].Env[i].Name
		envVarValue := r.dbPod.Spec.Containers[0].Env[i].Value

		var cfgName, cfgKey string
		if r.dbPod.Spec.Containers[0].Env[i].ValueFrom != nil {
			cfgName = r.dbPod.Spec.Containers[0].Env[i].ValueFrom.ConfigMapKeyRef.Name
			cfgKey = r.dbPod.Spec.Containers[0].Env[i].ValueFrom.ConfigMapKeyRef.Key
		}

		switch envVarName {
		case utils.GetConfigMapEnvVarKey(db.Spec.ConfigMapDatabaseNameParam, db.Spec.DatabaseNameParam):
			database = envVarValue
			if database == "" {
				database, err := r.getValueFromConfigMap(cfgName, bkp.Namespace, cfgKey)
				if database == "" || err != nil {
					err := fmt.Errorf("Unable to get the database name to add in the secret")
					return nil, err
				}
			}
		case utils.GetConfigMapEnvVarKey(db.Spec.ConfigMapDatabaseUserParam, db.Spec.DatabaseUserParam):
			user = envVarValue
			if user == "" {
				user, err := r.getValueFromConfigMap(cfgName, bkp.Namespace, cfgKey)
				if user == "" || err != nil {
					err := fmt.Errorf("Unable to get the database user to add in the secret")
					return nil, err
				}
			}
		case utils.GetConfigMapEnvVarKey(db.Spec.ConfigMapDatabasePasswordParam, db.Spec.DatabasePasswordParam):
			pwd = envVarValue
			if pwd == "" {
				pwd, err := r.getValueFromConfigMap(cfgName, bkp.Namespace, cfgKey)
				if pwd == "" || err != nil {
					err := fmt.Errorf("Unable to get the pwd user to add in the secret")
					return nil, err
				}
			}
		}
	}
	return getDDBSecretData(user, pwd, database, host, superuser, bkp), nil
}

func getDDBSecretData(user string, pwd string, database string, host string, superuser string, bkp *v1alpha1.Backup) map[string][]byte {
	return map[string][]byte{
		"POSTGRES_USERNAME":  []byte(user),
		"POSTGRES_PASSWORD":  []byte(pwd),
		"POSTGRES_DATABASE":  []byte(database),
		"POSTGRES_HOST":      []byte(host),
		"POSTGRES_SUPERUSER": []byte(superuser),
		"VERSION":            []byte(bkp.Spec.DatabaseVersion),
	}
}

func (r *ReconcileBackup) getValueFromConfigMap(configMapName, configMapNamespace, configMapKey string) (string, error) {
	// search for ConfigMap
	cfg, err := r.fetchConfigMap(configMapName, configMapNamespace)
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
		awsSecretName = encSecretPrefix + bkp.Name
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
