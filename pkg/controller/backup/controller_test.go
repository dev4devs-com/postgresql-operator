package backup

import (
	"context"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func TestReconcileBackup_Update(t *testing.T) {
	type fields struct {
		createdInstance  *v1alpha1.Backup
		instanceToUpdate *v1alpha1.Backup
		scheme           *runtime.Scheme
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Should successfully update the instance",
			fields: fields{
				createdInstance:  &bkpInstance,
				instanceToUpdate: &bkpInstance,
			},
			wantErr: false,
		},
		{
			name: "Should give an error when the namespace is not found",
			fields: fields{
				createdInstance:  &bkpInstance,
				instanceToUpdate: &bkpInstanceNonDefaultNamespace,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.fields.createdInstance}

			r := buildReconcileWithFakeClientWithMocks(objs)

			err := r.update(tt.fields.instanceToUpdate)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcileBackup_Update.update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestReconcileBackup_NotFound(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&bkpInstance,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      bkpInstance.Name,
			Namespace: "unknown",
		},
	}

	res, err := r.Reconcile(req)
	if err == nil {
		t.Error("should fail since the instance do not exist in the <unknown> nammespace")
	}

	if res.Requeue {
		t.Fatalf("did not expected reconcile to requeue.")
	}
}

func TestReconcileBackup_WithoutPodAndServiceDatabase(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&bkpInstance,
		&dbInstanceWithoutSpec,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      bkpInstance.Name,
			Namespace: "postgresql",
		},
	}

	res, err := r.Reconcile(req)
	if err == nil {
		t.Fatalf("should ono create the cronjob because the db pod and service was not found")
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}

func TestReconcileBackup_WithoutDBInstance(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&bkpInstance,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      bkpInstance.Name,
			Namespace: bkpInstance.Namespace,
		},
	}

	res, err := r.Reconcile(req)
	if err == nil {
		t.Fatalf("should fail since has database CR was not appplied")
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}

func TestReconcileBackup_Success(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&bkpInstance,
		&dbInstanceWithoutSpec,
		&podDatabase,
		&serviceDatabase,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      bkpInstance.Name,
			Namespace: "postgresql",
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	awsSecret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: awsSecretPrefix + bkpInstance.Name, Namespace: bkpInstance.Namespace}, awsSecret)
	if err != nil {
		t.Fatalf("error to get aws secret: (%v)", err)
	}

	dbSecret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: dbSecretPrefix + bkpInstance.Name, Namespace: bkpInstance.Namespace}, dbSecret)
	if err != nil {
		t.Fatalf("error to get db secret: (%v)", err)
	}

	encSecret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: encSecretPrefix + bkpInstance.Name, Namespace: bkpInstance.Namespace}, encSecret)
	if err == nil {
		t.Fatalf("error because should not found encripty secret: (%v)", err)
	}

	cronJob := &v1beta1.CronJob{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: bkpInstance.Name, Namespace: bkpInstance.Namespace}, cronJob)
	if err != nil {
		t.Fatalf("error to get cronJob: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}

func TestReconcileBackup_MissingPodDatabase(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&bkpInstance,
		&dbInstanceWithoutSpec,
		&serviceDatabase,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      bkpInstance.Name,
			Namespace: "postgresql",
		},
	}

	res, err := r.Reconcile(req)
	if err == nil {
		t.Fatalf("should not create the cronjob because the db pod was not found")
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}

func TestReconcileBackup_MissingServiceDatabase(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&bkpInstance,
		&dbInstanceWithoutSpec,
		&podDatabase,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      bkpInstance.Name,
			Namespace: "postgresql",
		},
	}

	res, err := r.Reconcile(req)
	if err == nil {
		t.Fatalf("should not create the cronjob because the service pod was not found")
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}

