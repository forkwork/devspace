package generate

import (
	"fmt"
	"path/filepath"
	"strings"

	"k8s.io/gengo/v2/types"
)

// IsAPIResource returns true if t has a +resource comment tag
func IsAPIResource(t *types.Type) bool {
	for _, c := range t.CommentLines {
		if strings.Contains(c, "+resource") || strings.Contains(c, "+kubebuilder:resource") {
			return true
		}
	}
	return false
}

// IsNonNamespaced returns true if t has a +nonNamespaced comment tag
func IsNonNamespaced(t *types.Type) bool {
	if !IsAPIResource(t) {
		return false
	}

	for _, c := range t.CommentLines {
		if strings.Contains(c, "+genclient:nonNamespaced") {
			return true
		}
	}

	for _, c := range t.SecondClosestCommentLines {
		if strings.Contains(c, "+genclient:nonNamespaced") {
			return true
		}
	}

	return false
}

// IsAPISubresource returns true if t has a +subresource-request comment tag
func IsAPISubresource(t *types.Type) bool {
	for _, c := range t.CommentLines {
		if strings.Contains(c, "+subresource-request") {
			return true
		}
	}
	return false
}

// HasSubresource returns true if t is an APIResource with one or more Subresources
func HasSubresource(t *types.Type) bool {
	if !IsAPIResource(t) {
		return false
	}
	for _, c := range t.CommentLines {
		if strings.Contains(c, "+subresource") {
			return true
		}
	}
	return false
}

func IsUnversioned(t *types.Type, group string) bool {
	return IsApisDir(filepath.Base(filepath.Dir(t.Name.Package))) && GetGroup(t) == group
}

func IsVersioned(t *types.Type, group string) bool {
	return GetGroup(t) == group
}

func GetVersion(t *types.Type, group string) string {
	if !IsVersioned(t, group) {
		panic(fmt.Errorf("Cannot get version for unversioned type %v", t.Name))
	}
	return filepath.Base(t.Name.Package)
}

func GetGroup(t *types.Type) string {
	return filepath.Base(GetGroupPackage(t))
}

func GetGroupPackage(t *types.Type) string {
	if IsApisDir(filepath.Base(filepath.Dir(t.Name.Package))) {
		return t.Name.Package
	}
	return filepath.Dir(t.Name.Package)
}

func GetKind(t *types.Type, group string) string {
	if !IsVersioned(t, group) && !IsUnversioned(t, group) {
		panic(fmt.Errorf("Cannot get kind for type not in group %v", t.Name))
	}
	return t.Name.Name
}

// IsApisDir returns true if a directory path is a Kubernetes api directory
func IsApisDir(dir string) bool {
	return dir == "apis" || dir == "api"
}

// Comments is a structure for using comment tags on go structs and fields
type Comments []string

// GetTags returns the value for the first comment with a prefix matching "+name="
// e.g. "+name=foo\n+name=bar" would return "foo"
func (c Comments) GetTag(name, sep string) string {
	for _, c := range c {
		prefix := fmt.Sprintf("+%s%s", name, sep)
		if strings.HasPrefix(c, prefix) {
			return strings.Replace(c, prefix, "", 1)
		}
	}
	return ""
}

func (c Comments) HasTag(name string) bool {
	for _, c := range c {
		prefix := fmt.Sprintf("+%s", name)
		if strings.HasPrefix(c, prefix) {
			return true
		}
	}
	return false
}

// GetTags returns the value for all comments with a prefix and separator.  E.g. for "name" and "="
// "+name=foo\n+name=bar" would return []string{"foo", "bar"}
func (c Comments) GetTags(name, sep string) []string {
	tags := []string{}
	for _, c := range c {
		prefix := fmt.Sprintf("+%s%s", name, sep)
		if strings.HasPrefix(c, prefix) {
			tags = append(tags, strings.Replace(c, prefix, "", 1))
		}
	}
	return tags
}
