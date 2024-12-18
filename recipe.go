package recipemark

import (
	"html/template"
	"time"

	"github.com/yuin/goldmark/ast"
)

type Ingredient struct {
	Quantity    string
	Name        string
	Preparation template.HTML
}

type Chunk struct {
	Section    string
	Type       string
	Ingredient Ingredient
	Content    template.HTML
}

type Recipe struct {
	Doc         *ast.Document `json:"-"`
	Html        template.HTML `json:"-"`
	Published   time.Time
	Image       string
	Tags        []string
	Name        string
	Description template.HTML
	Author      string
	Cuisine     string
	Category    string
	Yield       string
	Method      string
	PrepTime    string
	CookTime    string
	TotalTime   string
	Meta        map[string]interface{}
	Ingredients []string
	Chunks      []*Chunk
}
