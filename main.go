package main

import (
	"context"
	"fmt"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func WatchEvents() {

	// Location of kubeconfig file
	kubeconfig := os.Getenv("HOME") + "/.kube/config"

	// Create a Config (k8s.io/client-go/rest)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// Create an API Clientset (k8s.io/client-go/kubernetes)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create a CoreV1Client (k8s.io/client-go/kubernetes/typed/core/v1)
	coreV1Client := clientset.CoreV1()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	evts, err := coreV1Client.Events("").Watch(ctx, metaV1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// List each event (k8s.io/api/core/v1)
	for {
		result := <- evts.ResultChan()
		fmt.Printf("result: %v\n",result)
	}

}



func main() {

	// Location of kubeconfig file
	kubeconfig := os.Getenv("HOME") + "/.kube/config"

	// Create a Config (k8s.io/client-go/rest)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// Create an API Clientset (k8s.io/client-go/kubernetes)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create a CoreV1Client (k8s.io/client-go/kubernetes/typed/core/v1)
	coreV1Client := clientset.CoreV1()
	// Create an AppsV1Client (k8s.io/client-go/kubernetes/typed/apps/v1)
	appsV1Client := clientset.AppsV1()

	//-------------------------------------------------------------------------//
	// List pods (all namespaces)
	//-------------------------------------------------------------------------//

	// Get a *PodList (k8s.io/api/core/v1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pods, err := coreV1Client.Pods("").List(ctx, metaV1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// List each Pod (k8s.io/api/core/v1)
	for i, pod := range pods.Items {
		fmt.Printf("Pod %d: %s\n", i+1, pod.ObjectMeta.Name)
	}

	//-------------------------------------------------------------------------//
	// List events (all namespaces)
	// k8s removes events after 1 hr.
	//-------------------------------------------------------------------------//

	// TODO: Clean this up
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	evts, err := coreV1Client.Events("").List(ctx2, metaV1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// List each event (k8s.io/api/core/v1)
	for i, v := range evts.Items {
		fmt.Printf("Event %d: %s ", i+1, v.ObjectMeta.Name)
		fmt.Printf("TimeStamp %v: Msg: %s\n", v.ObjectMeta.CreationTimestamp, v.Message)
		fmt.Printf("     FirstTimestamp: %v   Count: %v \n", v.FirstTimestamp, v.Count)
	}

	//-------------------------------------------------------------------------//
	// List nodes
	//-------------------------------------------------------------------------//

	// Get a *NodeList (k8s.io/api/core/v1)
	nodes, err := coreV1Client.Nodes().List(ctx, metaV1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// For each Node (k8s.io/api/core/v1)
	for i, node := range nodes.Items {
		fmt.Printf("Node %d: %s\n", i+1, node.ObjectMeta.Name)
	}

	//-------------------------------------------------------------------------//
	// List deployments (all namespaces)
	//-------------------------------------------------------------------------//

	// Get a *DeploymentList (k8s.io/api/apps/v1)
	deployments, err := appsV1Client.Deployments("").List(ctx, metaV1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// For each Deployment (k8s.io/api/apps/v1)
	for i, deployment := range deployments.Items {
		fmt.Printf("Deployment %d: %s\n", i+1, deployment.ObjectMeta.Name)
	}

	WatchEvents()
}
