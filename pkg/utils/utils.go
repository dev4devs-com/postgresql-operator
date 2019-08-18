package utils

import "github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"

func GetLabels(name string) map[string]string {
	return map[string]string{"app": "postgresqloperator", "cr": name}
}

// GetAWSSecretName returns the name of the secret
// NOTE: The user can just inform the name and namespace of the Secret which is already applied in the cluster OR
// the data required for the operator be able to create one in the same namespace where the backup is applied
func GetAWSSecretName(bkp *v1alpha1.Backup) string {
	if IsAwsKeySetupByName(bkp) {
		return bkp.Spec.AwsCredentialsSecretName
	}
	return AwsSecretPrefix + bkp.Name
}

// GetAwsSecretNamespace returns the namespace where the secret is applied already
// NOTE: The user can just inform the name and namespace of the Secret which is already applied in the cluster OR
// the data required for the operator be able to create one in the same namespace where the backup is applied
func GetAwsSecretNamespace(bkp *v1alpha1.Backup) string {
	if IsAwsKeySetupByName(bkp) && bkp.Spec.AwsCredentialsSecretNamespace != "" {
		return bkp.Spec.AwsCredentialsSecretNamespace
	}
	return bkp.Namespace
}

// GetEncSecretNamespace returns the namespace where the secret is applied already
// NOTE: The user can just inform the name and namespace of the Secret which is already applied in the cluster OR
// the data required for the operator be able to create one in the same namespace where the backup is applied
func GetEncSecretNamespace(bkp *v1alpha1.Backup) string {
	if IsEncKeySetupByNameAndNamaspace(bkp) {
		return bkp.Spec.EncryptKeySecretNamespace
	}
	return bkp.Namespace
}

// GetEncSecretName returns the name of the secret
// NOTE: The user can just inform the name and namespace of the Secret which is already applied in the cluster OR
// the data required for the operator be able to create one in the same namespace where the backup is applied
func GetEncSecretName(bkp *v1alpha1.Backup) string {
	if IsEncKeySetupByName(bkp) {
		return bkp.Spec.EncryptKeySecretName
	}
	return EncSecretPrefix + bkp.Name
}

// IsEncryptionKeyOptionConfig returns true when the CR has the configuration to allow it be used
func IsEncryptionKeyOptionConfig(bkp *v1alpha1.Backup) bool {
	return bkp.Spec.AwsCredentialsSecretName != "" ||
		(bkp.Spec.GpgTrustModel != "" && bkp.Spec.GpgEmail != "" && bkp.Spec.GpgPublicKey != "")
}

// IsEncKeySetupByName returns true when it is setup to get an pre-existing secret applied in the cluster.
// NOTE: The user can just inform the name of the Secret which is already applied in the cluster OR
// the data required for the operator be able to create one
func IsEncKeySetupByName(bkp *v1alpha1.Backup) bool {
	return bkp.Spec.EncryptKeySecretName != ""
}

// IsAwsKeySetupByName returns true when it is setup to get an pre-existing secret applied in the cluster.
// NOTE: The user can just inform the name of the Secret which is already applied in the cluster OR
// the data required for the operator be able to create one
func IsAwsKeySetupByName(bkp *v1alpha1.Backup) bool {
	return bkp.Spec.AwsCredentialsSecretName != ""
}

func IsEncKeySetupByNameAndNamaspace(bkp *v1alpha1.Backup) bool {
	return IsEncKeySetupByName(bkp) && bkp.Spec.EncryptKeySecretNamespace != ""
}
