package helpers

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// ---- Math AST nodes ----

var KindMathInline = gast.NewNodeKind("MathInline")
var KindMathBlock = gast.NewNodeKind("MathBlock")

type MathInline struct {
	gast.BaseInline
	Formula []byte
}

func (n *MathInline) Kind() gast.NodeKind { return KindMathInline }
func (n *MathInline) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

type MathBlock struct {
	gast.BaseBlock
	Formula []byte
}

func (n *MathBlock) Kind() gast.NodeKind { return KindMathBlock }
func (n *MathBlock) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

// ---- Inline math parser: $...$ ----

type mathInlineParser struct{}

func (p *mathInlineParser) Trigger() []byte { return []byte{'$'} }

func (p *mathInlineParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	line, _ := block.PeekLine()
	if len(line) < 3 || line[0] != '$' || line[1] == '$' {
		return nil
	}
	end := bytes.IndexByte(line[1:], '$')
	if end < 0 {
		return nil
	}
	node := &MathInline{Formula: append([]byte{}, line[1:end+1]...)}
	block.Advance(end + 2)
	return node
}

// ---- Block math parser: $$\n...\n$$ ----

type mathBlockParser struct{}

func (p *mathBlockParser) Trigger() []byte { return []byte{'$'} }

func (p *mathBlockParser) Open(parent gast.Node, reader text.Reader, pc parser.Context) (gast.Node, parser.State) {
	line, _ := reader.PeekLine()
	if !bytes.Equal(bytes.TrimSpace(line), []byte("$$")) {
		return nil, parser.NoChildren
	}
	return &MathBlock{}, parser.Continue | parser.NoChildren
}

func (p *mathBlockParser) Continue(node gast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, seg := reader.PeekLine()
	if bytes.Equal(bytes.TrimSpace(line), []byte("$$")) {
		reader.Advance(seg.Len()) // consume closing delimiter
		return parser.Close
	}
	// goldmark auto-advances content lines; don't advance here
	node.(*MathBlock).Formula = append(node.(*MathBlock).Formula, bytes.TrimRight(line, "\r\n")...)
	return parser.Continue | parser.NoChildren
}

func (p *mathBlockParser) Close(node gast.Node, reader text.Reader, pc parser.Context) {
	n := node.(*MathBlock)
	n.Formula = bytes.TrimRight(n.Formula, "\r\n")
}

func (p *mathBlockParser) CanInterruptParagraph() bool { return false }
func (p *mathBlockParser) CanAcceptIndentedLine() bool { return false }

// ---- Trilium node renderer ----

type triliumRenderer struct{}

func (r *triliumRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindMathInline, r.renderMathInline)
	reg.Register(KindMathBlock, r.renderMathBlock)
	reg.Register(gast.KindList, r.renderList)
	reg.Register(gast.KindListItem, r.renderListItem)
	reg.Register(east.KindTaskCheckBox, r.renderTaskCheckBox)
}

func (r *triliumRenderer) renderMathInline(w util.BufWriter, _ []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		fmt.Fprintf(w, `<span class="math-tex">\(%s\)</span>`, node.(*MathInline).Formula)
	}
	return gast.WalkSkipChildren, nil
}

func (r *triliumRenderer) renderMathBlock(w util.BufWriter, _ []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		fmt.Fprintf(w, "<span class=\"math-tex\">\\[%s\\]</span>\n", node.(*MathBlock).Formula)
	}
	return gast.WalkSkipChildren, nil
}

// taskCheckBox returns the TaskCheckBox inside a ListItem (through TextBlock/Paragraph).
func taskCheckBox(listItem gast.Node) *east.TaskCheckBox {
	fc := listItem.FirstChild() // TextBlock or Paragraph
	if fc == nil {
		return nil
	}
	tc := fc.FirstChild()
	if tc != nil && tc.Kind() == east.KindTaskCheckBox {
		return tc.(*east.TaskCheckBox)
	}
	return nil
}

func (r *triliumRenderer) renderList(w util.BufWriter, _ []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*gast.List)
	tag := "ul"
	if n.IsOrdered() {
		tag = "ol"
	}
	isTaskList := !n.IsOrdered() && taskCheckBox(n.FirstChild()) != nil
	if entering {
		switch {
		case isTaskList:
			w.WriteString("<ul class=\"todo-list\">\n")
		case n.IsOrdered() && n.Start != 1:
			fmt.Fprintf(w, "<ol start=\"%d\">\n", n.Start)
		default:
			fmt.Fprintf(w, "<%s>\n", tag)
		}
	} else {
		fmt.Fprintf(w, "</%s>\n", tag)
	}
	return gast.WalkContinue, nil
}

func (r *triliumRenderer) renderListItem(w util.BufWriter, _ []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	tc := taskCheckBox(node)
	if tc == nil {
		if entering {
			w.WriteString("<li>\n")
		} else {
			w.WriteString("</li>\n")
		}
		return gast.WalkContinue, nil
	}
	if entering {
		checked := tc.IsChecked
		w.WriteString("<li>\n<label class=\"todo-list__label\">")
		if checked {
			w.WriteString(`<input type="checkbox" disabled="disabled" checked="checked">`)
		} else {
			w.WriteString(`<input type="checkbox" disabled="disabled">`)
		}
		w.WriteString(`<span class="todo-list__label__description">`)
	} else {
		w.WriteString("</span></label>\n</li>\n")
	}
	return gast.WalkContinue, nil
}

func (r *triliumRenderer) renderTaskCheckBox(_ util.BufWriter, _ []byte, _ gast.Node, _ bool) (gast.WalkStatus, error) {
	return gast.WalkSkipChildren, nil
}

// ---- Extension ----

type triliumMDExtension struct{}

func (e *triliumMDExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(util.Prioritized(&mathInlineParser{}, 500)),
		parser.WithBlockParsers(util.Prioritized(&mathBlockParser{}, 500)),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(util.Prioritized(&triliumRenderer{}, 1)),
	)
}

// ---- Public API ----

var triliumMDConverter = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		&triliumMDExtension{},
	),
	goldmark.WithRendererOptions(
		html.WithUnsafe(),
	),
)

func ConvertMDToHTML(md string) (string, error) {
	var buf bytes.Buffer
	if err := triliumMDConverter.Convert([]byte(md), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
