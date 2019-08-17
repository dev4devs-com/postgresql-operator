package backup

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//Returns the buildDatabaseSecret object for the PostgreSQL Backup
func buildSecret(bkp *v1alpha1.Backup, prefix string, secretData map[string][]byte, secretStringData map[string]string, scheme *runtime.Scheme) *corev1.Secret {
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
	controllerutil.SetControllerReference(bkp, secret, scheme)
	return secret
}
