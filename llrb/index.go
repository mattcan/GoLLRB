package llrb

// IndexOfGreaterOrEqual returns the index of item as if the tree were a
// sorted array. If the same item does not exist in the tree, it returns
// the index of the closest item greater than item. If the tree does not
// contain any item greater than or equal to item, it returns t.Len().
func (t *LLRB) IndexOfGreaterOrEqual(item Item) int {
	var i int
	h := t.root
	for h != nil {
		switch {
		case less(item, h.Item):
			h = h.Left
		case less(h.Item, item):
			if h.Left != nil {
				i += h.Left.count
			}
			i++
			h = h.Right
		default:
			if h.Left != nil {
				i += h.Left.count
			}
			return i
		}
	}
	return i
}

// IndexOfLessOrEqual returns the index of item as if the tree were a
// sorted array. If the same item does not exist in the tree, it returns
// the index of the closest item less than item. If the tree does not
// contain any item less than or equal to item, it returns -1.
func (t *LLRB) IndexOfLessOrEqual(item Item) int {
	var i int
	h := t.root
	for h != nil {
		switch {
		case less(item, h.Item):
			h = h.Left
		case less(h.Item, item):
			if h.Left != nil {
				i += h.Left.count
			}
			i++
			h = h.Right
		default:
			if h.Left != nil {
				i += h.Left.count
			}
			return i
		}
	}
	return i - 1
}

// IndexOfLessOrEqual returns the index of item as if the tree were a
// sorted array. The second return variable indicates success. It is
// false if the requested item does not exist in the tree.
func (t *LLRB) IndexOf(item Item) (int, bool) {
	var i int
	h := t.root
	for h != nil {
		switch {
		case less(item, h.Item):
			h = h.Left
		case less(h.Item, item):
			if h.Left != nil {
				i += h.Left.count
			}
			i++
			h = h.Right
		default:
			if h.Left != nil {
				i += h.Left.count
			}
			return i, true
		}
	}
	return 0, false
}
