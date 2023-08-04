package main

import (
	"fmt"
	"time"

	mar "go-angi/api/myapigroup/v1alpha1"
	"go-angi/pkg/controller"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	controllerAgentName = "myapp-controller"
)

func main() {

	clientset, dynamicClient := controller.InitKubernetesClient()

	// Create the informer for the custom resource
	myAppResourceGVR := schema.GroupVersionResource{
		Group:    "my.api.group",   // Replace with your custom resource's API group
		Version:  "v1alpha1",       // Replace with your custom resource's version
		Resource: "myappresources", // Replace with your custom resource's plural form
	}

	informerFactory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynamicClient, 10*time.Second, "default", nil)
	// Create the informer for the custom resource
	myAppResourceInformer := informerFactory.ForResource(myAppResourceGVR).Informer()

	stopCh := make(chan struct{})
	defer close(stopCh)

	// Start the informer to begin watching MyAppResource objects
	go myAppResourceInformer.Run(stopCh)

	// Define the callback functions for add, update, and delete events
	myAppResourceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("running myAppResource add callback")
			myAppResource, err := toMyAppResource(obj)
			if err != nil {
				return
			}
			controller.HandleMyAppResourceAdd(myAppResource, clientset)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("Running myAppResource update callback")
			//fmt.Printf("Old object: %+v\n", oldObj)
			//fmt.Printf("New object: %+v\n", newObj)
			oldMyAppResource, err := toMyAppResource(oldObj)
			if err != nil {
				return
			}
			newMyAppResource, err := toMyAppResource(newObj)
			if err != nil {
				return
			}
			controller.HandleMyAppResourceUpdate(oldMyAppResource, newMyAppResource, clientset)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("running myAppResource add callback")
			myAppResource, err := toMyAppResource(obj)
			if err != nil {
				return
			}
			controller.HandleMyAppResourceDelete(myAppResource, clientset)
		},
	})
	fmt.Println("MyAppResource controller started")

	// Wait forever or until an error occurs
	select {}
}

func toMyAppResource(obj interface{}) (*mar.MyAppResource, error) {
	unstructuredObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("invalid object type received")
	}

	// Unmarshal the unstructured object to your custom MyAppResource type
	myAppResource := &mar.MyAppResource{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.Object, myAppResource); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the object: %v", err)
	}

	return myAppResource, nil
}
