package service

import (
	goctx "context"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//FetchService returns the Service resource with the name in the namespace
func FetchService(name, namespace string, client client.Client) (*corev1.Service, error) {
	service := &corev1.Service{}
	err := client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, service)
	return service, err
}

//FetchService returns the Deployment resource with the name in the namespace
func FetchDeployment(name, namespace string, client client.Client) (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{}
	err := client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, deployment)
	return deployment, err
}

//FetchPersistentVolumeClaim returns the PersistentVolumeClaim resource with the name in the namespace
func FetchPersistentVolumeClaim(name, namespace string, client client.Client) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{}
	err := client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, pvc)
	return pvc, err
}

//FetchCronJob returns the CronJob resource with the name in the namespace
func FetchCronJob(name, namespace string, client client.Client) (*v1beta1.CronJob, error) {
	cronJob := &v1beta1.CronJob{}
	err := client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, cronJob)
	return cronJob, err
}

//FetchSecret returns the Secret resource with the name in the namespace
func FetchSecret(namespace, name string, client client.Client) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, secret)
	return secret, err
}

//FetchSecret returns the ConfigMap resource with the name in the namespace
func FetchConfigMap(name, namespace string, client client.Client) (*corev1.ConfigMap, error) {
	cfg := &corev1.ConfigMap{}
	err := client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, cfg)
	return cfg, err
}
