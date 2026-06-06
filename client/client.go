package client 

import (
	"os"
	"path/filepath"
	"errors"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/rest"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Type int

const (
	CommandLine Type = iota
	InCluster               
)

type ClientOpts struct {
	ClientType Type
	KubeconfigPath string
}

func (c *ClientOpts) Kubeconfig () (string, error) {
	if (c.KubeconfigPath == "") {
		if home := os.Getenv("HOME"); home != "" {
			c.KubeconfigPath = filepath.Join(home, ".kube", "config")
		} else {
			return "", errors.New("No file found")
		}
	}
	// if path is set verify that it exists
	return c.KubeconfigPath, nil 
}

type Client struct {
	corev1.CoreV1Interface
}

func CommandLineClient (clientOpts ClientOpts) (clientset *kubernetes.Clientset, err error) {
	var config *rest.Config
	var kubeconfig string

	kubeconfig, err = clientOpts.Kubeconfig() 
	if err != nil {
		return
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return
	}

	config.QPS = 50
	config.Burst = 100

	clientset, err = kubernetes.NewForConfig(config)

	return
}

func InClusterClient (clientOpts ClientOpts) (clientset *kubernetes.Clientset, err error) {
		var config *rest.Config
		config, err = rest.InClusterConfig()
		if err != nil {
			return
		}

	clientset, err = kubernetes.NewForConfig(config)

	return
}

func MakeClient (clientOpts ClientOpts) (client *Client, err error) {
	var clientset *kubernetes.Clientset
	if clientOpts.ClientType == CommandLine { 
		clientset, err = CommandLineClient (clientOpts)
		if err != nil {
			return
		}
	} else if clientOpts.ClientType == InCluster { 
		clientset, err = InClusterClient (clientOpts)
		if err != nil {
			return
		}
	}

	client = &Client{CoreV1Interface: clientset.CoreV1()}

	return
}
