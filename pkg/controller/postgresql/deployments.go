package postgresql

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//buildDBDeployment returns the deployment object for the PostgreSQL
func (r *ReconcilePostgresql) buildDBDeployment(db *v1alpha1.Postgresql) *appsv1.Deployment {
	ls := getDBLabels(db.Name)
	auto := true
	replicas := db.Spec.Size
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      db.Name,
			Namespace: db.Namespace,
			Labels:    ls,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           db.Spec.Image,
						Name:            db.Spec.ContainerName,
						ImagePullPolicy: db.Spec.ContainerImagePullPolicy,
						Ports: []corev1.ContainerPort{{
							ContainerPort: db.Spec.DatabasePort,
							Protocol:      "TCP",
						}},
						Env: []corev1.EnvVar{
							buildDatabaseNameEnvVar(db),
							buildDatabaseUserEnvVar(db),
							buildDatabasePasswordEnvVar(db),
							{
								Name:  "PGDATA",
								Value: "/var/lib/pgsql/data/pgdata",
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      db.Name,
								MountPath: "/var/lib/pgsql/data",
							},
						},
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								Exec: &corev1.ExecAction{
									Command: []string{
										"/usr/libexec/check-container",
										"'--live'",
									},
								},
							},
							FailureThreshold:    3,
							InitialDelaySeconds: 120,
							PeriodSeconds:       10,
							TimeoutSeconds:      10,
							SuccessThreshold:    1,
						},
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								Exec: &corev1.ExecAction{
									Command: []string{
										"/usr/libexec/check-container",
									},
								},
							},
							FailureThreshold:    3,
							InitialDelaySeconds: 5,
							PeriodSeconds:       10,
							TimeoutSeconds:      1,
							SuccessThreshold:    1,
						},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(db.Spec.DatabaseMemoryLimit),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(db.Spec.DatabaseMemoryRequest),
							},
						},
						TerminationMessagePath: "/dev/termination-log",
					}},
					DNSPolicy:     corev1.DNSClusterFirst,
					RestartPolicy: corev1.RestartPolicyAlways,
					Volumes: []corev1.Volume{
						{
							Name: db.Name,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: db.Name,
								},
							},
						},
					},
					AutomountServiceAccountToken: &auto,
				},
			},
		},
	}
	// Set PostgreSQL db as the owner and controller
	controllerutil.SetControllerReference(db, dep, r.scheme)
	return dep
}
