package connectors

import (
	"bytes"
	"strings"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/strikethrough"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
	"golang.org/x/net/html"
)

func newTrilliumConverter() *converter.Converter {
	return converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
			strikethrough.NewStrikethroughPlugin(),
			table.NewTablePlugin(),
			&trilliumPlugin{},
		),
	)
}

// triliumLangMap translates Trilium's internal MIME-based language identifiers
// to standard fenced-code-block language names.
var triliumLangMap = map[string]string{
	"text-x-csrc":          "c",
	"text-x-c++src":        "cpp",
	"text-x-python":        "python",
	"text-x-java":          "java",
	"text-x-javascript":    "javascript",
	"text-x-typescript":    "typescript",
	"text-x-go":            "go",
	"text-x-rust":          "rust",
	"text-x-shellscript":   "bash",
	"text-x-sh":            "bash",
	"text-x-yaml":          "yaml",
	"text-x-json":          "json",
	"text-x-xml":           "xml",
	"text-x-sql":           "sql",
	"text-x-markdown":      "markdown",
	"text-x-ruby":          "ruby",
	"text-x-php":           "php",
	"text-x-swift":         "swift",
	"text-x-kotlin":        "kotlin",
	"text-x-trilium-auto":  "mermaid",
	"application/javascript": "javascript",
	"application/typescript": "typescript",
	"application/json":     "json",
	"application/xml":      "xml",
}

type trilliumPlugin struct{}

func (p *trilliumPlugin) Name() string { return "trillium" }

func (p *trilliumPlugin) Init(conv *converter.Converter) error {
	// Must run before PriorityEarly (100) because the base plugin removes <input>
	// elements at PriorityEarly, which would destroy checkbox state before we read it.
	conv.Register.PreRenderer(p.preRender, converter.PriorityEarly-10)
	conv.Register.PostRenderer(p.postRenderTaskList, converter.PriorityStandard)
	conv.Register.RendererFor("figure", converter.TagTypeBlock, p.renderFigure, converter.PriorityStandard)
	conv.Register.RendererFor("span", converter.TagTypeInline, p.renderMathSpan, converter.PriorityStandard)
	return nil
}

// preRender fixes code language class names and transforms todo-list items.
func (p *trilliumPlugin) preRender(_ converter.Context, doc *html.Node) {
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "code":
				fixCodeLanguage(n)
			case "ul":
				if dom.HasClass(n, "todo-list") {
					transformTodoList(n)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
}

// fixCodeLanguage rewrites Trilium MIME-type class names (e.g. language-text-x-python)
// to standard names (e.g. language-python).
func fixCodeLanguage(n *html.Node) {
	for i, attr := range n.Attr {
		if attr.Key != "class" {
			continue
		}
		parts := strings.Fields(attr.Val)
		for j, part := range parts {
			trimmed := strings.TrimPrefix(part, "language-")
			if trimmed == part {
				continue
			}
			if mapped, ok := triliumLangMap[trimmed]; ok {
				parts[j] = "language-" + mapped
			}
		}
		n.Attr[i].Val = strings.Join(parts, " ")
	}
}

// transformTodoList rewrites CKEditor todo-list structure to plain <ul> with
// "[ ] " / "[x] " prefixes so the standard list renderer produces GFM task lists.
func transformTodoList(ul *html.Node) {
	for i, attr := range ul.Attr {
		if attr.Key == "class" {
			ul.Attr[i].Val = ""
		}
	}

	for li := ul.FirstChild; li != nil; li = li.NextSibling {
		if li.Type != html.ElementNode || li.Data != "li" {
			continue
		}

		checked := false
		descText := ""

		var find func(*html.Node)
		find = func(n *html.Node) {
			if n.Type == html.ElementNode {
				if n.Data == "input" {
					if dom.GetAttributeOr(n, "type", "") == "checkbox" {
						_, checked = dom.GetAttribute(n, "checked")
					}
				}
				if dom.HasClass(n, "todo-list__label__description") {
					descText = dom.CollectText(n)
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				find(c)
			}
		}
		find(li)

		for li.FirstChild != nil {
			li.RemoveChild(li.FirstChild)
		}

		prefix := "[ ] "
		if checked {
			prefix = "[x] "
		}
		li.AppendChild(&html.Node{
			Type: html.TextNode,
			Data: prefix + descText,
		})
	}
}

// postRenderTaskList unescapes the \[x] and \[ ] markers that the text
// escaper introduces for task-list prefixes written as plain text nodes.
func (p *trilliumPlugin) postRenderTaskList(_ converter.Context, content []byte) []byte {
	content = bytes.ReplaceAll(content, []byte(`\[x] `), []byte("[x] "))
	content = bytes.ReplaceAll(content, []byte(`\[ ] `), []byte("[ ] "))
	return content
}

// renderFigure unwraps <figure class="table"> so the table plugin can handle
// the inner <table>. Other <figure> elements are passed to the next renderer.
func (p *trilliumPlugin) renderFigure(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
	if !dom.HasClass(n, "table") {
		return converter.RenderTryNext
	}
	ctx.RenderChildNodes(ctx, w, n)
	return converter.RenderSuccess
}

// renderMathSpan converts <span class="math-tex">\(...\)</span> to $...$ and
// <span class="math-tex">\[...\]</span> to $$...$$.
func (p *trilliumPlugin) renderMathSpan(_ converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
	if !dom.HasClass(n, "math-tex") {
		return converter.RenderTryNext
	}

	text := dom.CollectText(n)

	if strings.HasPrefix(text, `\(`) && strings.HasSuffix(text, `\)`) {
		inner := strings.TrimSpace(text[2 : len(text)-2])
		w.WriteString("$")
		w.WriteString(inner)
		w.WriteString("$")
		return converter.RenderSuccess
	}

	if strings.HasPrefix(text, `\[`) && strings.HasSuffix(text, `\]`) {
		inner := strings.TrimSpace(text[2 : len(text)-2])
		w.WriteString("\n\n$$\n")
		w.WriteString(inner)
		w.WriteString("\n$$\n\n")
		return converter.RenderSuccess
	}

	return converter.RenderTryNext
}
