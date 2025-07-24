package formatter

import (
	"strings"
)

type FormatterNode struct {
	value    string
	children []*FormatterNode
}

type Formatter interface {
	Format() *FormatterNode
}

func Format(value string, children ...*FormatterNode) *FormatterNode {
	return &FormatterNode{
		value:    value,
		children: children,
	}
}

func FormatWithChildren[T Formatter](value string, children ...T) *FormatterNode {
	node := &FormatterNode{
		value:    value,
		children: make([]*FormatterNode, 0, len(children)),
	}

	for _, child := range children {
		node.children = append(node.children, child.Format())
	}

	return node
}

func (n *FormatterNode) Compact() string {
	if len(n.children) == 0 {
		return n.value
	}

	result := n.value + "("
	for i, child := range n.children {
		if i > 0 {
			result += ", "
		}
		result += child.Compact()
	}
	return result + ")"
}

func (n *FormatterNode) Detailed() string {
	return n.detailedWithIndent(0)
}

func (n *FormatterNode) detailedWithIndent(indent int) string {
	ind := strings.Repeat(" ", indent*2)
	rows := []string{
		ind + n.value,
	}

	for _, child := range n.children {
		rows = append(rows, child.detailedWithIndent(indent+1))
	}
	return strings.Join(rows, "\n")

}
