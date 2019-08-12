package postgresql

import (
	"reflect"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcilePostgreSQL_UpdateStatus(t *testing.T) {
	type fields struct {
		objs   []runtime.Object
		scheme *runtime.Scheme
	}
	type args struct {
		deploymentStatus *appsv1.Deployment
		serviceStatus    *corev1.Service
		pvcStatus        *corev1.PersistentVolumeClaim
		request          reconcile.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should return an error when no name found",
			fields: fields{
				objs:   []runtime.Object{&dbInstance},
				scheme: scheme.Scheme,
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      dbInstance.Name,
						Namespace: dbInstance.Namespace,
					},
				},
				deploymentStatus: &appsv1.Deployment{},
				serviceStatus:    &corev1.Service{},
				pvcStatus:        &corev1.PersistentVolumeClaim{},
			},
			wantErr: true,
		},
		{
			name: "Should update status",
			fields: fields{
				objs:   []runtime.Object{&dbInstance},
				scheme: scheme.Scheme,
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      dbInstance.Name,
						Namespace: dbInstance.Namespace,
					},
				},
				deploymentStatus: &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name: "DeploymentName",
					},
				},
				serviceStatus: &corev1.Service{},
				pvcStatus:     &corev1.PersistentVolumeClaim{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			reqLogger := log.WithValues("Request.Namespace", tt.args.request.Namespace, "Request.Name", tt.args.request.Name)

			if err := r.updateDBStatus(reqLogger, tt.args.deploymentStatus, tt.args.serviceStatus, tt.args.pvcStatus, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("TestReconcilePostgreSQL_UpdateStatus.updateDBStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReconcilePostgreSQL_UpdateDeploymentStatus(t *testing.T) {
	type fields struct {
		scheme *runtime.Scheme
		objs   []runtime.Object
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
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&dbInstance},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      dbInstance.Name,
						Namespace: dbInstance.Namespace,
					},
				},
			},
			wantErr: true,
			want:    reflect.TypeOf(&appsv1.Deployment{}),
		},
		{
			name: "Should find the deployment",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&dbInstance, &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name: "postgresql",
						Namespace: "postgresql",
					},
				}},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      dbInstance.Name,
						Namespace: dbInstance.Namespace,
					},
				},
			},
			wantErr: false,
			want:    reflect.TypeOf(&appsv1.Deployment{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			reqLogger := log.WithValues("Request.Namespace", tt.args.request.Namespace, "Request.Name", tt.args.request.Name)

			got, err := r.updateDeploymentStatus(reqLogger, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcilePostgreSQL_UpdateDeploymentStatus.updateDeploymentStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("TestReconcilePostgreSQL_UpdateDeploymentStatus.updateDeploymentStatus() = %v, want %v", gotType, tt.want)
			}
		})
	}
}

func TestReconcilePostgreSQL_UpdateServiceStatus(t *testing.T) {
	type fields struct {
		scheme *runtime.Scheme
		objs   []runtime.Object
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
			name: "Should not find the service",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&dbInstance},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      dbInstance.Name,
						Namespace: dbInstance.Namespace,
					},
				},
			},
			wantErr: true,
			want:    reflect.TypeOf(&corev1.Service{}),
		},
		{
			name: "Should find the service",
			fields: fields{
				scheme: scheme.Scheme,
				objs: []runtime.Object{&dbInstance, &corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "postgresql",
						Namespace: "postgresql",
					},
				}},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      dbInstance.Name,
						Namespace: dbInstance.Namespace,
					},
				},
			},
			wantErr: false,
			want:    reflect.TypeOf(&corev1.Service{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			reqLogger := log.WithValues("Request.Namespace", tt.args.request.Namespace, "Request.Name", tt.args.request.Name)

			got, err := r.updateServiceStatus(reqLogger, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcilePostgreSQL_UpdateServiceStatus.updateServiceStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("TestReconcilePostgreSQL_UpdateServiceStatus.updateServiceStatus() = %v, want %v", gotType, tt.want)
			}
		})
	}
}


func TestReconcilePostgreSQL_UpdatePVCStatus(t *testing.T) {
	type fields struct {
		scheme *runtime.Scheme
		objs   []runtime.Object
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
			name: "Should not find the pvc",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&dbInstance},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      dbInstance.Name,
						Namespace: dbInstance.Namespace,
					},
				},
			},
			wantErr: true,
			want:    reflect.TypeOf(&corev1.PersistentVolumeClaim{}),
		},
		{
			name: "Should find the pvc",
			fields: fields{
				scheme: scheme.Scheme,
				objs: []runtime.Object{&dbInstance, &corev1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "postgresql",
						Namespace: "postgresql",
					},
				}},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      dbInstance.Name,
						Namespace: dbInstance.Namespace,
					},
				},
			},
			wantErr: false,
			want:    reflect.TypeOf(&corev1.PersistentVolumeClaim{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			reqLogger := log.WithValues("Request.Namespace", tt.args.request.Namespace, "Request.Name", tt.args.request.Name)

			got, err := r.updatePvcStatus(reqLogger, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcilePostgreSQL_UpdatePVCStatus.updatePvcStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("TestReconcilePostgreSQL_UpdatePVCStatus.updatePvcStatus() = %v, want %v", gotType, tt.want)
			}
		})
	}
}

