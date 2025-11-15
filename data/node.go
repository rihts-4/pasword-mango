package data

// Node represent each character
type Node struct {
	//this is a single letter stored for example letter a,b,c,d,etc
	Char string
	//store all children  of a node
	//that is from a-z
	//a slice of Nodes(and each child will also have 26 children)
	Children [26]*Node
	// IsEnd will be true if the node represents the end of a word
	IsEnd bool
}

// / NewNode this will be used to initialize a new node with 26 children
// /each child should first be initialized to nil
func NewNode(char string) *Node {
	node := &Node{Char: char}
	node.IsEnd = false
	return node
}
