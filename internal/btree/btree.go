package btree

type TreeNode struct {
	Value []byte
	Left  *TreeNode
	Right *TreeNode
}

func New(val []byte) *TreeNode {
	return &TreeNode{
		Value: val,
		Left:  nil,
		Right: nil,
	}
}
