package resource

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//Returns the buildDatabaseSecret object for the PostgreSQL Backup
func NewBackupSecret(bkp *v1alpha1.Backup, prefix string, secretData map[string][]byte, secretStringData map[string]string, scheme *runtime.Scheme) *corev1.Secret {
	ls := utils.GetLabels(bkp.Name)

	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      prefix + bkp.Name,
			Namespace: bkp.Namespace,
			Labels:    ls,
		},
		Data: secretData,
		Type: "Opaque",
	}

	if secretStringData != nil && len(secretStringData) > 0 {
		secret.StringData = secretStringData
	}

	controllerutil.SetControllerReference(bkp, secret, scheme)
	return secret
}
