// Copyright 2010 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A Left-Leaning Red-Black (LLRB) implementation of 2-3 balanced binary search trees,
// based on the following work:
//
//   http://www.cs.princeton.edu/~rs/talks/LLRB/08Penn.pdf
//   http://www.cs.princeton.edu/~rs/talks/LLRB/LLRB.pdf
//   http://www.cs.princeton.edu/~rs/talks/LLRB/Java/RedBlackBST.java
//
//  2-3 trees (and the run-time equivalent 2-3-4 trees) are the de facto standard BST
//  algoritms found in implementations of Python, Java, and other libraries. The LLRB
//  implementation of 2-3 trees is a recent improvement on the traditional implementation,
//  observed and documented by Robert Sedgewick.
//
package llrb

// Tree is a Left-Leaning Red-Black (LLRB) implementation of 2-3 trees
type LLRB struct {
	root *Node
}

type Node struct {
	Key
	Value
	Left, Right *Node // Pointers to left and right child nodes
	count       int
	Black       bool // If set, the color of the link (incoming from the parent) is black
	// In the LLRB, new nodes are always red, hence the zero-value for node
}

type Key interface {
	Less(than Key) bool
}

type Value interface {
	Update(value Value)
}

//
func less(x, y Key) bool {
	if x == pinf {
		return false
	}
	if x == ninf {
		return true
	}
	return x.Less(y)
}

// Inf returns an Item that is "bigger than" any other item, if sign is positive.
// Otherwise  it returns an Item that is "smaller than" any other item.
func Inf(sign int) Key {
	if sign == 0 {
		panic("sign")
	}
	if sign > 0 {
		return pinf
	}
	return ninf
}

var (
	ninf = nInf{}
	pinf = pInf{}
)

type nInf struct{}

func (nInf) Less(Key) bool {
	return true
}

type pInf struct{}

func (pInf) Less(Key) bool {
	return false
}

// New() allocates a new tree
func New() *LLRB {
	return &LLRB{}
}

// SetRoot sets the root node of the tree.
// It is intended to be used by functions that deserialize the tree.
func (t *LLRB) SetRoot(r *Node) {
	t.root = r
}

// Root returns the root node of the tree.
// It is intended to be used by functions that serialize the tree.
func (t *LLRB) Root() *Node {
	return t.root
}

// Len returns the number of nodes in the tree.
func (t *LLRB) Len() int {
	if t.root != nil {
		return t.root.count
	}

	return 0
}

// Has returns true if the tree contains an element whose order is the same as that of key.
func (t *LLRB) Has(key Key) bool {
	k, _ := t.Get(key)
	return k != nil
}

// Get retrieves an element from the tree whose order is the same as that of key.
func (t *LLRB) Get(key Key) (Key, Value) {
	h := t.root
	for h != nil {
		switch {
		case less(key, h.Key):
			h = h.Left
		case less(h.Key, key):
			h = h.Right
		default:
			return h.Key, h.Value
		}
	}
	return nil, nil
}

// Min returns the minimum element in the tree.
func (t *LLRB) Min() (Key, Value) {
	h := t.root
	if h == nil {
		return nil, nil
	}
	for h.Left != nil {
		h = h.Left
	}
	return h.Key, h.Value
}

// Max returns the maximum element in the tree.
func (t *LLRB) Max() (Key, Value) {
	h := t.root
	if h == nil {
		return nil, nil
	}
	for h.Right != nil {
		h = h.Right
	}
	return h.Key, h.Value
}

func (t *LLRB) UpdateOrInsert(k Key, v Value) {
	if k == nil {
		panic("Cannot update or insert nil key")
	}

	t.root = t.updateOrInsert(t.root, k, v)
	t.root.Black = true
}

