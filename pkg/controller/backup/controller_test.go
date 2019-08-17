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

func TestReconcileBackup(t *testing.T) {
	type fields struct {
		objs []runtime.Object
	}
	type args struct {
		bkpInstance v1alpha1.Backup
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantRequeue   bool
		wantAwsSecret bool
		wantDBSecret  bool
		wantEncSecret bool
		wantCronJob   bool
		wantErr       bool
	}{
		{
			name: "Should work with default values as key of env vars variables",
			fields: fields{
				objs: []runtime.Object{&bkpInstanceWithMandatorySpec, &dbInstanceWithConfigMap, &podDatabaseConfigMap, &serviceDatabase, &configMapDefault},
			},
			args: args{
				bkpInstance: bkpInstanceWithMandatorySpec,
			},
			wantErr:       false,
			wantRequeue:   false,
			wantAwsSecret: true,
			wantDBSecret:  true,
			wantEncSecret: false,
			wantCronJob:   true,
		},
		{
			name: "Should fail with wrong key values mapped",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithMandatorySpec,
					&dbInstanceWithConfigMap,
					&podDatabaseConfigMap,
					&serviceDatabase,
					&configMapOtherKeyValuesInvalidKeys,
				},
			},
			args: args{
				bkpInstance: bkpInstanceWithMandatorySpec,
			},
			wantErr:       true,
			wantRequeue:   false,
			wantAwsSecret: false,
			wantDBSecret:  false,
			wantEncSecret: false,
			wantCronJob:   false,
		},
		{
			name: "Should work with customized keys for the db env vars",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithMandatorySpec,
					&dbInstanceWithConfigMapAndCustomizeKeys,
					&podDatabaseConfigMap,
					&serviceDatabase,
					&configMapOtherKeyValuesInvalidKeys,
				},
			},
			args: args{
				bkpInstance: bkpInstanceWithMandatorySpec,
			},
			wantErr:       false,
			wantRequeue:   false,
			wantAwsSecret: true,
			wantDBSecret:  true,
			wantEncSecret: false,
			wantCronJob:   true,
		},
		{
			name: "Should work with encryption secret data and create this secret",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithEncSecretData,
					&dbInstanceWithoutSpec,
					&podDatabase,
					&serviceDatabase,
				},
			},
			args: args{
				bkpInstance: bkpInstanceWithEncSecretData,
			},

			wantErr:       false,
			wantRequeue:   false,
			wantAwsSecret: true,
			wantDBSecret:  true,
			wantEncSecret: true,
			wantCronJob:   true,
		},
		{
			name: "Should work with secret names and found the secrets",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithSecretNames,
					&dbInstanceWithoutSpec,
					&podDatabase,
					&serviceDatabase,
					&awsSecretWithSecretNames,
					&encSecretWithSecretNames,
				},
			},
			args: args{
				bkpInstance: bkpInstanceWithSecretNames,
			},

			wantErr:       false,
			wantRequeue:   false,
			wantAwsSecret: true,
			wantDBSecret:  true,
			wantEncSecret: true,
			wantCronJob:   true,
		},
		{
			name: "Should fail when it is missing the pod database",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithMandatorySpec,
					&dbInstanceWithoutSpec,
					&serviceDatabase,
				},
			},
			args: args{
				bkpInstance: bkpInstanceWithMandatorySpec,
			},

			wantErr:       true,
			wantRequeue:   false,
			wantAwsSecret: false,
			wantDBSecret:  false,
			wantEncSecret: false,
			wantCronJob:   false,
		},
		{
			name: "Should fail when it is missing the service from database",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithMandatorySpec,
					&dbInstanceWithoutSpec,
					&podDatabase,
				},
			},
			args: args{
				bkpInstance: bkpInstanceWithMandatorySpec,
			},

			wantErr:       true,
			wantRequeue:   false,
			wantAwsSecret: false,
			wantDBSecret:  false,
			wantEncSecret: false,
			wantCronJob:   false,
		},
		{
			name: "Should fail when it is missing the service from database",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithMandatorySpec,
					&dbInstanceWithoutSpec,
					&podDatabase,
				},
			},
			args: args{
				bkpInstance: bkpInstanceWithMandatorySpec,
			},

			wantErr:       true,
			wantRequeue:   false,
			wantAwsSecret: false,
			wantDBSecret:  false,
			wantEncSecret: false,
			wantCronJob:   false,
		}, {
			name: "Should fail since has database CR was not applied",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithMandatorySpec,
				},
			},
			args: args{
				bkpInstance: bkpInstanceWithMandatorySpec,
			},

			wantErr:       true,
			wantRequeue:   false,
			wantAwsSecret: false,
			wantDBSecret:  false,
			wantEncSecret: false,
			wantCronJob:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs)

			// mock request to simulate Reconcile() being called on an event for a watched resource
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      tt.args.bkpInstance.Name,
					Namespace: tt.args.bkpInstance.Namespace,
				},
			}

			res, err := r.Reconcile(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcileBackup reconcile: error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			awsSecret := &corev1.Secret{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: getAWSSecretName(&tt.args.bkpInstance), Namespace: getAwsSecretNamespace(&tt.args.bkpInstance)}, awsSecret)
			if (err == nil) != tt.wantAwsSecret {
				t.Errorf("TestReconcileBackup to get aws secret error = %v, wantErr %v", err, tt.wantAwsSecret)
				return
			}

			dbSecret := &corev1.Secret{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: dbSecretPrefix + tt.args.bkpInstance.Name, Namespace: tt.args.bkpInstance.Namespace}, dbSecret)
			if (err == nil) != tt.wantDBSecret {
				t.Errorf("TestReconcileBackup to get db secret error = %v, wantErr %v", err, tt.wantDBSecret)
				return
			}

			encSecret := &corev1.Secret{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: getEncSecretName(&tt.args.bkpInstance), Namespace: getEncSecretNamespace(&tt.args.bkpInstance)}, encSecret)
			if (err == nil) != tt.wantEncSecret {
				t.Errorf("TestReconcileBackup to get enc secret error = %v, wantErr %v", err, tt.wantEncSecret)
				return
			}

			cronJob := &v1beta1.CronJob{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: tt.args.bkpInstance.Name, Namespace: tt.args.bkpInstance.Namespace}, cronJob)
			if (err == nil) != tt.wantCronJob {
				t.Errorf("TestReconcileBackup to get cronjob error = %v, wantErr %v", err, tt.wantCronJob)
				return
			}

			if (res.Requeue) != tt.wantRequeue {
				t.Errorf("TestReconcileBackup expect request to requeue res.Requeue = %v, wantRequeue %v", res.Requeue, tt.wantRequeue)
				return
			}
		})
	}
}
