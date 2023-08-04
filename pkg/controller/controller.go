package controller

import (
	"context"
	"fmt"
	mar "go-angi/api/myapigroup/v1alpha1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"reflect"
	"time"
)

func InitKubernetesClient() (*kubernetes.Clientset, dynamic.Interface) {
	// Load the in-cluster Kubernetes configuration.
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// Create the clientset.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create the dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset, dynamicClient
}

func HandleMyAppResourceAdd(myAppResource *mar.MyAppResource, clientset *kubernetes.Clientset) {
	fmt.Printf("Handling new MyAppResource %s\n", myAppResource.Name)

	var redisPod *v1.Pod = nil
	var redisService *v1.Service = nil
	if myAppResource.Spec.Redis.Enabled {
		fmt.Printf("Handling new MyAppResource %s\n", myAppResource.Name)

		redisPod = &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      myAppResource.Name + "-redis",
				Namespace: myAppResource.Namespace,
				Labels:    map[string]string{"app": "myapp-redis"},
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					v1.Container{
						Name:  "redis",
						Image: "redis:latest",
					},
				},
			},
		}

		_, err := clientset.CoreV1().Pods(myAppResource.Namespace).Create(context.TODO(), redisPod, metav1.CreateOptions{})
		if err != nil {
			if apierrors.IsAlreadyExists(err) {
				fmt.Printf("Redis pod %s already exists, proceeding\n", redisPod.Name)
			} else {
				fmt.Printf("Error creating Redis pod for MyAppResource %s: %v\n", myAppResource.Name, err)
				return
			}
		} else {
			fmt.Printf("Created Redis pod for MyAppResource %s\n", myAppResource.Name)
		}

		redisService = &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      myAppResource.Name + "-redis",
				Namespace: myAppResource.Namespace,
			},
			Spec: v1.ServiceSpec{
				Selector: map[string]string{
					"app": "myapp-redis",
				},
				Ports: []v1.ServicePort{
					v1.ServicePort{
						Port:       6379,
						TargetPort: intstr.FromInt(6379),
						Protocol:   v1.ProtocolTCP,
					},
				},
			},
		}

		_, err = clientset.CoreV1().Services(myAppResource.Namespace).Create(context.TODO(), redisService, metav1.CreateOptions{})
		if err != nil {
			if apierrors.IsAlreadyExists(err) {
				fmt.Printf("Service %s already exists, proceeding\n", redisService.Name)
			} else {
				fmt.Printf("Error creating redis service for MyAppResource %s: %v\n", myAppResource.Name, err)
			}
		} else {
			fmt.Printf("Created redis service for MyAppResource %s\n", myAppResource.Name)
		}
	}

	for i := 0; i < int(myAppResource.Spec.ReplicaCount); i++ {
		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      myAppResource.Name + fmt.Sprintf("-pod-%d", i),
				Namespace: myAppResource.Namespace,
				Labels:    map[string]string{"app": "myapp"},
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					v1.Container{
						Name:  "myapp-container",
						Image: myAppResource.Spec.Image.Repository + ":" + myAppResource.Spec.Image.Tag,
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceCPU: resource.MustParse(myAppResource.Spec.Resources.CpuRequest),
							},
							Limits: v1.ResourceList{
								v1.ResourceMemory: resource.MustParse(myAppResource.Spec.Resources.MemoryLimit),
							},
						},
						Env: []v1.EnvVar{
							v1.EnvVar{
								Name:  "PODINFO_UI_COLOR",
								Value: myAppResource.Spec.Ui.Color,
							},
							v1.EnvVar{
								Name:  "PODINFO_UI_MESSAGE",
								Value: myAppResource.Spec.Ui.Message,
							},
						},
					},
				},
			},
		}

		if redisPod != nil && redisService != nil {
			fmt.Printf("Setting PODINFO_CACHE_SERVER on pod %s\n", pod.Name)
			pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env,
				v1.EnvVar{
					Name:  "PODINFO_CACHE_SERVER",
					Value: "redis://" + redisService.Name + ":6379",
				})
		}

		_, err := clientset.CoreV1().Pods(myAppResource.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
		if err != nil {
			if apierrors.IsAlreadyExists(err) {
				fmt.Printf("Pod %s already exists, proceeding\n", pod.Name)
			} else {
				fmt.Printf("Error creating pod for MyAppResource %s: %v\n", myAppResource.Name, err)
			}
		} else {
			fmt.Printf("Created pod for MyAppResource %s\n", myAppResource.Name)
		}
	}

	// Create a new service for podinfo
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      myAppResource.Name + "-service",
			Namespace: myAppResource.Namespace,
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"app": "myapp",
			},
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Name:       "http",
					Protocol:   v1.ProtocolTCP,
					Port:       9898,
					TargetPort: intstr.FromInt(9898),
				},
			},
			Type: v1.ServiceTypeClusterIP, // Explicitly set the ServiceType as ClusterIP
		},
	}

	_, err := clientset.CoreV1().Services(myAppResource.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			fmt.Printf("Service %s already exists, proceeding\n", service.Name)
		} else {
			fmt.Printf("Error creating service for MyAppResource %s: %v\n", myAppResource.Name, err)
		}
	} else {
		fmt.Printf("Created service for MyAppResource %s\n", myAppResource.Name)
	}
}

