package main

type Global struct {
	Config string `help:"Alternate configuration file"`
}

type Build struct {
	Source      string `help:"The source directory where recipes are found"`
	Assets      string `help:"The source directory for templates and assets"`
	Destination string `help:"The destination directory to build to"`
}

type Serve struct {
	Listen      string `help:"The address:port to listen on when serving"`
	Destination string `help:"The destination directory to build to"`
}

type Config struct {
	Source      string `help:"The source directory where recipes are found"`
	Assets      string `help:"The source directory for templates and assets"`
	Listen      string `help:"The address:port to listen on when serving"`
	Destination string `help:"The destination directory to build to"`
}
