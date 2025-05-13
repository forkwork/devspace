/*
Copyright 2020 DevSpace Technologies Inc.

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

package main

import (
	apiserver "dev.khulnasoft.com/apiserver/pkg/server"
	"dev.khulnasoft.com/test/apis"
	"dev.khulnasoft.com/test/apis/test"
	testv1 "dev.khulnasoft.com/test/apis/test/v1"
	"dev.khulnasoft.com/test/pkg/openapi"
	_ "dev.khulnasoft.com/test/pkg/registry"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// +kubebuilder:scaffold:imports

	// Make sure dep tools picks up these dependencies
	//_ "github.com/go-openapi/loads"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Enable cloud provider auth
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	// API extensions are not in the above scheme set,
	// and must thus be added separately.
	_ = apiregistrationv1.AddToScheme(scheme)

	_ = test.AddToScheme(scheme)
	_ = testv1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	// Start the api server
	err := apiserver.StartAPIServer(&apiserver.StartOptions{
		Apis:                  apis.GetAllApiBuilders(),
		GetOpenAPIDefinitions: openapi.GetOpenAPIDefinitions,
	})
	if err != nil {
		panic(err)
	}
}
