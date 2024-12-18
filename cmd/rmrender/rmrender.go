package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/wttw/recipemark"
)

func main() {
	source, err := os.ReadFile("/Users/steve/recipes/caponata/caponata.md")
	if err != nil {
		log.Fatal(err)
	}
	p := recipemark.NewParser()
	recipe, err := p.Parse(source)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(recipe.Html)
	meta, err := json.MarshalIndent(recipe, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", meta)
}
