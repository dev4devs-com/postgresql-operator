package postgresql

import (
	"context"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/service"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcilePostgresql(t *testing.T) {
	type fields struct {
		scheme *runtime.Scheme
		objs   []runtime.Object
	}
	type args struct {
		dbInstance v1alpha1.Postgresql
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantRequeue    bool
		wantDeployment bool
		wantService    bool
		wantPVC        bool
		wantErr        bool
	}{
		{
			name: "Should work with default values",
			fields: fields{
				objs: []runtime.Object{&dbInstanceWithoutSpec},
			},
			args: args{
				dbInstance: dbInstanceWithoutSpec,
			},
			wantErr:        false,
			wantRequeue:    false,
			wantDeployment: true,
			wantService:    true,
			wantPVC:        true,
		},
		{
			name: "Should work when is using config map to create env vars",
			fields: fields{
				objs: []runtime.Object{
					&dbInstanceConfigMapSameKeys,
					&configMapSameKeyValues,
				},
			},
			args: args{
				dbInstance: dbInstanceConfigMapSameKeys,
			},
			wantErr:        false,
			wantRequeue:    false,
			wantDeployment: true,
			wantService:    true,
			wantPVC:        true,
		},
		{
			name:           "Should fail because is missing the instance",
			wantErr:        true,
			wantRequeue:    false,
			wantDeployment: false,
			wantService:    false,
			wantPVC:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			// mock request to simulate Reconcile() being called on an event for a watched resource
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      tt.args.dbInstance.Name,
					Namespace: tt.args.dbInstance.Namespace,
				},
			}

			res, err := r.Reconcile(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcilePostgresql reconcile: error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			deployment := &appsv1.Deployment{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: tt.args.dbInstance.Name, Namespace: tt.args.dbInstance.Namespace}, deployment)
			if (err == nil) != tt.wantDeployment {
				t.Errorf("TestReconcilePostgresql to get deployment error = %v, wantDeployment %v", err, tt.wantDeployment)
				return
			}

			service := &corev1.Service{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: tt.args.dbInstance.Name, Namespace: tt.args.dbInstance.Namespace}, service)
			if (err == nil) != tt.wantService {
				t.Errorf("TestReconcilePostgresql to get service error = %v, wantService %v", err, tt.wantService)
				return
			}

			pvc := &corev1.PersistentVolumeClaim{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: tt.args.dbInstance.Name, Namespace: tt.args.dbInstance.Namespace}, pvc)
			if (err == nil) != tt.wantPVC {
				t.Errorf("TestReconcilePostgresql to get service error = %v, wantPVC %v", err, tt.wantPVC)
				return
			}

			if (res.Requeue) != tt.wantRequeue {
				t.Errorf("TestReconcileBackup expect request to requeue res.Requeue = %v, wantRequeue %v", res.Requeue, tt.wantRequeue)
				return
			}
		})
	}
}

func TestReconcilePostgresql_EnsureReplicasSizeInstance(t *testing.T) {

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

	//Mock Replicas wrong size
	size := int32(3)
	deployment.Spec.Replicas = &size

	// Update
	err = r.client.Update(context.TODO(), deployment)
	if err != nil {
		t.Fatalf("fails when try to update deployment replicas: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	dep, err := service.FetchDeployment(dbInstanceWithoutSpec.Name, dbInstanceWithoutSpec.Namespace, r.client)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	if *dep.Spec.Replicas != 1 {
		t.Errorf("Replicas size was not respected got (%v), when is expected (%v)", *dep.Spec.Replicas, 1)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}
