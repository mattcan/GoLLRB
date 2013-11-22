# GoLLRB with indexing

GoLLRB is a Left-Leaning Red-Black (LLRB) implementation of 2-3 balanced binary
search trees in Go Language.

This fork adds three functions:

	IndexOfGreaterOrEqual(Item) int
	IndexOfLessOrEqual(Item) int
	IndexOf(Item) (int, bool)

## Rationale

You can use these functions to determine how many items you will be iterating
over, before you perform the actual iteration. Perhaps you need to send a
header with the number of items. Or maybe you want to allocate a slice of the
correct size before iterating.

## Overhead

To provide these functions requires one int per node, so memory consumption grows
by 8\*n bytes. Index maintenance code adds a small CPU overhead to
insert and delete operations which do not affect asymptotic complexity. Insert and
delete are still O(log n), as are the new IndexOf... functions.

## Installation

`go get github.com/foobaz/GoLLRB/llrb`

## Example

This code:

        package main

        import (
                "fmt"
                "github.com/foobaz/GoLLRB/llrb"
        )   

        func main() {
                var t llrb.LLRB
                for i := 0; i <= 4; i += 2 { 
                        item := llrb.Int(i)
                        t.ReplaceOrInsert(item)
                }   

                fmt.Printf("%d items in tree\n", t.Len())

                for i := 0; i <= 4; i++ {
                        item := llrb.Int(i)
                        exact, ok := t.IndexOf(item)
                        if ok {
                                fmt.Printf("%d is at %d\n", i, exact)
                        } else {
                                min := t.IndexOfLessOrEqual(item)
                                max := t.IndexOfGreaterOrEqual(item)
                                fmt.Printf("%d is between %d and %d\n", i, min, max)
                        }   
                }   
        }

prints this output:

        3 items in tree
        0 is at 0
        1 is between 0 and 1
        2 is at 1
        3 is between 1 and 2
        4 is at 2

## About

GoLLRB was written by [Petar Maymounkov](http://pdos.csail.mit.edu/~petar/). 

The indexing additions were added by William MacKay.
