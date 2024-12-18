package ingredient

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Parser implements the goldmark parser.Parser interface
type Parser struct{}

// Trigger looks for the | separating an ingredient from a quantity
func (p *Parser) Trigger() []byte {
	return []byte{'|'}
}

// Parse implements the parser.Parser interface for pipe separated ingredients in markdown
func (p *Parser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	n := parent
	for n.Kind() != ast.KindListItem {
		n = n.Parent()
		if n == nil {
			return nil
		}
	}

	line, segment := block.PeekLine()
	fmt.Printf("line=[%s] segment=%#v\n", line, segment)
	pos := 1
	for ; pos < len(line); pos++ {
		if line[pos] == '|' {
			pos++
			break
		}
	}
	quantity := block.Value(text.NewSegment(segment.Start+1, segment.Start+pos-1))
	block.Advance(pos)
	return NewQuantity(strings.TrimSpace(string(quantity)))
}
