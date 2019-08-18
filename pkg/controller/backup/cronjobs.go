package backup

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//Returns the NewCronJob object for the PostgreSQL Backup
func buildCronJob(bkp *v1alpha1.Backup, scheme *runtime.Scheme) *v1beta1.CronJob {
	cron := &v1beta1.CronJob{
		ObjectMeta: v1.ObjectMeta{
			Name:      bkp.Name,
			Namespace: bkp.Namespace,
			Labels:    getBkpLabels(bkp.Name),
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
											Value: getAWSSecretName(bkp),
										},
										{
											Name:  "BACKEND_SECRET_NAMESPACE",
											Value: getAwsSecretNamespace(bkp),
										},
										{
											Name:  "ENCRYPTION_SECRET_NAME",
											Value: getEncSecretName(bkp),
										},
										{
											Name:  "ENCRYPTION_SECRET_NAMESPACE",
											Value: getEncSecretNamespace(bkp),
										},
										{
											Name:  "COMPONENT_SECRET_NAME",
											Value: dbSecretPrefix + bkp.Name,
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
	// Set PostgreSQL db as the owner and controller
	controllerutil.SetControllerReference(bkp, cron, scheme)
	return cron
}
