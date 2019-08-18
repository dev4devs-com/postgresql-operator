package postgresql

import (
	"reflect"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestUpdateDBStatus(t *testing.T) {
	type fields struct {
		objs []runtime.Object
	}
	type args struct {
		request reconcile.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should update status",
			fields: fields{
				objs: []runtime.Object{&dbInstanceWithoutSpec},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      dbInstanceWithoutSpec.Name,
						Namespace: dbInstanceWithoutSpec.Namespace,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			if err := r.updateDBStatus(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("TestUpdateDBStatus error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateDeploymentStatus(t *testing.T) {
	type fields struct {
		objs []runtime.Object
	}
	type args struct {
		request reconcile.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    reflect.Type
		wantErr bool
	}{
		{
			name: "Should not find the deployment",
			fields: fields{
				objs: []runtime.Object{&dbInstanceWithoutSpec},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      dbInstanceWithoutSpec.Name,
						Namespace: dbInstanceWithoutSpec.Namespace,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Should upddate with success",
			fields: fields{
				objs: []runtime.Object{&dbInstanceWithoutSpec, &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      dbInstanceWithoutSpec.Name,
						Namespace: dbInstanceWithoutSpec.Namespace,
					},
					Status: appsv1.DeploymentStatus{
						Replicas: 3,
					},
				}},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      dbInstanceWithoutSpec.Name,
						Namespace: dbInstanceWithoutSpec.Namespace,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			err := r.updateDeploymentStatus(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestUpdateDeploymentStatus) error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpdateServiceStatus(t *testing.T) {
	type fields struct {
		objs []runtime.Object
	}
	type args struct {
		request reconcile.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should not find the service",
			fields: fields{
				objs: []runtime.Object{&dbInstanceWithoutSpec},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      dbInstanceWithoutSpec.Name,
						Namespace: dbInstanceWithoutSpec.Namespace,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Should update with success",
			fields: fields{
				objs: []runtime.Object{&dbInstanceWithoutSpec, &corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      dbInstanceWithoutSpec.Name,
						Namespace: dbInstanceWithoutSpec.Namespace,
					},
					Status: corev1.ServiceStatus{
						LoadBalancer: corev1.LoadBalancerStatus{
							Ingress: []corev1.LoadBalancerIngress{
								corev1.LoadBalancerIngress{
									IP: "test",
								},
							},
						},
					},
				}},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      dbInstanceWithoutSpec.Name,
						Namespace: dbInstanceWithoutSpec.Namespace,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			err := r.updateServiceStatus(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestUpdateServiceStatus error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpdatePVCStatus(t *testing.T) {
	type fields struct {
		objs []runtime.Object
	}
	type args struct {
		request reconcile.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should not find the pvc",
			fields: fields{
				objs: []runtime.Object{&dbInstanceWithoutSpec},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      dbInstanceWithoutSpec.Name,
						Namespace: dbInstanceWithoutSpec.Namespace,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Should update with success",
			fields: fields{
				objs: []runtime.Object{&dbInstanceWithoutSpec, &corev1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "postgresql",
						Namespace: "postgresql",
					},
					Status: corev1.PersistentVolumeClaimStatus{
						Phase: "test",
					},
				}},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      dbInstanceWithoutSpec.Name,
						Namespace: dbInstanceWithoutSpec.Namespace,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			err := r.updatePvcStatus(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestUpdatePVCStatus error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
