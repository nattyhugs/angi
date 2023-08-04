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
			unstructuredObj, ok := obj.(*unstructured.Unstructured)
			if !ok {
				fmt.Println("Invalid object type received in AddFunc")
				return
			}

			// Unmarshal the unstructured object to your custom MyAppResource type
			myAppResource := &mar.MyAppResource{}
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.Object, myAppResource); err != nil {
				fmt.Println("Failed to unmarshal the object:", err)
				return
			}
			controller.HandleMyAppResourceAdd(myAppResource, clientset)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("running myAppResource update callback")
			unstructuredNewObj, ok := newObj.(*unstructured.Unstructured)
			if !ok {
				fmt.Println("Invalid object type received in UpdateFunc")
				return
			}

			// Unmarshal the unstructured object to your custom MyAppResource type
			newMyAppResource := &mar.MyAppResource{}
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredNewObj.Object, newMyAppResource); err != nil {
				fmt.Println("Failed to unmarshal the object:", err)
				return
			}

			unstructuredOldObj, ok := newObj.(*unstructured.Unstructured)
			if !ok {
				fmt.Println("Invalid object type received in UpdateFunc")
				return
			}

			// Unmarshal the unstructured object to your custom MyAppResource type
			oldMyAppResource := &mar.MyAppResource{}
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredOldObj.Object, oldMyAppResource); err != nil {
				fmt.Println("Failed to unmarshal the object:", err)
				return
			}

			controller.HandleMyAppResourceUpdate(oldMyAppResource, newMyAppResource, clientset)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("running myAppResource add callback")
			unstructuredObj, ok := obj.(*unstructured.Unstructured)
			if !ok {
				fmt.Println("Invalid object type received in AddFunc")
				return
			}

			// Unmarshal the unstructured object to your custom MyAppResource type
			myAppResource := &mar.MyAppResource{}
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.Object, myAppResource); err != nil {
				fmt.Println("Failed to unmarshal the object:", err)
				return
			}
			controller.HandleMyAppResourceDelete(myAppResource, clientset)
		},
	})
	fmt.Println("MyAppResource controller started")

	// Wait forever or until an error occurs
	select {}
}
