package backup

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

//buildReconcileWithFakeClientWithMocks return reconcile with fake client, schemes and mock objects
func buildReconcileWithFakeClientWithMocks(objs []runtime.Object) *ReconcileBackup {
	s := scheme.Scheme

	s.AddKnownTypes(v1alpha1.SchemeGroupVersion, &v1alpha1.Backup{})
	s.AddKnownTypes(v1alpha1.SchemeGroupVersion, &v1alpha1.Postgresql{})

	// create a fake client to mock API calls with the mock objects
	cl := fake.NewFakeClientWithScheme(s, objs...)

	// create a PostgreSQL object with the scheme and fake client
	return &ReconcileBackup{client: cl, scheme: s, dbPod: &podDatabase, dbService: &serviceDatabase}
}
