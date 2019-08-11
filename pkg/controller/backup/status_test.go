package backup

import (
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
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
