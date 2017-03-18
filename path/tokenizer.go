package path

import (
	"github.com/go-scim/scimify/adt"
	"strings"
	"errors"
	"fmt"
	"github.com/go-scim/scimify/filter"
)

func Tokenize(path string) (adt.Queue, error) {
	queue := adt.NewQueueWithoutLimit()

	for _, component := range strings.Split(path, ".") {
		if len(component) == 0 {
			return nil, errors.New("empty component in path")
		}

		switch {
		case isSimpleComponent(component):
			queue.Offer(adt.NewNode(filter.Token{
				Value:component,
				Type:filter.Path,
			}))

		case componentContainsFilter(component):
			simple, filterStr := splitFilterWithSimple(component)
			queue.Offer(adt.NewNode(filter.Token{
				Value:simple,
				Type:filter.Path,
			}))
			if tokens, err := filter.Tokenize(filterStr); err != nil {
				return nil, err
			} else if node, err := filter.Parse(tokens); err != nil {
				return nil, err
			} else {
				queue.Offer(node)
			}

		default:
			return nil, fmt.Errorf("invalid component in path: %s", component)
		}
	}

	return queue, nil
}

func splitFilterWithSimple(component string) (string, string) {
	return component[:strings.Index(component, "[")], component[strings.Index(component, "[")+1:len(component)-1]
}

func componentContainsFilter(component string) bool {
	return strings.Count(component, "[") == 1 &&
		strings.Count(component, "]") == 1 &&
		strings.Index(component, "[") < strings.Index(component, "]") - 1 &&
		strings.Index(component, "[") > 0 &&
		strings.Index(component, "]") == len(component) - 1
}

func isSimpleComponent(component string) bool {
	return strings.Count(component, "[") == 0 &&
		strings.Count(component, "]") == 0
}