package ingredient

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func init() {
	// Monkey with the Goldmark permitted attributes
	html.ListItemAttributeFilter.Add([]byte("itemprop"))
}

type Extender struct{}

func (e Extender) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(
		parser.WithInlineParsers(util.Prioritized(
			&Parser{}, 120),
		),
	)
	md.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewHTMLRenderer(), 500),
		),
	)
	md.Parser().AddOptions(
		parser.WithASTTransformers(util.Prioritized(&Transformer{}, 200)),
	)
}

type Transformer struct {
}

//https://www.delish.com/cooking/recipe-ideas/a22024047/best-szechuan-beef-recipe/
//
//Use the first image as the image for indexes etc.

// Transform walks the doc converting text before a Quantity tag to an Ingredient tag
// and adding itemprop tags
func (t *Transformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	firstBlockquote := true
	seenIngredients := false
	_ = ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Kind() == KindQuantity {
			seenIngredients = true
			addClasses(n, "quantity")
			parent := n.Parent()
			ingredient := New()
			peer := n.PreviousSibling()
			for peer != nil {
				parent.RemoveChild(parent, peer)
				ingredient.InsertBefore(ingredient, ingredient.FirstChild(), peer)
				peer = n.PreviousSibling()
			}
			trailing := ingredient.LastChild()
			var extra *ast.Text
			if trailing != nil && trailing.Kind() == ast.KindText {
				txt := trailing.(*ast.Text)
				seg := txt.Segment
				space := util.TrimRightSpaceLength(reader.Value(seg))
				if space > 0 {
					txt.Segment = seg.WithStop(seg.Stop - space)
					extra = ast.NewTextSegment(seg.WithStart(seg.Stop - space))
				}
			}
			parent.InsertBefore(parent, n, ingredient)
			addClasses(ingredient, "ingredient")
			if extra != nil {
				parent.InsertAfter(parent, ingredient, extra)
			}
			// find the parent <li>
			p := n
			for p != nil {
				if p.Kind() == ast.KindListItem {
					mergeAttribute(p, "itemprop", "recipeIngredient")
					break
				}
				p = p.Parent()
			}
		}
		if entering && n.Kind() == ast.KindParagraph && seenIngredients {
			// Add itemprop unless we have a blockquote parent
			p := n
			for {
				if p == nil {
					mergeAttribute(n, "itemprop", "recipeInstructions")
					break
				}
				if p.Kind() == ast.KindBlockquote {
					break
				}
				p = p.Parent()
			}
		}
		if entering && n.Kind() == ast.KindHeading {
			heading := n.(*ast.Heading)
			if heading.Level == 1 {
				mergeAttribute(n, "itemprop", "name")
			}
			//if heading.Level < 6 {
			//	heading.Level++
			//}
		}
		if entering && n.Kind() == ast.KindBlockquote && firstBlockquote {
			firstBlockquote = false
			mergeAttribute(n, "itemprop", "description")
		}
		return ast.WalkContinue, nil
	})
}
