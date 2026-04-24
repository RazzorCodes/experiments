package helpers

import (
	"strings"
	"testing"
)

func TestMDToHTMLBold(t *testing.T) {
	out, err := ConvertMDToHTML("**bold**")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "<strong>bold</strong>") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestMDToHTMLMathInline(t *testing.T) {
	out, err := ConvertMDToHTML("here is $x^2 + y^2$ inline")
	if err != nil {
		t.Fatal(err)
	}
	want := `<span class="math-tex">\(x^2 + y^2\)</span>`
	if !strings.Contains(out, want) {
		t.Errorf("want %q in output, got: %s", want, out)
	}
}

func TestMDToHTMLMathBlock(t *testing.T) {
	out, err := ConvertMDToHTML("$$\nx^2 + y^2 = 1\n$$")
	if err != nil {
		t.Fatal(err)
	}
	want := `<span class="math-tex">\[x^2 + y^2 = 1\]</span>`
	if !strings.Contains(out, want) {
		t.Errorf("want %q in output, got: %s", want, out)
	}
}

func TestMDToHTMLTaskList(t *testing.T) {
	md := "- [x] done\n- [ ] todo\n"
	out, err := ConvertMDToHTML(md)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `class="todo-list"`) {
		t.Errorf("missing todo-list class: %s", out)
	}
	if !strings.Contains(out, `checked="checked"`) {
		t.Errorf("missing checked attribute: %s", out)
	}
	if !strings.Contains(out, `class="todo-list__label__description"`) {
		t.Errorf("missing label description span: %s", out)
	}
}

func TestMDToHTMLTable(t *testing.T) {
	md := "| A | B |\n|---|---|\n| 1 | 2 |\n"
	out, err := ConvertMDToHTML(md)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "<table>") || !strings.Contains(out, "<td>") {
		t.Errorf("missing table tags: %s", out)
	}
}

func TestMDToHTMLStrikethrough(t *testing.T) {
	out, err := ConvertMDToHTML("~~struck~~")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "<del>struck</del>") {
		t.Errorf("unexpected output: %s", out)
	}
}
