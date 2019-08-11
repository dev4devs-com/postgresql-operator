package backup

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	awsSecretPrefix     = "aws-"
	dbSecretPrefix      = "db-"
	encryptionKeySecret = "encryption-"
)

//Returns the buildDatabaseSecret object for the PostgreSQL Backup
func (r *ReconcileBackup) buildSecret(bkp *v1alpha1.Backup, prefix string, secretData map[string][]byte, secretStringData map[string]string) *corev1.Secret {
	ls := getBkpLabels(bkp.Name)

	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      prefix + bkp.Name,
			Namespace: bkp.Namespace,
			Labels:    ls,
		},
		Data: secretData,
		Type: "Opaque",
	}

	// Add string data
	if secretStringData != nil && len(secretStringData) > 0 {
		secret.StringData = secretStringData
	}

	// Set Backup as the owner and controller
	controllerutil.SetControllerReference(bkp, secret, r.scheme)
	return secret
}
