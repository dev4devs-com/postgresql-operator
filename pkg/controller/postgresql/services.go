package postgresql

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	"k8s.io/apimachinery/pkg/util/intstr"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Returns the service object for the PostgreSQL
func (r *ReconcilePostgresql) buildDBService(db *v1alpha1.Postgresql) *corev1.Service {
	ls := getDBLabels(db.Name)
	ser := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      db.Name,
			Namespace: db.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Type:     corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name: db.Name,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: db.Spec.DatabasePort,
					},
					Port:     db.Spec.DatabasePort,
					Protocol: "TCP",
				},
			},
		},
	}
	// Set PostgreSQL db as the owner and controller
	controllerutil.SetControllerReference(db, ser, r.scheme)
	return ser
}
