package tree

type color bool

const (
	red   = true
	black = false
)

type Compare[T any] func(this, that T) int

type RedBlackTree[T any] struct {
	compare Compare[T]
	root    *node[T]
}

// Iterator is an iterator over treap elements. It can be used like this:
//
//     for it := v.Iterator(); it.HasElem(); it.Next() {
//         elem := it.Elem()
//         // do something with elem...
//     }
type Iterator[T any] interface {
	// Elem returns the element at the current position.
	Elem() T
	// HasElem returns whether the iterator is pointing to an element.
	HasElem() bool
	// Next moves the iterator to the next position.
	Next()
}

type node[T any] struct {
	item  T
	c     color
	left  *node[T]
	right *node[T]
}

func (n node[T]) copyWithEntry(item T) *node[T] {
	n.item = item
	return &n
}

func (n node[T]) copyWithLeft(left *node[T]) *node[T] {
	n.left = left
	return &n
}

func (n node[T]) copyWithRight(right *node[T]) *node[T] {
	n.right = right
	return &n
}

func (n *node[T]) copyWithColor(c color) *node[T] {
	if n.c == c {
		return n
	}
	return &node[T]{
		item:  n.item,
		c:     c,
		left:  n.left,
		right: n.right,
	}
}

func newNode[T any](item T) *node[T] {
	return &node[T]{
		item: item,
		c:    red,
	}
}

func (n *node[T]) withMaybeNewLeft(c Compare[T], inserting *node[T]) *node[T] {
	if inserting == nil {
		return n
	}
	if n.left == nil {
		return n.copyWithLeft(inserting).balance()
	}
	return n.left.upsert(c, inserting)
}

func (n *node[T]) withMaybeNewRight(c Compare[T], inserting *node[T]) *node[T] {
	if inserting == nil {
		return n
	}
	if n.right == nil {
		return n.copyWithRight(inserting).balance()
	}
	return n.right.upsert(c, inserting)
}

func (n *node[T]) combinedChildren(c Compare[T]) *node[T] {
	if n.left == nil {
		return n.right
	}
	return n.left.upsert(c, n.right)
}

func (n *node[T]) upsert(c Compare[T], inserting *node[T]) *node[T] {
	if inserting == nil {
		return n
	}
	cmp := c(inserting.item, n.item)
	if cmp == 0 {
		return n.copyWithEntry(inserting.item).upsert(c, inserting.right).balance().upsert(c, inserting.left).balance()
	}
	if cmp == -1 {
		if n.left == nil {
			return n.copyWithLeft(inserting).balance()
		}
		return n.copyWithLeft(n.left.upsert(c, inserting.copyWithRight(nil).balance())).balance().upsert(c, inserting.right)
	}
	if n.right == nil {
		return n.copyWithRight(inserting).balance()
	}
	return n.copyWithRight(n.right.upsert(c, inserting.copyWithLeft(nil).balance())).balance().upsert(c, inserting.left)
}

func (n *node[T]) intersection(c Compare[T], intersecting *node[T]) *node[T] {
	if n == nil {
		return nil
	}
	if intersecting == nil {
		return nil
	}
	cmp := c(intersecting.item, n.item)
	if cmp == 0 {
		rightIntersection := n.right.intersection(c, intersecting.right).balance()
		leftIntersection := n.left.intersection(c, intersecting.left).balance()
		return n.copyWithRight(rightIntersection).balance().copyWithLeft(leftIntersection).balance()
	}
	if cmp == -1 {
		leftIntersect := n.left.intersection(c, intersecting.copyWithRight(nil).balance())
		if leftIntersect == nil {
			return n.intersection(c, intersecting.right)
		}
		return leftIntersect.copyWithRight(n.intersection(c, intersecting.right)).balance()
	}
	rightIntersect := n.right.intersection(c, intersecting.copyWithLeft(nil).balance())
	if rightIntersect == nil {
		return n.intersection(c, intersecting.left)
	}
	return rightIntersect.copyWithLeft(n.intersection(c, intersecting.left)).balance()
}

func (n *node[T]) subtract(c Compare[T], subtracting *node[T]) *node[T] {
	if subtracting == nil {
		return n
	}
	cmp := c(subtracting.item, n.item)
	if cmp == 0 {
		return n.combinedChildren(c).subtract(c, subtracting.right).subtract(c, subtracting.left)
	}
	if cmp == -1 {
		if n.left == nil {
			return n
		}
		return n.copyWithLeft(n.left.subtract(c, subtracting.copyWithRight(nil).balance())).balance().subtract(c, subtracting.right)
	}
	if n.right == nil {
		return n
	}
	return n.copyWithRight(n.right.subtract(c, subtracting.copyWithLeft(nil).balance())).balance().subtract(c, subtracting.left)
}

