package main

import (
	"github.com/wttw/recipemark/site"
	"os"
)

func (b Build) Run() error {
	// fmt.Printf("Building %#v\n", b)
	builder := site.NewBuilder(os.DirFS(b.Source), os.DirFS(b.Assets), site.NewDestFS(b.Destination))
	return builder.Build()
}
