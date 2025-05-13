// Api versions allow the api contract for a resource to be changed while keeping
// backward compatibility by support multiple concurrent versions
// of the same resource

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=dev.khulnasoft.com/test/apis/test
// +k8s:defaulter-gen=TypeMeta
// +groupName=test.loft.sh
package v1 // import "dev.khulnasoft.com/test/apis/test/v1"
