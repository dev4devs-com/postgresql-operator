package backup

import (
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func TestUpdateBackupStatus(t *testing.T) {
	type fields struct {
		objs   []runtime.Object
		scheme *runtime.Scheme
	}
	type args struct {
		cronJobStatus *v1beta1.CronJob
		request       reconcile.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should update status without enc secret",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithMandatorySpec,
					&awsSecretWithMadatorySpec,
					&cronJobWithMadatorySpec,
					&dbSecretWithMadatorySpec,
				},
				scheme: scheme.Scheme,
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithMandatorySpec.Name,
						Namespace: bkpInstanceWithMandatorySpec.Namespace,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should update status with enc secret",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithSecretNames,
					&awsSecretWithSecretNames,
					&croJobWithSecretNames,
					&encSecretWithSecretNames,
					&dbSecretWithSecretNames,
				},
				scheme: scheme.Scheme,
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithSecretNames.Name,
						Namespace: bkpInstanceWithSecretNames.Namespace,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should return error when not found aws secret",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithMandatorySpec,
					&cronJobWithMadatorySpec,
				},
				scheme: scheme.Scheme,
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithMandatorySpec.Name,
						Namespace: bkpInstanceWithMandatorySpec.Namespace,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Should return error when not found the cronjob",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithMandatorySpec,
					&awsSecretWithMadatorySpec,
				},
				scheme: scheme.Scheme,
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithMandatorySpec.Name,
						Namespace: bkpInstanceWithMandatorySpec.Namespace,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			if err := r.updateBackupStatus(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("TestUpdateBackupStatus error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateCronJobStatus(t *testing.T) {
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
		wantErr bool
	}{
		{
			name: "Should not find the cronJob",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstanceWithMandatorySpec},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithMandatorySpec.Name,
						Namespace: bkpInstanceWithMandatorySpec.Namespace,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Should find the cronJob",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstanceWithMandatorySpec, &cronJobWithMadatorySpec},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithMandatorySpec.Name,
						Namespace: bkpInstanceWithMandatorySpec.Namespace,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			err := r.updateCronJobStatus(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestUpdateCronJobStatus error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpdateAWSSecretStatus(t *testing.T) {
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
			name: "Should fail since the aws secret was not found",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstanceWithMandatorySpec},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithMandatorySpec.Name,
						Namespace: bkpInstanceWithMandatorySpec.Namespace,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Should works with success",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstanceWithMandatorySpec, &awsSecretWithMadatorySpec},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithMandatorySpec.Name,
						Namespace: bkpInstanceWithMandatorySpec.Namespace,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should works with success when the aws secret is informed by the user",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstanceWithSecretNames, &awsSecretWithSecretNames},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithSecretNames.Name,
						Namespace: bkpInstanceWithSecretNames.Namespace,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			err := r.updateAWSSecretStatus(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestUpdateAWSSecretStatus error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpdateEncSecretStatus(t *testing.T) {
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
			name: "Should fail since the enc secret was not found",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstanceWithSecretNames},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithSecretNames.Name,
						Namespace: bkpInstanceWithSecretNames.Namespace,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Should works with success since the Backup CR was not customized to use secret",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstanceWithMandatorySpec},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithMandatorySpec.Name,
						Namespace: bkpInstanceWithMandatorySpec.Namespace,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should works with success when the enc secret is informed by the user",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstanceWithSecretNames, &encSecretWithSecretNames},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithSecretNames.Name,
						Namespace: bkpInstanceWithSecretNames.Namespace,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			err := r.updateEncSecretStatus(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestUpdateEncSecretStatus error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpdateDBSecretStatus(t *testing.T) {
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
			name: "Should not find the dbSecret",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstanceWithMandatorySpec},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithMandatorySpec.Name,
						Namespace: bkpInstanceWithMandatorySpec.Namespace,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Should find the dbSecret",
			fields: fields{
				scheme: scheme.Scheme,
				objs: []runtime.Object{&bkpInstanceWithMandatorySpec, &corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name:      "db-postgresql-backup",
						Namespace: bkpInstanceWithMandatorySpec.Namespace,
					},
				}},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithMandatorySpec.Name,
						Namespace: bkpInstanceWithMandatorySpec.Namespace,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			err := r.updateDBSecretStatus(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestUpdateDBSecretStatus error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpdatePodDatabaseFoundStatus(t *testing.T) {
	type fields struct {
		objs []runtime.Object
	}
	type args struct {
		request reconcile.Request
		pod     corev1.Pod
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should not find the Pod",
			fields: fields{
				objs: []runtime.Object{&bkpInstanceWithMandatorySpec},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithMandatorySpec.Name,
						Namespace: bkpInstanceWithMandatorySpec.Namespace,
					},
				},
				pod: corev1.Pod{},
			},
			wantErr: false,
		},
		{
			name: "Should find the Pod",
			fields: fields{
				objs: []runtime.Object{&bkpInstanceWithMandatorySpec, &corev1.Pod{
					ObjectMeta: v1.ObjectMeta{
						Name:      "postgresql",
						Namespace: "postgresql",
					},
				}},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithMandatorySpec.Name,
						Namespace: bkpInstanceWithMandatorySpec.Namespace,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			err := r.updatePodDatabaseFoundStatus(tt.args.request, &tt.args.pod)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestUpdatePodDatabaseFoundStatus error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestUpdateServiceDatabaseFoundStatus(t *testing.T) {
	type fields struct {
		scheme *runtime.Scheme
		objs   []runtime.Object
	}
	type args struct {
		request reconcile.Request
		service corev1.Service
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should not find the Service",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstanceWithMandatorySpec},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithMandatorySpec.Name,
						Namespace: bkpInstanceWithMandatorySpec.Namespace,
					},
				},
				service: corev1.Service{},
			},
			wantErr: false,
		},
		{
			name: "Should find the Service",
			fields: fields{
				scheme: scheme.Scheme,
				objs: []runtime.Object{&bkpInstanceWithMandatorySpec, &corev1.Service{
					ObjectMeta: v1.ObjectMeta{
						Name:      "postgresql",
						Namespace: "postgresql",
					},
				}},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstanceWithMandatorySpec.Name,
						Namespace: bkpInstanceWithMandatorySpec.Namespace,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			err := r.updateServiceDatabaseFoundStatus(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestUpdateServiceDatabaseFoundStatus error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
