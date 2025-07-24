package formatter

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

func (n *FormatterNode) String() string {
	if len(n.children) == 0 {
		return n.value
	}

	result := n.value + "("
	for i, child := range n.children {
		if i > 0 {
			result += ", "
		}
		result += child.String()
	}
	return result + ")"
}