func (n *node[T]) balance() *node[T] {
	if n == nil {
		return nil
	}
	if n.c == red {
		return n
	}
	// 4 cases:
	// https://www.cs.tufts.edu/comp/150FP/archive/chris-okasaki/redblack99.pdf
	// page 3
	// top case
	if n.left != nil && n.left.isRed() && n.left.right != nil && n.left.right.isRed() {
		return &node[T]{
			item: n.left.right.item,
			c:    red,
			left: &node[T]{
				item:  n.left.item,
				c:     black,
				left:  n.left.left,
				right: n.left.right.left,
			},
			right: &node[T]{
				item:  n.item,
				c:     black,
				left:  n.left.right.right,
				right: n.right,
			},
		}
	}
	// left case
	if n.left != nil && n.left.isRed() && n.left.left != nil && n.left.left.isRed() {
		return &node[T]{
			item: n.left.item,
			c:    red,
			left: n.left.left.copyWithColor(black),
			right: &node[T]{
				item:  n.item,
				c:     black,
				left:  n.left.right,
				right: n.right,
			},
		}
	}
	// right case
	if n.right != nil && n.right.isRed() && n.right.right != nil && n.right.right.isRed() {
		return &node[T]{
			item: n.right.item,
			c:    red,
			left: &node[T]{
				item:  n.item,
				c:     black,
				left:  n.left,
				right: n.right.left,
			},
			right: n.right.right.copyWithColor(black),
		}
	}
	// botttom case
	if n.right != nil && n.right.isRed() && n.right.left != nil && n.right.left.isRed() {
		return &node[T]{
			item: n.right.left.item,
			c:    red,
			left: &node[T]{
				item:  n.item,
				c:     black,
				left:  n.left,
				right: n.right.left.left,
			},
			right: &node[T]{
				item:  n.right.item,
				c:     black,
				left:  n.right.left.right,
				right: n.right.right,
			},
		}
	}
	return n
}

func (n *node[T]) isRed() bool {
	return n.c == red
}

/*
Invariant 1. No red node has a red parent.
Invariant 2. Every path from the root to an empty node contains the same number
of black nodes.

Empty nodes are black
*/

func NewRedBlack[T any](compareFn Compare[T]) *RedBlackTree[T] {
	return &RedBlackTree[T]{
		root:    nil,
		compare: compareFn,
	}
}

func (r *RedBlackTree[T]) Upsert(item T) *RedBlackTree[T] {
	if r.root == nil {
		return &RedBlackTree[T]{
			root: &node[T]{
				item: item,
				c:    black,
			},
			compare: r.compare,
		}
	}
	return &RedBlackTree[T]{
		root:    r.root.upsert(r.compare, newNode(item)).copyWithColor(black),
		compare: r.compare,
	}
}

func (r *RedBlackTree[T]) Delete(item T) *RedBlackTree[T] {
	if r.root == nil {
		return r
	}
	return &RedBlackTree[T]{
		root:    r.root.subtract(r.compare, newNode(item)).copyWithColor(black),
		compare: r.compare,
	}
}

func (r *RedBlackTree[T]) Subtract(other *RedBlackTree[T]) *RedBlackTree[T] {
	if r.root == nil {
		return r
	}
	return &RedBlackTree[T]{
		root:    r.root.subtract(r.compare, other.root).copyWithColor(black),
		compare: r.compare,
	}
}

func (r *RedBlackTree[T]) Union(other *RedBlackTree[T]) *RedBlackTree[T] {
	// note that union will overwrite shared keys with the values in other
	if r.root == nil {
		return other
	}
	return &RedBlackTree[T]{
		root:    r.root.upsert(r.compare, other.root).copyWithColor(black),
		compare: r.compare,
	}
}

func (r *RedBlackTree[T]) Intersection(other *RedBlackTree[T]) *RedBlackTree[T] {
	if r.root == nil {
		return nil
	}
	return &RedBlackTree[T]{
		root:    r.root.intersection(r.compare, other.root).copyWithColor(black),
		compare: r.compare,
	}

}

func (r *RedBlackTree[T]) Iterator() Iterator[T] {
	n := r.root
	if n == nil {
		empty := nodeIterator[T]{current: nil}
		return &empty
	}
	var stack []*node[T]
	iter := nodeIterator[T]{
		current: n,
	}

	for n.left != nil {
		n = n.left
		stack = append(stack, iter.current)
		iter.current = n
	}
	iter.unprocStack = stack
	return &iter
}

type nodeIterator[T any] struct {
	unprocStack []*node[T]
	current     *node[T]
}

func (n *nodeIterator[T]) Elem() T {
	if n.current == nil {
		var zeroVal T
		return zeroVal
	}
	return n.current.item
}

func (n *nodeIterator[T]) HasElem() bool {
	return n.current != nil
}

func (n *nodeIterator[T]) Next() {
	if n.current == nil {
		return
	}

	var cursor *node[T]
	if n.current != nil {
		cursor = n.current.right
	}
	for cursor != nil {
		n.unprocStack = append(n.unprocStack, cursor)
		cursor = cursor.left
	}

	n.current = nil
	if len(n.unprocStack) >= 1 {
		n.current = n.unprocStack[len(n.unprocStack)-1]
		n.unprocStack = n.unprocStack[:len(n.unprocStack)-1]
	}
}
