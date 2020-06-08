/*
Copyright 2019 The OpenEBS Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package k8s

import (
	"github.com/openebs/node-disk-manager/integration_tests/utils"
	"github.com/openebs/node-disk-manager/pkg/apis"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"time"

	//"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// Namespace is the default namespace
	Namespace = "default"
	// WaitDuration is the default wait duration
	WaitDuration time.Duration = 5 * time.Second
	// Running is the active/running status of pod
	Running = "Running"
)

// K8sClient is the client used for etcd operations
type K8sClient struct {
	config        *rest.Config
	ClientSet     *kubernetes.Clientset
	APIextClient  *apiextensionsclient.Clientset
	RunTimeClient client.Client
}

// GetClientSet generates the client-set from the config file on the host
// Three clients are generated by the function:
// 1. Client-set from client-go which is used for operations on pods/nodes and
//    other first-class k8s objects. While using runtime client for this operations,
//    the podList etc retrieved were not up-to-date
// 2. Client from apiextensions which is used for CRUD operations on CRDs
// 3. Runtime client from the controller-runtime which will be used for
//    CRUD operations related to custom resources
func GetClientSet() (K8sClient, error) {
	clientSet := K8sClient{}
	kubeConfigPath, err := utils.GetConfigPath()
	if err != nil {
		return clientSet, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return clientSet, err
	}
	clientSet.config = config
	// client-go clientSet
	clientSet.ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		return clientSet, err
	}

	// client for creating CRDs
	clientSet.APIextClient, err = apiextensionsclient.NewForConfig(config)
	if err != nil {
		return clientSet, err
	}

	// controller-runtime client
	mgr, err := manager.New(config, manager.Options{Namespace: Namespace, MetricsBindAddress: "0"})
	if err != nil {
		return clientSet, err
	}

	// add to scheme
	scheme := mgr.GetScheme()
	if err = apis.AddToScheme(scheme); err != nil {
		return clientSet, err
	}

	clientSet.RunTimeClient, err = client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return clientSet, err
	}
	return clientSet, nil
}

func (k *K8sClient) RegenerateClient() error {
	// controller-runtime client
	mgr, err := manager.New(k.config, manager.Options{Namespace: Namespace, MetricsBindAddress: "0"})
	if err != nil {
		return err
	}

	// add to scheme
	scheme := mgr.GetScheme()
	if err = apis.AddToScheme(scheme); err != nil {
		return err
	}

	k.RunTimeClient, err = client.New(k.config, client.Options{Scheme: scheme})
	if err != nil {
		return err
	}
	return nil
}
