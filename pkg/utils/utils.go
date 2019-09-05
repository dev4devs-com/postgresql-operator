package utils

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"


)

func GetLabels(name string) map[string]string {
	return map[string]string{"owner": "postgresqloperator", "cr": name}
}

// GetAWSSecretName returns the name of the secret
// NOTE: The user can just inform the name and namespace of the Secret which is already applied in the cluster OR
// the data required for the operator be able to create one in the same namespace where the backup is applied
func GetAWSSecretName(bkp *v1alpha1.Backup) string {
	if IsAwsKeySetupByName(bkp) {
		return bkp.Spec.AwsSecretName
	}
	return AwsSecretPrefix + bkp.Name
}

// GetAwsSecretNamespace returns the namespace where the secret is applied already
// NOTE: The user can just inform the name and namespace of the Secret which is already applied in the cluster OR
// the data required for the operator be able to create one in the same namespace where the backup is applied
func GetAwsSecretNamespace(bkp *v1alpha1.Backup) string {
	if IsAwsKeySetupByName(bkp) && bkp.Spec.AwsSecretNamespace != "" {
		return bkp.Spec.AwsSecretNamespace
	}
	return bkp.Namespace
}

// GetEncSecretNamespace returns the namespace where the secret is applied already
// NOTE: The user can just inform the name and namespace of the Secret which is already applied in the cluster OR
// the data required for the operator be able to create one in the same namespace where the backup is applied
func GetEncSecretNamespace(bkp *v1alpha1.Backup) string {
	if IsEncKeySetupByNameAndNamespace(bkp) {
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
	return bkp.Spec.AwsSecretName != "" ||
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
	return bkp.Spec.AwsSecretName != ""
}

// IsEncKeySetupByNameAndNamespace it will return true when the Enc Key is setup by using an preexisting
// secret applied in the cluster.
func IsEncKeySetupByNameAndNamespace(bkp *v1alpha1.Backup) bool {
	return IsEncKeySetupByName(bkp) && bkp.Spec.EncryptKeySecretNamespace != ""
}

func GetLoggerByRequestAndController(request reconcile.Request, controllerName string) logr.Logger {
	var log = logf.Log.WithName(controllerName)
	return log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
}


