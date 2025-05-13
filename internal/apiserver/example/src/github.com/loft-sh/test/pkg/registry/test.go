package registry

import (
	"context"

	"dev.khulnasoft.com/test/apis/test"
	testv1 "dev.khulnasoft.com/test/apis/test/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
)

func init() {
	test.NewClusterRoleRESTFunc = NewREST
}

// NewREST implements the storage interface
func NewREST() rest.Storage {
	return &REST{
		TableConvertor: rest.NewDefaultTableConvertor(testv1.Resource("clusterroles")),
	}
}

type REST struct {
	rest.TableConvertor
}

func (r *REST) New() runtime.Object {
	return &testv1.ClusterRole{}
}

func (r *REST) Destroy() {}

func (r *REST) NamespaceScoped() bool {
	return false
}

func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	// return object
	return &testv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}, nil
}
