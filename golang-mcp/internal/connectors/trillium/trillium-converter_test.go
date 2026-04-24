package connectors

import (
	"os"
	"strings"
	"testing"
)

func convertHTML(t *testing.T, input string) string {
	t.Helper()
	md, err := triliumConv.ConvertString(input)
	if err != nil {
		t.Fatalf("ConvertString: %v", err)
	}
	return md
}

func TestConvertFile(t *testing.T) {
	input, err := os.ReadFile("../../../file.html")
	if err != nil {
		t.Skip("file.html not present:", err)
	}
	md := convertHTML(t, string(input))
	t.Log(md)
}

func TestMathInline(t *testing.T) {
	md := convertHTML(t, `<p>Inline: <span class="math-tex">\(E = mc^2\)</span></p>`)
	if !strings.Contains(md, "$E = mc^2$") {
		t.Errorf("expected inline math $E = mc^2$, got:\n%s", md)
	}
}

func TestMathBlock(t *testing.T) {
	md := convertHTML(t, `<span class="math-tex">\[\int_0^\infty e^{-x} dx = 1\]</span>`)
	if !strings.Contains(md, "$$") {
		t.Errorf("expected block math $$, got:\n%s", md)
	}
	if !strings.Contains(md, `\int_0^\infty e^{-x} dx = 1`) {
		t.Errorf("expected math content preserved, got:\n%s", md)
	}
}

func TestTodoListChecked(t *testing.T) {
	input := `<ul class="todo-list">
		<li><label class="todo-list__label">
			<input type="checkbox" checked="checked" disabled="disabled">
			<span class="todo-list__label__description">Done task</span>
		</label></li>
		<li><label class="todo-list__label">
			<input type="checkbox" disabled="disabled">
			<span class="todo-list__label__description">Pending task</span>
		</label></li>
	</ul>`
	md := convertHTML(t, input)
	if !strings.Contains(md, "[x] Done task") {
		t.Errorf("expected [x] Done task, got:\n%s", md)
	}
	if !strings.Contains(md, "[ ] Pending task") {
		t.Errorf("expected [ ] Pending task, got:\n%s", md)
	}
}

func TestTableConversion(t *testing.T) {
	input := `<figure class="table"><table>
		<thead><tr><th>Syntax</th><th>Description</th></tr></thead>
		<tbody><tr><td>Header</td><td>Title</td></tr></tbody>
	</table></figure>`
	md := convertHTML(t, input)
	if !strings.Contains(md, "|") {
		t.Errorf("expected GFM pipe table, got:\n%s", md)
	}
	if !strings.Contains(md, "Syntax") || !strings.Contains(md, "Description") {
		t.Errorf("expected table headers, got:\n%s", md)
	}
}

func TestCodeLanguageMapping(t *testing.T) {
	cases := []struct{ mime, want string }{
		{"text-x-python", "python"},
		{"text-x-csrc", "c"},
		{"text-x-trilium-auto", "mermaid"},
		{"text-x-go", "go"},
	}
	for _, tc := range cases {
		input := `<pre><code class="language-` + tc.mime + `">code here</code></pre>`
		md := convertHTML(t, input)
		if !strings.Contains(md, "```"+tc.want) {
			t.Errorf("mime %s: expected ```%s, got:\n%s", tc.mime, tc.want, md)
		}
	}
}

func TestStrikethrough(t *testing.T) {
	md := convertHTML(t, `<p><del>Strikethrough</del></p>`)
	if !strings.Contains(md, "~~Strikethrough~~") {
		t.Errorf("expected ~~Strikethrough~~, got:\n%s", md)
	}
}
