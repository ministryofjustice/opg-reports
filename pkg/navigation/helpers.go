package navigation

import (
	"log/slog"
	"net/http"
	"strings"
)

// Flat traverses the `tree` recursively calls itself on Children
// appending to the `flat` map passed
func Flat(tree []*Navigation, flat map[string]*Navigation) {
	slog.Debug("[navigation.Flat] traversing tree")
	for _, node := range tree {
		if node != nil {
			var key = node.Uri
			flat[key] = node
			// recurse if this node has children
			if len(node.Children) > 0 {
				slog.Debug("[navigation.Flat] recurse")
				Flat(node.Children, flat)
			}
		}
	}
	return
}

// ActivateTree traverses the tree and marks each item if they
// are within the requested uri or are an exact match
// Used so we can mark the navigation items that relate to the
// active page stack
func ActivateTree(tree []*Navigation, request *http.Request) {
	var url string = request.URL.Path
	// remove leading and trailing url, re-add leading
	url = strings.TrimSuffix(url, "/")
	url = "/" + strings.TrimPrefix(url, "/")

	for _, node := range tree {
		// clean up the node uri as well to remove
		// extra slashes
		var nodeUrl = node.Uri
		nodeUrl = strings.TrimSuffix(nodeUrl, "/")
		nodeUrl = "/" + strings.TrimPrefix(nodeUrl, "/")

		if strings.HasPrefix(url, nodeUrl) {
			node.Display.InUri = true
		}
		if url == nodeUrl {
			node.Display.IsActive = true
		}
		// recurse
		if len(node.Children) > 0 {
			ActivateTree(node.Children, request)
		}
	}
	return
}
