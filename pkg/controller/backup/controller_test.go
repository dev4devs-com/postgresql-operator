package backup

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	"k8s.io/api/batch/v1beta1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
	"context"
	corev1 "k8s.io/api/core/v1"

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

func TestReconcileBackup_Create(t *testing.T) {
	objs := []runtime.Object{&bkpInstance}
	r := buildReconcileWithFakeClientWithMocks(objs)
	dataDBSecret, _ := r.buildDBSecretData(&bkpInstance, &dbInstanceWithoutSpec)
	awsDataSecret := buildAwsSecretData(&bkpInstance)
	type fields struct {
		scheme *runtime.Scheme
	}
	type args struct {
		instance   *v1alpha1.Backup
		kind       string
		secretData map[string][]byte
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "Should create and return a new cronJob",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       cronJob,
				secretData: nil,
			},
			wantErr: false,
		},
		{
			name: "Should create and return a new DB Secret",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       dbSecret,
				secretData: dataDBSecret,
			},
			wantErr: false,
		},
		{
			name: "Should create and return a new Aws cronJob",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       swsSecret,
				secretData: awsDataSecret,
			},
			wantErr: false,
		},
		{
			name: "Should fail to create an unknown type",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance: &bkpInstance,
				kind:     "UNKNOWN",
			},
			wantErr:   false,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("TestReconcileBackup_Create.create() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			err := r.create(tt.args.instance, &dbInstanceWithoutSpec, tt.args.kind)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcileBackup_Create.create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestReconcileBackup_BuildFactory(t *testing.T) {
	objs := []runtime.Object{&bkpInstance}
	r := buildReconcileWithFakeClientWithMocks(objs)
	dataDBSecret, _ := r.buildDBSecretData(&bkpInstance, &dbInstanceWithoutSpec)
	awsDataSecret := buildAwsSecretData(&bkpInstance)
	encDataSecret, _ := buildEncSecretData(&bkpInstanceWithSecretNames)

	type fields struct {
		scheme *runtime.Scheme
	}
	type args struct {
		instance   *v1alpha1.Backup
		kind       string
		secretData map[string][]byte
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      reflect.Type
		wantPanic bool
	}{
		{
			name: "Should create a cronJob",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       cronJob,
				secretData: nil,
			},
			want: reflect.TypeOf(&v1beta1.CronJob{}),
		},
		{
			name: "Should create a DB Secret",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       dbSecret,
				secretData: dataDBSecret,
			},
			want: reflect.TypeOf(&v1.Secret{}),
		},
		{
			name: "Should create a Aws Secret",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       swsSecret,
				secretData: awsDataSecret,
			},
			want: reflect.TypeOf(&v1.Secret{}),
		},
		{
			name: "Should create a Enc Secret",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstanceWithSecretNames,
				kind:       encSecret,
				secretData: encDataSecret,
			},
			want: reflect.TypeOf(&v1.Secret{}),
		},
		{
			name: "Should panic when trying to create unrecognized object type",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance: &bkpInstance,
				kind:     "UNDEFINED",
			},
			wantPanic: true,
			want:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// testing if the panic() function is executed
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("TestReconcileBackup_BuildFactory.buildFactory() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			got, _ := r.buildFactory(tt.args.instance, &dbInstanceWithoutSpec, tt.args.kind)

			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("TestReconcileBackup_BuildFactory.buildFactory() = %v, want %v", gotType, tt.want)
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