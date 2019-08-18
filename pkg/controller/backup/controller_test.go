package backup

import (
	"context"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/utils"
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
			name: "Should fail with wrong database name key mapped when it will build the db data secret",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithMandatorySpec,
					&dbInstanceWithConfigMap,
					&podDatabaseConfigMap,
					&serviceDatabase,
					&configMapInvalidDatabaseKey,
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
			name: "Should fail with wrong database user key mapped when it will build the db data secret",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithMandatorySpec,
					&dbInstanceWithConfigMap,
					&podDatabaseConfigMap,
					&serviceDatabase,
					&configMapInvalidUserKey,
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
			name: "Should fail with wrong database pwd key mapped when it will build the db data secret",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithMandatorySpec,
					&dbInstanceWithConfigMap,
					&podDatabaseConfigMap,
					&serviceDatabase,
					&configMapInvalidPwdKey,
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
					&configMapInvalidDatabaseKey,
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
			name: "Should fail when the aws secret informed by the user do not exist",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithSecretNames,
					&dbInstanceWithoutSpec,
					&podDatabase,
					&serviceDatabase,
					&encSecretWithSecretNames,
				},
			},
			args: args{
				bkpInstance: bkpInstanceWithSecretNames,
			},

			wantErr:       true,
			wantRequeue:   false,
			wantAwsSecret: false,
			wantDBSecret:  true,
			wantEncSecret: true,
			wantCronJob:   true,
		},
		{
			name: "Should fail when the enc secret informed by the user do not exist",
			fields: fields{
				objs: []runtime.Object{
					&bkpInstanceWithSecretNames,
					&dbInstanceWithoutSpec,
					&podDatabase,
					&serviceDatabase,
					&awsSecretWithSecretNames,
				},
			},
			args: args{
				bkpInstance: bkpInstanceWithSecretNames,
			},

			wantErr:       true,
			wantRequeue:   false,
			wantAwsSecret: true,
			wantDBSecret:  true,
			wantEncSecret: false,
			wantCronJob:   false,
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
		{
			name: "Should fail because is missing the PostgreSQL CR",
			fields: fields{
				objs: []runtime.Object{&bkpInstanceWithMandatorySpec, &podDatabaseConfigMap, &serviceDatabase, &configMapDefault},
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
			name: "Should fail because is missing the Backup CR",
			fields: fields{
				objs: []runtime.Object{&dbInstanceWithoutSpec, &podDatabaseConfigMap, &serviceDatabase, &configMapDefault},
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
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: utils.GetAWSSecretName(&tt.args.bkpInstance), Namespace: utils.GetAwsSecretNamespace(&tt.args.bkpInstance)}, awsSecret)
			if (err == nil) != tt.wantAwsSecret {
				t.Errorf("TestReconcileBackup to get aws secret error = %v, wantErr %v", err, tt.wantAwsSecret)
				return
			}

			dbSecret := &corev1.Secret{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: utils.DbSecretPrefix + tt.args.bkpInstance.Name, Namespace: tt.args.bkpInstance.Namespace}, dbSecret)
			if (err == nil) != tt.wantDBSecret {
				t.Errorf("TestReconcileBackup to get db secret error = %v, wantErr %v", err, tt.wantDBSecret)
				return
			}

			encSecret := &corev1.Secret{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: utils.GetEncSecretName(&tt.args.bkpInstance), Namespace: utils.GetEncSecretNamespace(&tt.args.bkpInstance)}, encSecret)
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
