package main

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/roundbrackets/cluster-gossip/client"
)

func main() {

	client, err := MakeClient(ClientOpts{ClientType: CommandLine})
	if err != nil {
		fmt.Printf("Error encountered: %v\n", err)
		return
	}

	namespaces, err := client.Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error encountered: %v\n", err)
		return
	}
	for _, namespace := range namespaces.Items {
		fmt.Printf("%-40s\n", namespace.Name);
		fmt.Printf("%-40s\n", strings.Repeat("=", len(namespace.Name)));

		fmt.Println("\nEvents\n------");
		events, err := client.Events(namespace.Name).List(context.TODO(), metav1.ListOptions{})
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
		pods, err := client.Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{})
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
