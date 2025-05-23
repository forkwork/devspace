// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "dev.khulnasoft.com/api/v4/pkg/apis/management/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// LoftUpgradeLister helps list LoftUpgrades.
// All objects returned here must be treated as read-only.
type LoftUpgradeLister interface {
	// List lists all LoftUpgrades in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.LoftUpgrade, err error)
	// Get retrieves the LoftUpgrade from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.LoftUpgrade, error)
	LoftUpgradeListerExpansion
}

// loftUpgradeLister implements the LoftUpgradeLister interface.
type loftUpgradeLister struct {
	listers.ResourceIndexer[*v1.LoftUpgrade]
}

// NewLoftUpgradeLister returns a new LoftUpgradeLister.
func NewLoftUpgradeLister(indexer cache.Indexer) LoftUpgradeLister {
	return &loftUpgradeLister{listers.New[*v1.LoftUpgrade](indexer, v1.Resource("loftupgrade"))}
}