func TestReconcileBackup_WithSecretNames(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&bkpInstanceWithSecretNames,
		&dbInstanceWithoutSpec,
		&podDatabase,
		&serviceDatabase,
		&awsSecretMock,
		&encSecretMock,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      bkpInstance.Name,
			Namespace: "postgresql",
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	dbSecret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: dbSecretPrefix + bkpInstance.Name, Namespace: bkpInstance.Namespace}, dbSecret)
	if err != nil {
		t.Fatalf("error to get db secret: (%v)", err)
	}

	cronJob := &v1beta1.CronJob{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: bkpInstance.Name, Namespace: bkpInstance.Namespace}, cronJob)
	if err != nil {
		t.Fatalf("error to get cronJob: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}

func TestReconcileBackup_WithEncSecretData(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&bkpInstanceWithEncSecretData,
		&dbInstanceWithoutSpec,
		&podDatabase,
		&serviceDatabase,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      bkpInstance.Name,
			Namespace: "postgresql",
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	awsSecret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: awsSecretPrefix + bkpInstance.Name, Namespace: bkpInstance.Namespace}, awsSecret)
	if err != nil {
		t.Fatalf("error to get aws secret: (%v)", err)
	}

	dbSecret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: dbSecretPrefix + bkpInstance.Name, Namespace: bkpInstance.Namespace}, dbSecret)
	if err != nil {
		t.Fatalf("error to get db secret: (%v)", err)
	}

	encSecret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: encSecretPrefix + bkpInstance.Name, Namespace: bkpInstance.Namespace}, encSecret)
	if err != nil {
		t.Fatalf("error to get enc secret: (%v)", err)
	}

	cronJob := &v1beta1.CronJob{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: bkpInstance.Name, Namespace: bkpInstance.Namespace}, cronJob)
	if err != nil {
		t.Fatalf("error to get cronJob: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}

func TestReconcileBackup_WithConfigMap(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&bkpInstance,
		&dbInstanceWithConfigMap,
		&podDatabaseConfigMap,
		&serviceDatabase,
		&configMapOtherKeyValues,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      bkpInstance.Name,
			Namespace: "postgresql",
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	awsSecret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: awsSecretPrefix + bkpInstance.Name, Namespace: bkpInstance.Namespace}, awsSecret)
	if err != nil {
		t.Fatalf("error to get aws secret: (%v)", err)
	}

	dbSecret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: dbSecretPrefix + bkpInstance.Name, Namespace: bkpInstance.Namespace}, dbSecret)
	if err != nil {
		t.Fatalf("error to get db secret: (%v)", err)
	}

	encSecret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: encSecretPrefix + bkpInstance.Name, Namespace: bkpInstance.Namespace}, encSecret)
	if err == nil {
		t.Fatalf("error because should not found encripty secret: (%v)", err)
	}

	cronJob := &v1beta1.CronJob{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: bkpInstance.Name, Namespace: bkpInstance.Namespace}, cronJob)
	if err != nil {
		t.Fatalf("error to get cronJob: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}

func TestReconcileBackup_WithConfigMapAndWrongKeyValues(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&bkpInstance,
		&dbInstanceWithConfigMap,
		&podDatabaseConfigMap,
		&serviceDatabase,
		&configMapOtherKeyValuesInvalidKeys,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      bkpInstance.Name,
			Namespace: "postgresql",
		},
	}

	res, err := r.Reconcile(req)
	if err.Error() != "Unable to get the database name to add in the secret" {
		t.Fatalf("reconcile: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}

func TestReconcileBackup_WithConfigMapAndCustomizedKeys(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&bkpInstance,
		&dbInstanceWithConfigMapAndCustomizeKeys,
		&podDatabaseConfigMap,
		&serviceDatabase,
		&configMapOtherKeyValuesInvalidKeys,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      bkpInstance.Name,
			Namespace: "postgresql",
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	awsSecret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: awsSecretPrefix + bkpInstance.Name, Namespace: bkpInstance.Namespace}, awsSecret)
	if err != nil {
		t.Fatalf("error to get aws secret: (%v)", err)
	}

	dbSecret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: dbSecretPrefix + bkpInstance.Name, Namespace: bkpInstance.Namespace}, dbSecret)
	if err != nil {
		t.Fatalf("error to get db secret: (%v)", err)
	}

	encSecret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: encSecretPrefix + bkpInstance.Name, Namespace: bkpInstance.Namespace}, encSecret)
	if err == nil {
		t.Fatalf("error because should not found encripty secret: (%v)", err)
	}

	cronJob := &v1beta1.CronJob{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: bkpInstance.Name, Namespace: bkpInstance.Namespace}, cronJob)
	if err != nil {
		t.Fatalf("error to get cronJob: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}
