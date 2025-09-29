package main

import (
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// MarkdownElement represents a parsed element from markdown with hierarchy info.
type MarkdownElement struct {
	Type           string             `json:"type"` // e.g., "heading", "list", "codeblock", "text"
	Content        string             `json:"content"`
	Depth          int                `json:"depth"` // hierarchy depth for headings/lists
	Children       []*MarkdownElement `json:"children,omitempty"`
	Path           []string           `json:"path,omitempty"` // path of headings leading to this element
	Lang           string             `json:"lang,omitempty"` // optional: code block language
	Index          int                `json:"index,omitempty"`
	SecondaryIndex int                `json:"secondary_index,omitempty"`
}

// parseMarkdownFile parses the markdown file preserving heading/list hierarchy and code blocks.
// This implementation uses goldmark to build an AST and walks it to produce logical units:
// headings, list items, paragraphs (text), and code blocks.
func parseMarkdownFile(filePath string) ([]*MarkdownElement, error) {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	md := goldmark.New()
	reader := text.NewReader(src)
	doc := md.Parser().Parse(reader)

	var elements []*MarkdownElement
	var headingStack []string
	listDepth := 0

	// Helper: recursively extract text content from a node's children
	var extractText func(node ast.Node) string
	extractText = func(node ast.Node) string {
		var b strings.Builder
		for c := node.FirstChild(); c != nil; c = c.NextSibling() {
			// collect plain text nodes
			if t, ok := c.(*ast.Text); ok {
				b.Write(t.Segment.Value(src))
			} else {
				// recurse into children (covers emphasis, strong, links, etc.)
				b.WriteString(extractText(c))
			}
		}
		return b.String()
	}

	// Helper: collect lines from a code block node
	collectCode := func(node ast.Node) string {
		var b strings.Builder
		for i := 0; i < node.Lines().Len(); i++ {
			seg := node.Lines().At(i)
			b.Write(seg.Value(src))
		}
		return b.String()
	}

	// Walk the AST and build elements
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		switch node := n.(type) {
		case *ast.Heading:
			if entering {
				level := node.Level
				content := strings.TrimSpace(extractText(node))
				// Adjust heading stack to reflect current level
				if level <= len(headingStack) {
					headingStack = headingStack[:level-1]
				}
				headingStack = append(headingStack, content)
				elements = append(elements, &MarkdownElement{
					Type:    "heading",
					Content: content,
					Depth:   level,
					Path:    append([]string{}, headingStack...),
				})
			}
			// no special action on exit (heading handled on enter)
		case *ast.List:
			if entering {
				listDepth++
			} else {
				if listDepth > 0 {
					listDepth--
				}
			}
		case *ast.ListItem:
			if entering {
				// Compute list depth by counting ancestor List nodes (more robust than a global counter)
				depth := 0
				for p := node.Parent(); p != nil; p = p.Parent() {
					if _, ok := p.(*ast.List); ok {
						depth++
					}
				}

				// Prefer the paragraph child as the list item's content, and avoid including nested list text.
				content := ""
				for c := node.FirstChild(); c != nil; c = c.NextSibling() {
					if para, ok := c.(*ast.Paragraph); ok {
						content = strings.TrimSpace(extractText(para))
						break
					}
					if txt, ok := c.(*ast.Text); ok && content == "" {
						content = strings.TrimSpace(string(txt.Segment.Value(src)))
					}
				}
				if content == "" {
					// Fallback: concatenate non-list children
					var b strings.Builder
					for c := node.FirstChild(); c != nil; c = c.NextSibling() {
						if _, ok := c.(*ast.List); ok {
							continue
						}
						b.WriteString(extractText(c))
					}
					content = strings.TrimSpace(b.String())
				}

				// Build list path by appending "list-item" depth times to heading path
				listPath := append([]string{}, headingStack...)
				for i := 0; i < depth; i++ {
					listPath = append(listPath, "list-item")
				}
				elements = append(elements, &MarkdownElement{
					Type:    "list",
					Content: content,
					Depth:   depth,
					Path:    listPath,
				})
			}
		case *ast.FencedCodeBlock:
			if entering {
				content := collectCode(node)
				lang := string(node.Language(src))
				elements = append(elements, &MarkdownElement{
					Type:    "codeblock",
					Content: content,
					Depth:   len(headingStack),
					Path:    append([]string{}, headingStack...),
					Lang:    lang,
				})
			}
			// skip walking into lines
			return ast.WalkSkipChildren, nil
		case *ast.CodeBlock:
			if entering {
				content := collectCode(node)
				elements = append(elements, &MarkdownElement{
					Type:    "codeblock",
					Content: content,
					Depth:   len(headingStack),
					Path:    append([]string{}, headingStack...),
				})
			}
			return ast.WalkSkipChildren, nil
		case *ast.Paragraph:
			if entering {
				content := strings.TrimSpace(extractText(node))
				// avoid creating empty text nodes
				if content != "" {
					elements = append(elements, &MarkdownElement{
						Type:    "text",
						Content: content,
						Depth:   len(headingStack),
						Path:    append([]string{}, headingStack...),
					})
				}
			}
			// let walker continue into children (extractText already reads children)
		}
		return ast.WalkContinue, nil
	})

	return elements, nil
}
