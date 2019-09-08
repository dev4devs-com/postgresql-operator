package resource

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/utils"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//Returns the NewBackupCronJob object for the Database Backup
func NewBackupCronJob(bkp *v1alpha1.Backup, scheme *runtime.Scheme) *v1beta1.CronJob {
	cron := &v1beta1.CronJob{
		ObjectMeta: v1.ObjectMeta{
			Name:      bkp.Name,
			Namespace: bkp.Namespace,
			Labels:    utils.GetLabels(bkp.Name),
		},
		Spec: v1beta1.CronJobSpec{
			Schedule: bkp.Spec.Schedule,
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							ServiceAccountName: "postgresql-operator",
							Containers: []corev1.Container{
								{
									Name:    bkp.Name,
									Image:   bkp.Spec.Image,
									Command: []string{"/opt/intly/tools/entrypoint.sh", "-c", "postgres", "-n", bkp.Namespace, "-b", "s3", "-e", ""},
									Env: []corev1.EnvVar{
										{
											Name:  "BACKEND_SECRET_NAME",
											Value: utils.GetAWSSecretName(bkp),
										},
										{
											Name:  "BACKEND_SECRET_NAMESPACE",
											Value: utils.GetAwsSecretNamespace(bkp),
										},
										{
											Name:  "ENCRYPTION_SECRET_NAME",
											Value: utils.GetEncSecretName(bkp),
										},
										{
											Name:  "ENCRYPTION_SECRET_NAMESPACE",
											Value: utils.GetEncSecretNamespace(bkp),
										},
										{
											Name:  "COMPONENT_SECRET_NAME",
											Value: utils.DbSecretPrefix + bkp.Name,
										},
										{
											Name:  "COMPONENT_SECRET_NAMESPACE",
											Value: bkp.Namespace,
										},
										{
											Name:  "PRODUCT_NAME",
											Value: bkp.Spec.ProductName,
										},
									},
								},
							},
							RestartPolicy: corev1.RestartPolicyOnFailure,
						},
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(bkp, cron, scheme)
	return cron
}
