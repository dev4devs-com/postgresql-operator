package postgresql

import (
	"context"
	"github.com/dev4devs-com/postgresql-operator/pkg/config"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcilePostgreSQL_Success(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&dbInstance,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      dbInstance.Name,
			Namespace: dbInstance.Namespace,
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	deployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), req.NamespacedName, deployment)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	service := &corev1.Service{}
	err = r.client.Get(context.TODO(), req.NamespacedName, service)
	if err != nil {
		t.Fatalf("get service: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	pvc := &corev1.PersistentVolumeClaim{}
	err = r.client.Get(context.TODO(), req.NamespacedName, pvc)
	if err != nil {
		t.Fatalf("get pvc: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}

func TestReconcilePostgreSQL_NotFound(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&dbInstance,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      dbInstance.Name,
			Namespace: "unknown",
		},
	}

	res, err := r.Reconcile(req)
	if err == nil {
		t.Error("should fail since the instance do not exist in the <unknown> nammespace")
	}

	if res.Requeue {
		t.Fatalf("did not expected reconcile to requeue.")
	}
}

func TestReconcilePostgreSQL_UsingConfigMapToCreateEnvVars(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&dbInstanceConfigMapSameKeys,
		&configMapSameKeyValues,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      dbInstance.Name,
			Namespace: dbInstance.Namespace,
		},
	}

	_ = r.client.Create(context.TODO(), &configMapSameKeyValues)

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	deployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), req.NamespacedName, deployment)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom == nil {
		t.Error("deployment envvar did not came from service instance config map")
	}

	if deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.ConfigMapKeyRef.Name != configMapSameKeyValues.Name {
		t.Fatalf("deployment envvar did not came from service instance config map: (%v,%v)", deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.ConfigMapKeyRef.Name, configMapSameKeyValues.Name)
	}

	if res.Requeue {
		t.Error("reconcile did not requeue request as expected")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	service := &corev1.Service{}
	err = r.client.Get(context.TODO(), req.NamespacedName, service)
	if err != nil {
		t.Fatalf("get service: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	pvc := &corev1.PersistentVolumeClaim{}
	err = r.client.Get(context.TODO(), req.NamespacedName, pvc)
	if err != nil {
		t.Fatalf("get pvc: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}

func TestReconcilePostgreSQ_ReplicasSizes(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&dbInstance,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      dbInstance.Name,
			Namespace: dbInstance.Namespace,
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	deployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), req.NamespacedName, deployment)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	//Mock Replicas wrong size
	size := int32(3)
	deployment.Spec.Replicas = &size

	// Update
	err = r.client.Update(context.TODO(), deployment)
	if err != nil {
		t.Fatalf("fails when ttry to update deployment replicas: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	deployment = &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), req.NamespacedName, deployment)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	if *deployment.Spec.Replicas != dbInstance.Spec.Size {
		t.Error("Replicas size was not respected")
	}
}

func TestReconcilePostgreSQL_Reconcile_InstanceWithoutSpec(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&dbInstanceWithoutSpec,
	}

	r := buildReconcileWithFakeClientWithMocks(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      dbInstanceWithoutSpec.Name,
			Namespace: dbInstanceWithoutSpec.Namespace,
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	deployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), req.NamespacedName, deployment)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	// Check if the quantity of Replicas for this deployment is equals the specification
	if *deployment.Spec.Replicas != config.NewPostgreSQLConfig().Size {
		t.Errorf("dep size (%d) is not the expected size (%d)", deployment.Spec.Replicas, config.NewPostgreSQLConfig().Size)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	service := &corev1.Service{}
	err = r.client.Get(context.TODO(), req.NamespacedName, service)
	if err != nil {
		t.Fatalf("get service: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	pvc := &corev1.PersistentVolumeClaim{}
	err = r.client.Get(context.TODO(), req.NamespacedName, pvc)
	if err != nil {
		t.Fatalf("get pvc: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}
