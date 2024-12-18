package ingredient

import (
	"github.com/yuin/goldmark/renderer/html"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// HTMLRenderer struct is a renderer.NodeRenderer implementation for the extension.
type HTMLRenderer struct{}

// NewHTMLRenderer builds a new HTMLRenderer with given options and returns it.
func NewHTMLRenderer() renderer.NodeRenderer {
	return &HTMLRenderer{}
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs.
func (r *HTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindQuantity, r.renderQuantity)
	reg.Register(KindIngredient, r.renderIngredient)
}

var IngredientAttributeFilter = html.GlobalAttributeFilter.Extend([]byte("itemprop"), []byte("content"))

func (r *HTMLRenderer) renderIngredient(w util.BufWriter, source []byte,
	node ast.Node, entering bool) (ast.WalkStatus, error) {
	var err error
	if entering {
		_, err = w.WriteString(`<span`)
		if err == nil {
			html.RenderAttributes(w, node, IngredientAttributeFilter)
			_, err = w.WriteString(`>`)
		}
		return ast.WalkContinue, err
	}

	_, err = w.WriteString("</span>")
	return ast.WalkContinue, err
}

func (r *HTMLRenderer) renderQuantity(w util.BufWriter, source []byte,
	node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		return ast.WalkContinue, nil
	}

	_, err := w.WriteString(`<span`)
	if err == nil {
		html.RenderAttributes(w, node, IngredientAttributeFilter)
		_, err = w.WriteString(`>`)
		if err == nil {
			_, err = w.WriteString(node.(*Quantity).Value)
			if err == nil {
				_, err = w.WriteString(`</span>`)
			}
		}
	}
	return ast.WalkContinue, err
}
