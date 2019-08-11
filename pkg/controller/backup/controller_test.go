package backup

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	"k8s.io/api/batch/v1beta1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"reflect"
	"testing"
)

func TestReconcileBackup_Update(t *testing.T) {
	type fields struct {
		createdInstance  *v1alpha1.Backup
		instanceToUpdate *v1alpha1.Backup
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
				createdInstance:  &bkpInstance,
				instanceToUpdate: &bkpInstance,
			},
			wantErr: false,
		},
		{
			name: "Should give an error when the namespace is not found",
			fields: fields{
				createdInstance:  &bkpInstance,
				instanceToUpdate: &bkpInstanceNonDefaultNamespace,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.fields.createdInstance}

			r := buildReconcileWithFakeClientWithMocks(objs)

			reqLogger := log.WithValues("Request.Namespace", tt.fields.instanceToUpdate.Namespace, "Request.Name", tt.fields.createdInstance.Name)

			err := r.update(tt.fields.instanceToUpdate, reqLogger)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcileBackup_Update.update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestReconcileBackup_Create(t *testing.T) {
	objs := []runtime.Object{&bkpInstance}
	r := buildReconcileWithFakeClientWithMocks(objs)
	dataDBSecret, _ := r.buildDBSecretData(&bkpInstance)
	awsDataSecret := buildAwsSecretData(&bkpInstance)
	type fields struct {
		scheme *runtime.Scheme
	}
	type args struct {
		instance   *v1alpha1.Backup
		kind       string
		secretData map[string][]byte
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "Should create and return a new CronJob",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       CronJob,
				secretData: nil,
			},
			wantErr: false,
		},
		{
			name: "Should create and return a new DB Secret",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       DBSecret,
				secretData: dataDBSecret,
			},
			wantErr: false,
		},
		{
			name: "Should create and return a new Aws CronJob",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       AwsSecret,
				secretData: awsDataSecret,
			},
			wantErr: false,
		},
		{
			name: "Should fail to create an unknown type",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance: &bkpInstance,
				kind:     "UNKNOWN",
			},
			wantErr:   false,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			reqLogger := log.WithValues("Request.Namespace", tt.args.instance.Namespace, "Request.Name", tt.args.instance.Name)

			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("TestReconcileBackup_Create.create() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			err := r.create(tt.args.instance, tt.args.kind, reqLogger)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReconcileBackup_Create.create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestReconcileBackup_BuildFactory(t *testing.T) {
	objs := []runtime.Object{&bkpInstance}
	r := buildReconcileWithFakeClientWithMocks(objs)
	dataDBSecret, _ := r.buildDBSecretData(&bkpInstance)
	awsDataSecret := buildAwsSecretData(&bkpInstance)

	type fields struct {
		scheme *runtime.Scheme
	}
	type args struct {
		instance   *v1alpha1.Backup
		kind       string
		secretData map[string][]byte
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      reflect.Type
		wantPanic bool
	}{
		{
			name: "Should create a CronJob",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       CronJob,
				secretData: nil,
			},
			want: reflect.TypeOf(&v1beta1.CronJob{}),
		},
		{
			name: "Should create a DB Secret",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       DBSecret,
				secretData: dataDBSecret,
			},
			want: reflect.TypeOf(&v1.Secret{}),
		},
		{
			name: "Should create a Aws Secret",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       AwsSecret,
				secretData: awsDataSecret,
			},
			want: reflect.TypeOf(&v1.Secret{}),
		},
		{
			name: "Should panic when trying to create unrecognized object type",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance: &bkpInstance,
				kind:     "UNDEFINED",
			},
			wantPanic: true,
			want:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqLogger := log.WithValues("Request.Namespace", tt.args.instance.Namespace, "Request.Name", tt.args.instance.Name)

			// testing if the panic() function is executed
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("TestReconcileBackup_BuildFactory.buildFactory() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			got, _ := r.buildFactory(tt.args.instance, tt.args.kind, reqLogger)

			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("TestReconcileBackup_BuildFactory.buildFactory() = %v, want %v", gotType, tt.want)
			}
		})
	}
}
