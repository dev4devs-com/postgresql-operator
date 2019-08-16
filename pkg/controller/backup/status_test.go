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
		dbSecret      *corev1.Secret
		awsSecret     *corev1.Secret
		dbPod         *corev1.Pod
		dbService     *corev1.Service
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
				objs:   []runtime.Object{&bkpInstance},
				scheme: scheme.Scheme,
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
				cronJobStatus: &v1beta1.CronJob{},
				dbSecret:      &corev1.Secret{},
				awsSecret:     &corev1.Secret{},
				dbPod:         &corev1.Pod{},
				dbService:     &corev1.Service{},
			},

			wantErr: true,
		},
		{
			name: "Should update status without enc secret",
			fields: fields{
				objs:   []runtime.Object{&bkpInstance},
				scheme: scheme.Scheme,
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
				cronJobStatus: &v1beta1.CronJob{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
				dbSecret: &corev1.Secret{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
				awsSecret: &corev1.Secret{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
				dbPod: &corev1.Pod{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
				dbService: &corev1.Service{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
			},
			wantErr: false,
		},
		{
			name: "Should return error when not found secret by name",
			fields: fields{
				objs:   []runtime.Object{&bkpInstanceWithSecretNames},
				scheme: scheme.Scheme,
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
				cronJobStatus: &v1beta1.CronJob{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
				dbSecret: &corev1.Secret{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
				dbPod: &corev1.Pod{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
				dbService: &corev1.Service{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			if err := r.updateBackupStatus(tt.args.cronJobStatus, tt.args.dbSecret, tt.args.awsSecret, tt.args.dbPod, tt.args.dbService, tt.args.request); (err != nil) != tt.wantErr {
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
		want    reflect.Type
		wantErr bool
	}{
		{
			name: "Should not find the cronJob",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstance},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
			},
			wantErr: true,
			want:    reflect.TypeOf(&v1beta1.CronJob{}),
		},
		{
			name: "Should find the cronJob",
			fields: fields{
				scheme: scheme.Scheme,
				objs: []runtime.Object{&bkpInstance, &v1beta1.CronJob{
					ObjectMeta: v1.ObjectMeta{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				}},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
			},
			wantErr: false,
			want:    reflect.TypeOf(&v1beta1.CronJob{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			got, err := r.updateCronJobStatus(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestUpdateCronJobStatus error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("TestUpdateCronJobStatus got = %v, want %v", gotType, tt.want)
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
			name: "Should not find the AWSSecret since it was not created",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstance},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
			},
			wantErr: true,
			want:    reflect.TypeOf(&corev1.Secret{}),
		},
		{
			name: "Should not find the AWSSecret",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstanceWithSecretNames},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
			},
			wantErr: true,
			want:    reflect.TypeOf(&corev1.Secret{}),
		},
		{
			name: "Should find the AWSSecret",
			fields: fields{
				scheme: scheme.Scheme,
				objs: []runtime.Object{&bkpInstanceWithSecretNames, &corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name:      bkpInstanceWithSecretNames.Spec.AwsCredentialsSecretName,
						Namespace: bkpInstance.Namespace,
					},
				}},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
			},
			wantErr: false,
			want:    reflect.TypeOf(&corev1.Secret{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			got, err := r.updateAWSSecretStatus(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestUpdateAWSSecretStatus error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("TestUpdateAWSSecretStatus got = %v, want %v", gotType, tt.want)
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
			name: "Should not return error since has not enc secret configured",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstance},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
			},
			wantErr: false,
			want:    reflect.TypeOf(&corev1.Secret{}),
		},
		{
			name: "Should not find the encSecret",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstanceWithSecretNames},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
			},
			wantErr: true,
			want:    reflect.TypeOf(&corev1.Secret{}),
		},
		{
			name: "Should find the encSecret",
			fields: fields{
				scheme: scheme.Scheme,
				objs: []runtime.Object{&bkpInstanceWithSecretNames, &corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name:      bkpInstanceWithSecretNames.Spec.EncryptionKeySecretName,
						Namespace: bkpInstance.Namespace,
					},
				}},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
			},
			wantErr: false,
			want:    reflect.TypeOf(&corev1.Secret{}),
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
				objs:   []runtime.Object{&bkpInstance},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
			},
			wantErr: true,
			want:    reflect.TypeOf(&corev1.Secret{}),
		},
		{
			name: "Should find the dbSecret",
			fields: fields{
				scheme: scheme.Scheme,
				objs: []runtime.Object{&bkpInstance, &corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name:      "db-postgresql-backup",
						Namespace: bkpInstance.Namespace,
					},
				}},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
			},
			wantErr: false,
			want:    reflect.TypeOf(&corev1.Secret{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			got, err := r.updateDBSecretStatus(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestUpdateDBSecretStatus error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("TestUpdateDBSecretStatus got = %v, want %v", gotType, tt.want)
			}
		})
	}
}

func TestUpdatePodDatabaseFoundStatus(t *testing.T) {
	type fields struct {
		scheme *runtime.Scheme
		objs   []runtime.Object
	}
	type args struct {
		request reconcile.Request
		pod     corev1.Pod
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    reflect.Type
		wantErr bool
	}{
		{
			name: "Should not find the Pod",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstance},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
				pod: corev1.Pod{},
			},
			wantErr: false,
		},
		{
			name: "Should find the Pod",
			fields: fields{
				scheme: scheme.Scheme,
				objs: []runtime.Object{&bkpInstance, &corev1.Pod{
					ObjectMeta: v1.ObjectMeta{
						Name:      "postgresql",
						Namespace: "postgresql",
					},
				}},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
			},
			wantErr: false,
			want:    reflect.TypeOf(&corev1.Pod{}),
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
		want    reflect.Type
		wantErr bool
	}{
		{
			name: "Should not find the Service",
			fields: fields{
				scheme: scheme.Scheme,
				objs:   []runtime.Object{&bkpInstance},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
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
				objs: []runtime.Object{&bkpInstance, &corev1.Service{
					ObjectMeta: v1.ObjectMeta{
						Name:      "postgresql",
						Namespace: "postgresql",
					},
				}},
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
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

			err := r.updateServiceDatabaseFoundStatus(tt.args.request, &tt.args.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestUpdateServiceDatabaseFoundStatus error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
