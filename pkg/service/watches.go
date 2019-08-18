package service

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

//Watch for changes to secondary resource and create the owner Backup

func Watch(c controller.Controller, obj runtime.Object, isConttroller bool, owner runtime.Object) error {
	err := c.Watch(&source.Kind{Type: obj}, &handler.EnqueueRequestForOwner{
		IsController: isConttroller,
		OwnerType:    owner,
	})
	return err
}
