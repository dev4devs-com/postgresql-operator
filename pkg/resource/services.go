package resource

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/util/intstr"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Returns the service object for the Database
func NewDatabaseService(db *v1alpha1.Database, scheme *runtime.Scheme) (*corev1.Service, error) {
	ls := utils.GetLabels(db.Name)
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
	// Set Database db as the owner and controller
	if err := controllerutil.SetControllerReference(db, ser, scheme); err != nil {
		return nil, err
	}
	return ser, nil
}
