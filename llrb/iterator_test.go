package llrb

import (
	"reflect"
	"testing"
)

func TestAscendGreaterOrEqual(t *testing.T) {
	tree := New()
	tree.InsertNoReplace(Int(4), int(1))
	tree.InsertNoReplace(Int(6), int(2))
	tree.InsertNoReplace(Int(1), int(3))
	tree.InsertNoReplace(Int(3), int(4))
	var ary []Key
	tree.AscendGreaterOrEqual(Int(-1), func(i Key) bool {
		ary = append(ary, i)
		return true
	})
	expected := []Key{Int(1), Int(3), Int(4), Int(6)}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	tree.AscendGreaterOrEqual(Int(3), func(i Key) bool {
		ary = append(ary, i)
		return true
	})
	expected = []Key{Int(3), Int(4), Int(6)}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	tree.AscendGreaterOrEqual(Int(2), func(i Key) bool {
		ary = append(ary, i)
		return true
	})
	expected = []Key{Int(3), Int(4), Int(6)}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
}

func TestDescendLessOrEqual(t *testing.T) {
	tree := New()
	tree.InsertNoReplace(Int(4), int(1))
	tree.InsertNoReplace(Int(6), int(2))
	tree.InsertNoReplace(Int(1), int(3))
	tree.InsertNoReplace(Int(3), int(4))
	var ary []Key
	tree.DescendLessOrEqual(Int(10), func(i Key) bool {
		ary = append(ary, i)
		return true
	})
	expected := []Key{Int(6), Int(4), Int(3), Int(1)}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	tree.DescendLessOrEqual(Int(4), func(i Key) bool {
		ary = append(ary, i)
		return true
	})
	expected = []Key{Int(4), Int(3), Int(1)}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	tree.DescendLessOrEqual(Int(5), func(i Key) bool {
		ary = append(ary, i)
		return true
	})
	expected = []Key{Int(4), Int(3), Int(1)}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
}