func HandleMyAppResourceDelete(myAppResource *mar.MyAppResource, clientset *kubernetes.Clientset) {
	fmt.Printf("Handling delete of MyAppResource %s\n", myAppResource.Name)

	// Delete Podinfo pods
	for i := 0; i < int(myAppResource.Spec.ReplicaCount); i++ {
		podName := myAppResource.Name + fmt.Sprintf("-pod-%d", i)
		err := clientset.CoreV1().Pods(myAppResource.Namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				fmt.Printf("Pod %s not found, probably already deleted\n", podName)
			} else {
				fmt.Printf("Error deleting Pod %s: %v\n", podName, err)
			}
		} else {
			fmt.Printf("Deleted Pod %s\n", podName)
		}
	}

	// Delete Podinfo service
	serviceName := myAppResource.Name + "-service"
	err := clientset.CoreV1().Services(myAppResource.Namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			fmt.Printf("Service %s not found, probably already deleted\n", serviceName)
		} else {
			fmt.Printf("Error deleting Service %s: %v\n", serviceName, err)
		}
	} else {
		fmt.Printf("Deleted Service %s\n", serviceName)
	}

	if myAppResource.Spec.Redis.Enabled {
		// Delete Redis pod
		redisPodName := myAppResource.Name + "-redis"
		err = clientset.CoreV1().Pods(myAppResource.Namespace).Delete(context.TODO(), redisPodName, metav1.DeleteOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				fmt.Printf("Redis Pod %s not found, probably already deleted\n", redisPodName)
			} else {
				fmt.Printf("Error deleting Redis Pod %s: %v\n", redisPodName, err)
			}
		} else {
			fmt.Printf("Deleted Redis Pod %s\n", redisPodName)
		}

		// Delete Redis service
		redisServiceName := myAppResource.Name + "-redis"
		err = clientset.CoreV1().Services(myAppResource.Namespace).Delete(context.TODO(), redisServiceName, metav1.DeleteOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				fmt.Printf("Redis Service %s not found, probably already deleted\n", redisServiceName)
			} else {
				fmt.Printf("Error deleting Redis Service %s: %v\n", redisServiceName, err)
			}
		} else {
			fmt.Printf("Deleted Redis Service %s\n", redisServiceName)
		}
	}
}

func HandleMyAppResourceUpdate(oldMyAppResource, newMyAppResource *mar.MyAppResource, clientset *kubernetes.Clientset) {
	fmt.Printf("Handling update of MyAppResource %s\nold color = %s, new color = %s\nold message = %s, new message = %s\n",
		newMyAppResource.Name,
		oldMyAppResource.Spec.Ui.Color,
		newMyAppResource.Spec.Ui.Color,
		oldMyAppResource.Spec.Ui.Message,
		newMyAppResource.Spec.Ui.Message)

	// Check if the spec has changed
	if reflect.DeepEqual(oldMyAppResource.Spec, newMyAppResource.Spec) {
		fmt.Printf("No changes to MyAppResource %s, skipping update\n", newMyAppResource.Name)
		return
	}

	// Delete the old resources
	HandleMyAppResourceDelete(oldMyAppResource, clientset)

	// Wait for a while to allow the resources to be deleted
	fmt.Printf("Waiting for deletion to complete before re-adding\n")
	time.Sleep(time.Duration(30) * time.Second)

	// Create the new resources
	HandleMyAppResourceAdd(newMyAppResource, clientset)
}
