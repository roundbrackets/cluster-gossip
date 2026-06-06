package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)


// https://github.com/kubernetes/client-go
// Server Version: v1.31.4+k3s1
// v0.x.y tags f
//go: downloading k8s.io/client-go v0.36.1
//go: added k8s.io/client-go v0.36.1
//

// https://pkg.go.dev/k8s.io/client-go#section-readme
// https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/main.go

func main() {
	// The default location for the kubeconfig file is in the user's home directory.
	var kubeconfig string
	if home := os.Getenv("HOME"); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// Create the client configuration from the kubeconfig file.
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("Error encountered: %v\n", err)
		return
	}

	// Configure client-side rate limiting.
	config.QPS = 50
	config.Burst = 100

	// A clientset contains clients for all the API groups and versions supported
	// by the cluster.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error encountered: %v\n", err)
		return
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error encountered: %v\n", err)
		return
	}
	for _, namespace := range namespaces.Items {
		fmt.Printf("%-40s\n", namespace.Name);
		fmt.Printf("%-40s\n", strings.Repeat("=", len(namespace.Name)));

		fmt.Println("\nEvents\n------");
		events, err := clientset.CoreV1().Events(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("  Error encountered: %v\n", err)
			return
		}
		for _, event := range events.Items {
			fmt.Printf("  %s\n", event.Message);
			fmt.Printf("    %-30s %-30s %-30s\n", event.Reason, event.Type, event.Action);
			fmt.Printf("    %-30s %-30s\n", event.FirstTimestamp, event.LastTimestamp);
			// EventTime, Series
		}

		fmt.Printf("  There are %d pods in the %s namespace\n", len(events.Items), namespace.Name)

		fmt.Println("Pods\n----");
		pods, err := clientset.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error encountered: %v\n", err)
			return
		}
		for _, pod := range pods.Items {
			fmt.Printf("  %s\n", pod.Name);
		}

		fmt.Printf("  There are %d pods in the %s namespace\n\n", len(pods.Items), namespace.Name)
	}

	fmt.Printf("There are %d namespaces\n", len(namespaces.Items))

}
