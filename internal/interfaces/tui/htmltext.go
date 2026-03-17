package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"golang.org/x/net/html"
)

var skipTags = map[string]bool{
	"script": true,
	"style":  true,
}

var blockTags = map[string]bool{
	"p": true, "div": true, "br": true,
	"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
	"li": true, "ul": true, "ol": true,
	"blockquote": true, "pre": true, "article": true, "section": true,
	"header": true, "footer": true, "main": true, "nav": true,
}

var headingTags = map[string]bool{
	"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
}

var boldTags = map[string]bool{
	"strong": true, "b": true,
}

var italicTags = map[string]bool{
	"em": true, "i": true,
}

var boldStyle = lipgloss.NewStyle().Bold(true)
var italicStyle = lipgloss.NewStyle().Italic(true)

func htmlToText(raw string, width int) string {
	if width < 20 {
		width = 20
	}

	doc, err := html.Parse(strings.NewReader(raw))
	if err != nil {
		return wordwrap.String(raw, width)
	}

	var sb strings.Builder
	walkNode(&sb, doc, false, false)

	// Deduplicate consecutive blank lines
	lines := strings.Split(sb.String(), "\n")
	out := make([]string, 0, len(lines))
	blankCount := 0
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			blankCount++
			if blankCount <= 1 {
				out = append(out, "")
			}
		} else {
			blankCount = 0
			out = append(out, wordwrap.String(l, width))
		}
	}

	return strings.TrimSpace(strings.Join(out, "\n"))
}

// walkNode recursively visits HTML nodes and writes text to sb.
func walkNode(sb *strings.Builder, n *html.Node, isBold, isItalic bool) {
	if n.Type == html.ElementNode {
		tag := n.Data
		if skipTags[tag] {
			return
		}

		if blockTags[tag] {
			// Ensure block elements start on a new line
			s := sb.String()
			if len(s) > 0 && !strings.HasSuffix(s, "\n") {
				sb.WriteByte('\n')
			}
		}

		if tag == "li" {
			sb.WriteString("• ")
		}

		childBold := isBold || boldTags[tag] || headingTags[tag]
		childItalic := isItalic || italicTags[tag]

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walkNode(sb, c, childBold, childItalic)
		}

		if blockTags[tag] || headingTags[tag] {
			if !strings.HasSuffix(sb.String(), "\n") {
				sb.WriteByte('\n')
			}
		}
		return
	}

	if n.Type == html.TextNode {
		text := n.Data
		if strings.TrimSpace(text) == "" && text != " " {
			// Collapse pure-whitespace text nodes to a single space only if
			// the builder doesn't already end with a newline or space.
			s := sb.String()
			if len(s) > 0 && s[len(s)-1] != '\n' && s[len(s)-1] != ' ' {
				sb.WriteByte(' ')
			}
			return
		}

		switch {
		case isBold && isItalic:
			sb.WriteString(boldStyle.Italic(true).Render(text))
		case isBold:
			sb.WriteString(boldStyle.Render(text))
		case isItalic:
			sb.WriteString(italicStyle.Render(text))
		default:
			sb.WriteString(text)
		}
		return
	}

	// For document / other node types, just recurse.
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkNode(sb, c, isBold, isItalic)
	}
}
