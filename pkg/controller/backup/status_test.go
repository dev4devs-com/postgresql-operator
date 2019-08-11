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

func TestBackup_UpdateBackupStatus(t *testing.T) {
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

			reqLogger := log.WithValues("Request.Namespace", tt.args.request.Namespace, "Request.Name", tt.args.request.Name)

			if err := r.updateBackupStatus(reqLogger, tt.args.cronJobStatus, tt.args.dbSecret, tt.args.awsSecret, tt.args.dbPod, tt.args.dbService, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("TestBackup_UpdateBackupStatus.updateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReconcilePostgreSQL_UpdateCronJobStatus(t *testing.T) {
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

			reqLogger := log.WithValues("Request.Namespace", tt.args.request.Namespace, "Request.Name", tt.args.request.Name)

			got, err := r.updateCronJobStatus(reqLogger, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcilePostgreSQL_UpdateCronJobStatus.updateCronJobStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("TestReconcilePostgreSQL_UpdateCronJobStatus.updateCronJobStatus() = %v, want %v", gotType, tt.want)
			}
		})
	}
}

func TestReconcilePostgreSQL_UpdateAWSSecretStatus(t *testing.T) {
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

			reqLogger := log.WithValues("Request.Namespace", tt.args.request.Namespace, "Request.Name", tt.args.request.Name)

			got, err := r.updateAWSSecretStatus(reqLogger, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcilePostgreSQL_UpdateCronJobStatus.updateAWSSecretStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("TestReconcilePostgreSQL_UpdateCronJobStatus.updateAWSSecretStatus() = %v, want %v", gotType, tt.want)
			}
		})
	}
}

func TestReconcilePostgreSQL_UpdateEncSecretStatus(t *testing.T) {
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

			reqLogger := log.WithValues("Request.Namespace", tt.args.request.Namespace, "Request.Name", tt.args.request.Name)

			err := r.updateEncSecretStatus(reqLogger, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcilePostgreSQL_UpdateEncSecretStatus.UpdateEncSecretStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestReconcilePostgreSQL_UpdateDBSecretStatus(t *testing.T) {
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
						Name:      "aws-postgresql-backup",
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

			reqLogger := log.WithValues("Request.Namespace", tt.args.request.Namespace, "Request.Name", tt.args.request.Name)

			got, err := r.updateAWSSecretStatus(reqLogger, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcilePostgreSQL_UpdateDBSecretStatus.UpdateDBSecretStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("TestReconcilePostgreSQL_UpdateDBSecretStatus.UpdateDBSecretStatus() = %v, want %v", gotType, tt.want)
			}
		})
	}
}

func TestReconcilePostgreSQL_UpdatePodDatabaseFoundStatus(t *testing.T) {
	type fields struct {
		scheme *runtime.Scheme
		objs   []runtime.Object
	}
	type args struct {
		request reconcile.Request
		pod corev1.Pod
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

			reqLogger := log.WithValues("Request.Namespace", tt.args.request.Namespace, "Request.Name", tt.args.request.Name)

			err := r.updatePodDatabaseFoundStatus(reqLogger, tt.args.request, &tt.args.pod)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcilePostgreSQL_UpdatePodDabaseFoundStatus.updatePodDatabaseFoundStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}


func TestReconcilePostgreSQL_UpdateServiceDatabaseFoundStatus(t *testing.T) {
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

			reqLogger := log.WithValues("Request.Namespace", tt.args.request.Namespace, "Request.Name", tt.args.request.Name)

			err := r.updateServiceDatabaseFoundStatus(reqLogger, tt.args.request, &tt.args.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcilePostgreSQL_UpdateServiceDatabaseFoundStatus.updateServiceDatabaseFoundStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}