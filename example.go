package main

import (
	"fmt"
	"math/rand"

	"github.com/slobdell/persistent-red-black/tree"
)

func main() {
	verifyDelete()
}

func verifyIntersection() {
	ints := []int{1, 2, 3, 4, 5, 6, 7, 10, 0}
	rb1 := tree.NewRedBlack[int](compare)
	for _, i := range ints {
		rb1 = rb1.Upsert(i)
	}

	ints2 := []int{11, 12, 13, 17, 11, 12, -1, 15, 10, 100, 2, 3}
	rb2 := tree.NewRedBlack[int](compare)
	for _, i := range ints2 {
		rb2 = rb2.Upsert(i)
	}

	rbFinal := rb1.Intersection(rb2)
	for it := rbFinal.Iterator(); it.HasElem(); it.Next() {
		i := it.Elem()
		fmt.Println(i)
	}
	if err := tree.ValidateInvariants[int](rbFinal); err != nil {
		panic(err)
	}
}

func verifyUnion() {
	ints := []int{1, 2, 3, 4, 5, 6, 7}
	rb1 := tree.NewRedBlack[int](compare)
	for _, i := range ints {
		rb1 = rb1.Upsert(i)
	}

	ints2 := []int{1, 2, 3, 7, 11, 12, -1, 15}
	rb2 := tree.NewRedBlack[int](compare)
	for _, i := range ints2 {
		rb2 = rb2.Upsert(i)
	}

	rbFinal := rb1.Union(rb2)
	for it := rbFinal.Iterator(); it.HasElem(); it.Next() {
		i := it.Elem()
		fmt.Println(i)
	}
	if err := tree.ValidateInvariants[int](rbFinal); err != nil {
		panic(err)
	}
}

func verifyDelete() {
	ints := []int{1, 2, 3, 4, 5, 6, 7}
	rb := tree.NewRedBlack[int](compare)
	for _, i := range ints {
		rb = rb.Upsert(i)
	}
	rb = rb.Delete(10)
	for it := rb.Iterator(); it.HasElem(); it.Next() {
		i := it.Elem()
		fmt.Println(i)
	}
	if err := tree.ValidateInvariants[int](rb); err != nil {
		panic(err)
	}
}

func verifySubtract() {
	ints := []int{1, 2, 3, 4, 5, 6, 7}
	rb1 := tree.NewRedBlack[int](compare)
	for _, i := range ints {
		rb1 = rb1.Upsert(i)
	}

	ints2 := []int{1, 6, 7, 8, 9, 10, 4}
	rb2 := tree.NewRedBlack[int](compare)
	for _, i := range ints2 {
		rb2 = rb2.Upsert(i)
	}

	rbFinal := rb1.Subtract(rb2)
	for it := rbFinal.Iterator(); it.HasElem(); it.Next() {
		i := it.Elem()
		fmt.Println(i)
	}
}

func verifySorted() {
	rb := tree.NewRedBlack[int](compare)
	var ints []int
	const sampleSize = 200000
	for i := 0; i < sampleSize; i++ {
		ints = append(ints, i)
	}
	rand.Shuffle(len(ints), func(i, j int) { ints[i], ints[j] = ints[j], ints[i] })
	for _, i := range ints {
		/*
			if idx%1000 == 0 {
				if err := tree.ValidateInvariants[int](rb); err != nil {
					panic(err)
				}
			}
		*/
		rb = rb.Upsert(i)
	}
	for it := rb.Iterator(); it.HasElem(); it.Next() {
		i := it.Elem()
		fmt.Println(i)
	}
	if err := tree.ValidateInvariants[int](rb); err != nil {
		panic(err)
	}
}

func compare(a, b int) int {
	if a < b {
		return -1
	}
	if a == b {
		return 0
	}
	return 1
}
