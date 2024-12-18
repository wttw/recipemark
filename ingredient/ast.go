package ingredient

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
)

// An Item struct represents an ingredient-style list item of Markdown text.
type Item struct {
	ast.BaseBlock

	// Offset is an offset position of this item.
	Offset int
}

// Dump implements Node.Dump.
func (n *Item) Dump(source []byte, level int) {
	m := map[string]string{
		"Offset": fmt.Sprintf("%d", n.Offset),
	}
	ast.DumpHelper(n, source, level, m, nil)
}

var KindIngredient = ast.NewNodeKind("ingredient")

// Kind implements Node.Kind.
func (n *Item) Kind() ast.NodeKind {
	return KindIngredient
}

// New returns a new Ingredient node.
func New() *Item {
	return &Item{
		BaseBlock: ast.BaseBlock{},
	}
}

type Quantity struct {
	ast.BaseBlock

	Value string
}

// Dump implements Node.Dump.
func (n *Quantity) Dump(source []byte, level int) {
	m := map[string]string{
		"Value": n.Value,
	}
	ast.DumpHelper(n, source, level, m, nil)
}

var KindQuantity = ast.NewNodeKind("quantity")

// Kind implements Node.Kind.
func (n *Quantity) Kind() ast.NodeKind {
	return KindQuantity
}

// NewQuantity returns a new Ingredient node.
func NewQuantity(val string) *Quantity {
	return &Quantity{
		BaseBlock: ast.BaseBlock{},
		Value:     val,
	}
}

// mergeAttribute appends additional values to an attribute
// formed of space separated values, e.g. class
func mergeAttribute(n ast.Node, attrName string, values string) {
	existing, ok := n.AttributeString(attrName)
	if !ok {
		n.SetAttributeString(attrName, []byte(values))
	}
	all := map[string]struct{}{}
	existingBytes, ok := existing.([]byte)
	existingString := string(existingBytes)
	var existingClasses []string
	if ok {
		for _, c := range strings.Fields(existingString) {
			all[c] = struct{}{}
			existingClasses = append(existingClasses, c)
		}
	}

	for _, newClass := range strings.Fields(values) {
		_, ok = all[newClass]
		if !ok {
			existingClasses = append(existingClasses, newClass)
		}
	}
	if existingClasses != nil {
		n.SetAttributeString(attrName, []byte(strings.Join(existingClasses, " ")))
	}
}

// addClasses appends HTML classes to a node
func addClasses(n ast.Node, classes string) {
	mergeAttribute(n, "class", classes)
}
