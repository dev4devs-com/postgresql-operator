package postgresql

import (
	"context"
	"github.com/dev4devs-com/postgresql-operator/pkg/config"
	"reflect"
	"testing"

	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcilePostgreSQL_Update(t *testing.T) {
	type fields struct {
		createdInstance  *v1alpha1.Postgresql
		instanceToUpdate *v1alpha1.Postgresql
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
				createdInstance:  &dbInstance,
				instanceToUpdate: &dbInstance,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.fields.createdInstance}

			r := buildReconcileWithFakeClientWithMocks(objs)

			err := r.update(tt.fields.instanceToUpdate)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcilePostgreSQL_Update.update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestReconcilePostgreSQL_Create(t *testing.T) {
	type fields struct {
		scheme *runtime.Scheme
	}
	type args struct {
		instance *v1alpha1.Postgresql
		kind     string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "Should create and return a new deployment",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance: &dbInstance,
				kind:     deployment,
			},
			wantErr: false,
		},
		{
			name: "Should fail to create an unknown type",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance: &dbInstance,
				kind:     "UNKNOWN",
			},
			wantErr:   false,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.args.instance}

			r := buildReconcileWithFakeClientWithMocks(objs)


			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("TestReconcilePostgreSQL_Create.create() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			err := r.create(tt.args.instance, tt.args.kind)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcilePostgreSQL_Create.create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestReconcilePostgreSQL_BuildFactory(t *testing.T) {
	type fields struct {
		scheme *runtime.Scheme
	}
	type args struct {
		instance *v1alpha1.Postgresql
		kind     string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      reflect.Type
		wantPanic bool
	}{
		{
			name: "Should create a deployment",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance: &dbInstance,
				kind:     deployment,
			},
			want: reflect.TypeOf(&appsv1.Deployment{}),
		},
		{
			name: "Should create a Persistent Volume Claim",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance: &dbInstance,
				kind:     pvc,
			},
			want: reflect.TypeOf(&corev1.PersistentVolumeClaim{}),
		},
		{
			name: "Should create a service",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance: &dbInstance,
				kind:     service,
			},
			want: reflect.TypeOf(&corev1.Service{}),
		},
		{
			name: "Should panic when trying to create unrecognized object type",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance: &dbInstance,
				kind:     "UNDEFINED",
			},
			wantPanic: true,
			want:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.args.instance}
			r := buildReconcileWithFakeClientWithMocks(objs)

			// testing if the panic() function is executed
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("TestReconcilePostgreSQL_BuildFactory.buildFactory() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			got := r.buildFactory(tt.args.instance, tt.args.kind)

			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("TestReconcilePostgreSQL_BuildFactory.buildFactory() = %v, want %v", gotType, tt.want)
			}
		})
	}
}

func TestReconcilePostgreSQL_Success(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&dbInstance,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      dbInstance.Name,
			Namespace: dbInstance.Namespace,
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	deployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), req.NamespacedName, deployment)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	service := &corev1.Service{}
	err = r.client.Get(context.TODO(), req.NamespacedName, service)
	if err != nil {
		t.Fatalf("get service: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	pvc := &corev1.PersistentVolumeClaim{}
	err = r.client.Get(context.TODO(), req.NamespacedName, pvc)
	if err != nil {
		t.Fatalf("get pvc: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}

func TestReconcilePostgreSQL_NotFound(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&dbInstance,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      dbInstance.Name,
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

func TestReconcilePostgreSQL_UsingConfigMapToCreateEnvVars(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&dbInstanceConfigMapSameKeys,
		&configMapSameKeyValues,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      dbInstance.Name,
			Namespace: dbInstance.Namespace,
		},
	}

	_ = r.client.Create(context.TODO(), &configMapSameKeyValues)

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	deployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), req.NamespacedName, deployment)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom == nil {
		t.Error("deployment envvar did not came from service instance config map")
	}

	if deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.ConfigMapKeyRef.Name != configMapSameKeyValues.Name {
		t.Fatalf("deployment envvar did not came from service instance config map: (%v,%v)", deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.ConfigMapKeyRef.Name, configMapSameKeyValues.Name)
	}

	if res.Requeue {
		t.Error("reconcile did not requeue request as expected")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	service := &corev1.Service{}
	err = r.client.Get(context.TODO(), req.NamespacedName, service)
	if err != nil {
		t.Fatalf("get service: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	pvc := &corev1.PersistentVolumeClaim{}
	err = r.client.Get(context.TODO(), req.NamespacedName, pvc)
	if err != nil {
		t.Fatalf("get pvc: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}

func TestReconcilePostgreSQ_ReplicasSizes(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&dbInstance,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      dbInstance.Name,
			Namespace: dbInstance.Namespace,
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	deployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), req.NamespacedName, deployment)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	//Mock Replicas wrong size
	size := int32(3)
	deployment.Spec.Replicas = &size

	// Update
	err = r.client.Update(context.TODO(), deployment)
	if err != nil {
		t.Fatalf("fails when ttry to update deployment replicas: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	deployment = &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), req.NamespacedName, deployment)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	if *deployment.Spec.Replicas != dbInstance.Spec.Size {
		t.Error("Replicas size was not respected")
	}
}

func TestReconcilePostgreSQL_Reconcile_InstanceWithoutSpec(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&dbInstanceWithoutSpec,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      dbInstanceWithoutSpec.Name,
			Namespace: dbInstanceWithoutSpec.Namespace,
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	deployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), req.NamespacedName, deployment)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	// Check if the quantity of Replicas for this deployment is equals the specification
	if *deployment.Spec.Replicas != config.NewPostgreSQLConfig().Size {
		t.Errorf("dep size (%d) is not the expected size (%d)", deployment.Spec.Replicas, config.NewPostgreSQLConfig().Size)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	service := &corev1.Service{}
	err = r.client.Get(context.TODO(), req.NamespacedName, service)
	if err != nil {
		t.Fatalf("get service: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	pvc := &corev1.PersistentVolumeClaim{}
	err = r.client.Get(context.TODO(), req.NamespacedName, pvc)
	if err != nil {
		t.Fatalf("get pvc: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}
