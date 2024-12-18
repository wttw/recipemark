package site

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/wttw/recipemark"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DestFS interface {
	WriteFile(path string, content []byte) error
}

func NewDestFS(dir string) DestFS {
	return &destFS{dir}
}

type destFS struct {
	Dir string
}

func (d *destFS) WriteFile(path string, content []byte) error {
	perm := os.FileMode(0644)
	fullPath := filepath.Join(d.Dir, path)
	dir := filepath.Dir(fullPath)
	err := os.MkdirAll(dir, os.FileMode(0755)|os.ModeDir)
	if err != nil {
		return err
	}
	return os.WriteFile(fullPath, content, perm)
}

type SourceInfo struct {
	ModTime     time.Time
	Title       string
	TitleParsed bool
	Parsed      bool
}

type SourceWalker interface {
	Visit(key string, val SourceInfo) error
}

type SourceDB interface {
	Set(key string, value SourceInfo) error
	Get(key string) (SourceInfo, bool, error)
	Walk(SourceWalker) error
}

type SourceDBMap map[string]SourceInfo

func (s SourceDBMap) Set(key string, value SourceInfo) error {
	s[key] = value
	return nil
}

func (s SourceDBMap) Get(key string) (SourceInfo, bool, error) {
	val, ok := s[key]
	return val, ok, nil
}

func (s SourceDBMap) Walk(w SourceWalker) error {
	for key, val := range s {
		err := w.Visit(key, val)
		if err != nil {
			return err
		}
	}
	return nil
}

type Builder struct {
	Source    fs.FS
	Assets    fs.FS
	Dest      DestFS
	SourceDB  SourceDB
	Parser    *recipemark.Parser
	Templates *template.Template
}

func NewBuilder(recipes, assets fs.FS, dest DestFS) *Builder {
	return &Builder{
		Source:   recipes,
		Assets:   assets,
		Dest:     dest,
		SourceDB: SourceDBMap{},
		Parser:   recipemark.NewParser(),
	}
}

func (b *Builder) Build() error {
	// Find all the files and invalidate any stored metadata if they've changed
	err := fs.WalkDir(b.Source, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		fmt.Println()
		key := path
		si, ok, err := b.SourceDB.Get(key)
		if err != nil {
			return err
		}
		fi, err := d.Info()
		if err != nil {
			return err
		}
		if ok && fi.ModTime().After(si.ModTime) {
			ok = false
		}
		if !ok {
			err = b.SourceDB.Set(key, SourceInfo{
				ModTime: fi.ModTime(),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Printf("assets: %#v\n", b.Assets)
	// Setup our templates
	b.Templates, err = template.New("").Funcs(b.TemplateFunctions()).ParseFS(b.Assets, "*.tpl")
	if err != nil {
		return fmt.Errorf("failed to parse html templates: %w", err)
	}
	err = b.SourceDB.Walk(b)
	if err != nil {
		return err
	}
	return nil
}

func (b *Builder) Visit(key string, val SourceInfo) error {
	md, err := fs.ReadFile(b.Source, key)
	if err != nil {
		return fmt.Errorf("failed to read recipe file '%s': %w", key, err)
	}
	recipe, err := b.Parser.Parse(md)
	if err != nil {
		return fmt.Errorf("failed to parse recipe file '%s': %w", key, err)
	}
	fmt.Printf("Parsed %s: %v\n", key, recipe.Name)
	fmt.Printf("Meta: %#v\n", recipe.Meta)
	var buff bytes.Buffer
	err = b.Templates.ExecuteTemplate(&buff, "single.tpl", recipe)
	if err != nil {
		return fmt.Errorf("failed to render template 'single.tpl' for '%s': %w", key, err)
	}
	// dir, file := filepath.Split(key)
	err = b.Dest.WriteFile(strings.TrimSuffix(key, ".md")+".html", buff.Bytes())
	return err
}

func (b *Builder) TemplateFunctions() template.FuncMap {
	ret := sprig.HtmlFuncMap()
	ret["ImageSet"] = b.imageSet
	return ret
}

var widthRe = regexp.MustCompile(`^([0-9]+)w$`)

type imageVariant struct {
	Variant  string
	Filename string
}

func (b *Builder) imageSet(file string, flags string) (template.HTMLAttr, error) {
	variants := strings.Fields(flags)
	suffix := filepath.Ext(file)
	basefile := strings.TrimSuffix(file, suffix)
	filenames := make([]imageVariant, 0, len(variants))
	for _, v := range variants {
		matches := widthRe.FindStringSubmatch(v)
		if matches != nil {
			w, err := strconv.Atoi(matches[1])
			if err != nil {
				return "", fmt.Errorf("invalid imageSet flag '%s': %w", v, err)
			}
			filenames = append(filenames, imageVariant{
				Variant:  fmt.Sprintf("%dw", w),
				Filename: fmt.Sprintf("%s-%dw%s", basefile, w, suffix),
			})
		}
	}
	if len(filenames) == 0 {
		return "", errors.New("no valid variants in imageSet")
	}
	return template.HTMLAttr(fmt.Sprintf(`src="%s"`, filenames[0].Filename)), nil
}
