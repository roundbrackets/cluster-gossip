package main

import (
	"context"
	"fmt"
	"os"
	//"strings"
	//	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	c "github.com/roundbrackets/cluster-gossip/client"
)

func main() {

	client, err := c.MakeClient(c.ClientOpts{ClientType: c.CommandLine})
	if err != nil {
		fmt.Printf("Error encountered: %v\n", err)
		os.Exit(1)		
	}

	ns := namespaces (client) 
	for n := range ns {
		n.ToString()
	}
}

type Event struct {
	namespace string
	podCnt int
	eventCnt int
	events []string
}

func (e *Event) ToString () {
	fmt.Printf("%-30s %3d %3d\n", e.namespace, e.podCnt, e.eventCnt);

	for _, e := range e.events {
		fmt.Printf("  %s\n", e);
	}
}

func namespaces (client *c.Client) <- chan *Event {
	out := make(chan *Event)

	ns, err := client.Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error encountered: %v\n", err)
		return out 
	}

	go func() {
		for _, n := range ns.Items {
			event := &Event{namespace: n.Name} 
			event.podCnt = pods(client, n.Name)
			event.events, event.eventCnt = events(client, n.Name)
			out <- event 
		}
		close(out)
	}()

	return out
}

func pods (client *c.Client, namespace string) int {
	pods, err := client.Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error encountered: %v\n", err)
		return -1
	}

	return len(pods.Items)
}

func events (client *c.Client, namespace string) ([]string, int) {
	events, err := client.Events(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("  Error encountered: %v\n", err)
		return nil, -1
	}

	e := []string{}
	for _, event := range events.Items {
		e = append(e, event.Message)
	}

	return e, len(events.Items)

	/*
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
	*/
}