func (t *LLRB) updateOrInsert(h *Node, k Key, v Value) *Node {
	if h == nil {
		return newNode(k, v)
	}

	h = walkDownRot23(h)

	if less(k, h.Key) {
		h.Left = t.updateOrInsert(h.Left, k, v)
	} else if less(h.Key, k) {
		h.Right = t.updateOrInsert(h.Right, k, v)
	} else {
		h.Value.Update(v)
	}

	h.count++

	h = walkUpRot23(h)

	return h
}


// ReplaceOrInsert inserts item into the tree. If an existing
// element has the same order, it is removed from the tree and returned.
func (t *LLRB) ReplaceOrInsert(k Key, v Value) Key {
	if k == nil {
		panic("Cannot replace or insert nil key")
	}
	var replaced Key
	t.root, replaced = t.replaceOrInsert(t.root, k, v)
	t.root.Black = true
	return replaced
}

func (t *LLRB) replaceOrInsert(h *Node, k Key, v Value) (*Node, Key) {
	if h == nil {
		return newNode(k, v), k
	}

	h = walkDownRot23(h)

	var replaced Key
	if less(k, h.Key) { // BUG
		h.Left, replaced = t.replaceOrInsert(h.Left, k, v)
		if replaced == nil {
			h.count++
		}
	} else if less(h.Key, k) {
		h.Right, replaced = t.replaceOrInsert(h.Right, k, v)
		if replaced == nil {
			h.count++
		}
	} else {
		replaced, h.Key, h.Value = h.Key, k, v
	}

	h = walkUpRot23(h)

	return h, replaced
}

// InsertNoReplace inserts item into the tree. If an existing
// element has the same order, both elements remain in the tree.
func (t *LLRB) InsertNoReplace(k Key, v Value) {
	if k == nil {
		panic("Cannot insert nil key")
	}
	t.root = t.insertNoReplace(t.root, k, v)
	t.root.Black = true
}

func (t *LLRB) insertNoReplace(h *Node, k Key, v Value) *Node {
	if h == nil {
		return newNode(k, v)
	}

	h = walkDownRot23(h)

	if less(k, h.Key) {
		h.Left = t.insertNoReplace(h.Left, k, v)
	} else {
		h.Right = t.insertNoReplace(h.Right, k, v)
	}
	h.count++

	return walkUpRot23(h)
}

// Rotation driver routines for 2-3 algorithm

func walkDownRot23(h *Node) *Node { return h }

func walkUpRot23(h *Node) *Node {
	if isRed(h.Right) && !isRed(h.Left) {
		h = rotateLeft(h)
	}

	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}

	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}

	return h
}

// Rotation driver routines for 2-3-4 algorithm

func walkDownRot234(h *Node) *Node {
	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}

	return h
}

func walkUpRot234(h *Node) *Node {
	if isRed(h.Right) && !isRed(h.Left) {
		h = rotateLeft(h)
	}

	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}

	return h
}

// DeleteMin deletes the minimum element in the tree and returns the
// deleted item or nil otherwise.
func (t *LLRB) DeleteMin() Key {
	var deleted Key
	t.root, deleted = deleteMin(t.root)
	if t.root != nil {
		t.root.Black = true
	}
	return deleted
}

// deleteMin code for LLRB 2-3 trees
func deleteMin(h *Node) (*Node, Key) {
	if h == nil {
		return nil, nil
	}
	if h.Left == nil {
		return nil, h.Key
	}

	if !isRed(h.Left) && !isRed(h.Left.Left) {
		h = moveRedLeft(h)
	}

	var deleted Key
	h.Left, deleted = deleteMin(h.Left)
	if deleted != nil {
		h.count--
	}

	return fixUp(h), deleted
}

// DeleteMax deletes the maximum element in the tree and returns
// the deleted item or nil otherwise
func (t *LLRB) DeleteMax() Key {
	var deleted Key
	t.root, deleted = deleteMax(t.root)
	if t.root != nil {
		t.root.Black = true
	}
	return deleted
}

