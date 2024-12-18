package recipemark

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark/ast"
	"html/template"
	"strings"

	hashtag "github.com/13rac1/goldmark-hashtag"
	wikilink "github.com/13rac1/goldmark-wikilink"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	"github.com/wttw/recipemark/ingredient"
)

type Parser struct {
	md goldmark.Markdown
}

// NewParser returns a preconfigured goldmark.Markdown parser
func NewParser() *Parser {
	p := &Parser{
		md: goldmark.New(
			goldmark.WithExtensions(
				extension.Table,
				extension.Footnote,
				emoji.Emoji,
				wikilink.New(),
				hashtag.New(),
				meta.Meta,
				ingredient.Extender{},
			),
		),
	}
	return p
}

func (p *Parser) Parse(source []byte) (Recipe, error) {
	r := Recipe{}
	reader := text.NewReader(source)
	context := parser.NewContext()
	node := p.md.Parser().Parse(reader, parser.WithContext(context))
	r.Doc = node.OwnerDocument()
	metaData, err := meta.TryGet(context)
	if err != nil {
		return Recipe{}, fmt.Errorf("invalid metadata: %v", err)
	}
	if metaData == nil {
		metaData = map[string]interface{}{}
	}
	r.Doc.Dump(source, 0)

	// Walk top level of the AST, to broadly categorize into sections
	n := r.Doc.FirstChild()
	var seenIngredients = false
	var section string
	for n != nil {
		var kind string
		if n.Kind() == ast.KindList {
			seenIngredients = true
		}
		if n.Kind() == ast.KindHeading {
			heading := n.(*ast.Heading)
			switch heading.Level {
			case 1:
				kind = "h1"
				_, ok := metaData["name"]
				if !ok {
					metaData["name"] = string(n.Text(source))
				}
			case 2:
				kind = "h2"
			case 3:
				kind = "h3"
			default:
				kind = "h"
			}
		} else {
			kind = "romance"
		}
		if !seenIngredients {
			section = "romance"
			_, ok := metaData["description"]
			if !ok {
				switch n.Kind() {
				case ast.KindParagraph, ast.KindBlockquote:
					var buff bytes.Buffer
					err = p.md.Renderer().Render(&buff, source, n)
					if err != nil {
						return Recipe{}, err
					}
					metaData["description"] = buff.String()
				}
			}
		} else {
			switch n.Kind() {
			case ast.KindHeading:
				// already handled above
			case ast.KindList:
				kind = "ingredients"
				section = "ingredients"
			case ast.KindParagraph:
				kind = "step"
				section = "method"
			case ast.KindBlockquote:
				kind = "note"
				section = "method"
			default:
				kind = "content"
			}
		}
		var buff bytes.Buffer
		err := p.md.Renderer().Render(&buff, source, n)
		if err != nil {
			return Recipe{}, err
		}
		r.Chunks = append(r.Chunks, &Chunk{
			Section: section,
			Type:    kind,
			Content: template.HTML(buff.Bytes()), //nolint:gosec
		})
		n = n.NextSibling()
		// fmt.Printf("%s: %s\n\n", kind, buff.String())
	}
	// Attach headers to the following section
	for i := len(r.Chunks); i > 0; i-- {
		switch r.Chunks[i-1].Type {
		case "h2", "h3", "h":
			r.Chunks[i-1].Section = r.Chunks[i].Section
		}
	}
	for _, c := range r.Chunks {
		fmt.Printf("[%s] [%s]\n", c.Section, c.Type)
	}
	var buff bytes.Buffer
	err = p.md.Renderer().Render(&buff, source, r.Doc)
	if err != nil {
		return Recipe{}, err
	}
	r.Html = template.HTML(buff.String()) //nolint:gosec

	_ = ast.Walk(r.Doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if n.Kind() == ingredient.KindIngredient && entering {
			r.Ingredients = append(r.Ingredients, string(n.Text(source)))
		}
		if n.Kind() == ast.KindImage {
			_, ok := metaData["image"]
			if !ok {
				im := n.(*ast.Image)
				if len(im.Destination) != 0 {
					metaData["image"] = string(im.Destination)
				}
			}
		}
		return ast.WalkContinue, nil
	})
	r.Meta = metaData
	var description string
	handleMeta(&r.Name, "name", metaData)
	handleMeta(&description, "description", metaData)
	if description != "" {
		r.Description = template.HTML(description) //nolint:gosec
	}
	handleMeta(&r.Author, "author", metaData)
	handleMeta(&r.Image, "image", metaData)
	handleMeta(&r.Cuisine, "cuisine", metaData)
	handleMeta(&r.Category, "category", metaData)
	handleMeta(&r.Yield, "yield", metaData)
	handleMeta(&r.Method, "method", metaData)
	handleMeta(&r.PrepTime, "prepTime", metaData)
	handleMeta(&r.CookTime, "cookTime", metaData)
	handleMeta(&r.TotalTime, "totalTime", metaData)
	return r, nil
}

func handleMeta(dst *string, name string, meta map[string]interface{}) {
	val, ok := meta[name]
	if ok {
		v, ok := val.(string)
		if ok {
			*dst = v
			return
		}
		*dst = fmt.Sprintf("%v", val)
	}
}

func hasAttributeTag(n ast.Node, attr, tag string) bool {
	a, ok := n.AttributeString(attr)
	var s string
	if !ok {
		return false
	}

	switch v := a.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return false
	}

	tags := strings.Fields(s)
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}
