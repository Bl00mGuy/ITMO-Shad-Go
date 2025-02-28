//go:build !solution

package treeiter

func DoInOrder[Node interface {
	Left() *Node
	Right() *Node
}](root *Node, cb func(node *Node)) {
	if root == nil {
		return
	}

	DoInOrder((*root).Left(), cb)
	cb(root)
	DoInOrder((*root).Right(), cb)
}
