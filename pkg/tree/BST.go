package tree

import (
	. "gocode/types"
	"math"
)

/**
 * 98. Validate Binary Search Tree
 */
func isValidBST(root *TreeNode) bool {
    var validate func(*TreeNode, int, int) bool
    validate = func(node *TreeNode, lower int, upper int) bool {
        if node == nil {
            return true
        }

        if (lower < node.Val) && (node.Val < upper) {
            return validate(node.Left, lower, node.Val) && validate(node.Right, node.Val, upper)
        } else {
            return false
        }
    }

    return validate(root, math.MinInt, math.MaxInt)
}

/**
 * 230. Kth Smallest Element in a BST
 */
func kthSmallest(root *TreeNode, k int) int {
    stack := []*TreeNode{}
    curr := root
    for root != nil {
        for curr != nil {
            stack = append(stack, curr)
            curr = curr.Left
        }
        curr = stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        k--
        if k == 0 {
            return curr.Val
        }

        curr = curr.Right
    }
    return -1
}

/**
 * 235. Lowest Common Ancestor of a Binary Search Tree
 */
func lowestCommonAncestorRecursive(root, p, q *TreeNode) *TreeNode {
    if p.Val < root.Val && q.Val < root.Val {
        return lowestCommonAncestorRecursive(root.Left, p, q)
    }
    if p.Val > root.Val && q.Val > root.Val {
        return lowestCommonAncestorRecursive(root.Right, p, q)
    }
    return root
}

func lowestCommonAncestorIterative(root, p, q *TreeNode) *TreeNode {
    for root != nil {
        if p.Val < root.Val && q.Val < root.Val {
            root = root.Left
        } else if p.Val > root.Val && q.Val > root.Val {
            root = root.Right
        } else {
            return root
        }
    }
    return nil
}

/**
 * 236. Lowest Common Ancestor of a Binary Tree
 */
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
    if root == nil || root == p || root == q {
        return root
    }

    // Process left and then right nodes
    left := lowestCommonAncestor(root.Left, p, q)
    right := lowestCommonAncestor(root.Right, p, q)

    // if both left and right aren't null, that means we found the targets on both sides of trees, means we need to return root
    if left != nil && right != nil {
        return root
    }
    // if we find in left, return left
    if left != nil {
        return left
    }
    return right
}
