package postgresql

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//Returns the Deployment object for the PostgreSQL
func (r *ReconcilePostgresql) buildPVCForDB(db *v1alpha1.Postgresql) *corev1.PersistentVolumeClaim {
	ls := getDBLabels(db.Name)
	pv := &corev1.PersistentVolumeClaim{
		ObjectMeta: v1.ObjectMeta{
			Name:      db.Name,
			Namespace: db.Namespace,
			Labels:    ls,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(db.Spec.DatabaseStorageRequest),
				},
			},
		},
	}
	// Set PostgreSQL db as the owner and controller
	controllerutil.SetControllerReference(db, pv, r.scheme)
	return pv
}