func deleteMax(h *Node) (*Node, Key) {
	if h == nil {
		return nil, nil
	}
	if isRed(h.Left) {
		h = rotateRight(h)
	}
	if h.Right == nil {
		return nil, h.Key
	}
	if !isRed(h.Right) && !isRed(h.Right.Left) {
		h = moveRedRight(h)
	}
	var deleted Key
	h.Right, deleted = deleteMax(h.Right)
	if deleted != nil {
		h.count--
	}

	return fixUp(h), deleted
}

// Delete deletes an item from the tree whose key equals key.
// The deleted item is return, otherwise nil is returned.
func (t *LLRB) Delete(k Key) Key {
	var deleted Key
	t.root, deleted = t.delete(t.root, k)
	if t.root != nil {
		t.root.Black = true
	}
	return deleted
}

func (t *LLRB) delete(h *Node, k Key) (*Node, Key) {
	var deleted Key
	if h == nil {
		return nil, nil
	}
	if less(k, h.Key) {
		if h.Left == nil { // item not present. Nothing to delete
			return h, nil
		}
		if !isRed(h.Left) && !isRed(h.Left.Left) {
			h = moveRedLeft(h)
		}
		h.Left, deleted = t.delete(h.Left, k)
	} else {
		if isRed(h.Left) {
			h = rotateRight(h)
		}
		// If @k equals @h.Key and no right children at @h
		if !less(h.Key, k) && h.Right == nil {
			return nil, h.Key
		}
		// PETAR: Added 'h.Right != nil' below
		if h.Right != nil && !isRed(h.Right) && !isRed(h.Right.Left) {
			h = moveRedRight(h)
		}
		// If @k equals @h.Key, and (from above) 'h.Right != nil'
		if !less(h.Key, k) {
			var subDeleted Key
			h.Right, subDeleted = deleteMin(h.Right)
			if subDeleted == nil {
				panic("logic")
			}
			deleted, h.Key = h.Key, subDeleted
		} else { // Else, @k is bigger than @h.Key
			h.Right, deleted = t.delete(h.Right, k)
		}
	}
	if deleted != nil {
		h.count--
	}

	return fixUp(h), deleted
}

// Internal node manipulation routines

func newNode(key Key, value Value) *Node { return &Node{Key: key, Value: value, count: 1} }

func isRed(h *Node) bool {
	if h == nil {
		return false
	}
	return !h.Black
}

func rotateLeft(h *Node) *Node {
	x := h.Right
	if x.Black {
		panic("rotating a black link")
	}

	x.count, h.count = h.count, h.count-x.count
	if x.Left != nil {
		h.count += x.Left.count
	}

	h.Right = x.Left
	x.Left = h
	x.Black = h.Black
	h.Black = false
	return x
}

func rotateRight(h *Node) *Node {
	x := h.Left
	if x.Black {
		panic("rotating a black link")
	}

	x.count, h.count = h.count, h.count-x.count
	if x.Right != nil {
		h.count += x.Right.count
	}

	h.Left = x.Right
	x.Right = h
	x.Black = h.Black
	h.Black = false
	return x
}

// REQUIRE: Left and Right children must be present
func flip(h *Node) {
	h.Black = !h.Black
	h.Left.Black = !h.Left.Black
	h.Right.Black = !h.Right.Black
}

// REQUIRE: Left and Right children must be present
func moveRedLeft(h *Node) *Node {
	flip(h)
	if isRed(h.Right.Left) {
		h.Right = rotateRight(h.Right)
		h = rotateLeft(h)
		flip(h)
	}
	return h
}

// REQUIRE: Left and Right children must be present
func moveRedRight(h *Node) *Node {
	flip(h)
	if isRed(h.Left.Left) {
		h = rotateRight(h)
		flip(h)
	}
	return h
}

func fixUp(h *Node) *Node {
	if isRed(h.Right) {
		h = rotateLeft(h)
	}

	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}

	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}

	return h
}
