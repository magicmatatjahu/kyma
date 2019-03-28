package client

import (
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewRestClientConfig(kubeconfigPath string) (*rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func NewRestClientConfigFromEnv() (*rest.Config, error) {
	kubeConfigPath := os.Getenv("KUBECONFIG")
	return NewRestClientConfig(kubeConfigPath)
}
