package main

import (
	"errors"
	"fmt"
	"sort"
	"text/template"
)

type Template struct {
	Text               string
	RequiresNumeric    bool
	RequiresComparable bool
	RequiresOrdered    bool
}

func getTemplate(name string) (result *template.Template, err error) {
	if isProjectionMethod(name) {
		return getProjectionTemplate(name)
	}
	return getStandardTemplate(name)
}

func getHeaderTemplate() *template.Template {
	return template.Must(template.New("header").Parse(header))
}

func getStandardTemplate(name string) (result *template.Template, err error) {
	t, found := standardTemplates[name]
	if found {
		result = template.Must(template.New(name).Parse(t.Text))
	} else {
		err = fmt.Errorf("%s is not a known method", name)
	}
	return
}

func isStandardMethod(s string) bool {
	_, ok := standardTemplates[s]
	return ok
}

func getStandardMethodKeys() (result []string) {
	for k := range standardTemplates {
		result = append(result, k)
	}
	sort.Strings(result)
	return
}

func getProjectionTemplate(name string) (result *template.Template, err error) {
	t, found := ProjectionMethods[name]
	if found {
		result = template.Must(template.New(name).Parse(t.Text))
	} else {
		err = errors.New(fmt.Sprintf("%s is not a known projection method", name))
	}
	return
}

func isProjectionMethod(s string) bool {
	_, ok := ProjectionMethods[s]
	return ok
}

func getProjectionMethodKeys() (result []string) {
	for k := range ProjectionMethods {
		result = append(result, k)
	}
	sort.Strings(result)
	return
}

func getSortSupportTemplate() *template.Template {
	return template.Must(template.New("sortSupport").Parse(sortSupport))
}

const header = `// This file was auto-generated using github.com/clipperhouse/gen
// Modifying this file is not recommended as it will likely be overwritten in the future

// Sort functions are a modification of http://golang.org/pkg/sort/#Sort
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package {{.Package.Name}}
{{if gt (len .Imports) 0}}
import ({{range .Imports}}
	"{{.}}"
{{end}})
{{end}}
// {{.Plural}} is a slice of type {{.Pointer}}{{.Name}}, for use with gen methods below. Use this type where you would use []{{.Pointer}}{{.Name}}. (This is required because slices cannot be method receivers.)
type {{.Plural}} []{{.Pointer}}{{.Name}}
`

var standardTemplates = map[string]*Template{

	"All": &Template{
		Text: `
// All verifies that all elements of {{.Plural}} return true for the passed func. See: http://clipperhouse.github.io/gen/#All
func (rcv {{.Plural}}) All(fn func({{.Pointer}}{{.Name}}) bool) bool {
	for _, v := range rcv {
		if !fn(v) {
			return false
		}
	}
	return true
}
`},

	"Any": &Template{
		Text: `
// Any verifies that one or more elements of {{.Plural}} return true for the passed func. See: http://clipperhouse.github.io/gen/#Any
func (rcv {{.Plural}}) Any(fn func({{.Pointer}}{{.Name}}) bool) bool {
	for _, v := range rcv {
		if fn(v) {
			return true
		}
	}
	return false
}
`},

	"Count": &Template{
		Text: `
// Count gives the number elements of {{.Plural}} that return true for the passed func. See: http://clipperhouse.github.io/gen/#Count
func (rcv {{.Plural}}) Count(fn func({{.Pointer}}{{.Name}}) bool) (result int) {
	for _, v := range rcv {
		if fn(v) {
			result++
		}
	}
	return
}
`},

	"Distinct": &Template{
		Text: `
// Distinct returns a new {{.Plural}} slice whose elements are unique. See: http://clipperhouse.github.io/gen/#Distinct
func (rcv {{.Plural}}) Distinct() (result {{.Plural}}) {
	appended := make(map[{{.Pointer}}{{.Name}}]bool)
	for _, v := range rcv {
		if !appended[v] {
			result = append(result, v)
			appended[v] = true
		}
	}
	return result
}
`,
		RequiresComparable: true,
	},

	"DistinctBy": &Template{
		Text: `
// DistinctBy returns a new {{.Plural}} slice whose elements are unique, where equality is defined by a passed func. See: http://clipperhouse.github.io/gen/#DistinctBy
func (rcv {{.Plural}}) DistinctBy(equal func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) (result {{.Plural}}) {
	for _, v := range rcv {
		eq := func(_app {{.Pointer}}{{.Name}}) bool {
			return equal(v, _app)
		}
		if !result.Any(eq) {
			result = append(result, v)
		}
	}
	return result
}
`},

	"Each": &Template{
		Text: `
// Each iterates over {{.Plural}} and executes the passed func against each element. See: http://clipperhouse.github.io/gen/#Each
func (rcv {{.Plural}}) Each(fn func({{.Pointer}}{{.Name}})) {
	for _, v := range rcv {
		fn(v)
	}
}
`},

	"First": &Template{
		Text: `
// First returns the first element that returns true for the passed func. Returns error if no elements return true. See: http://clipperhouse.github.io/gen/#First
func (rcv {{.Plural}}) First(fn func({{.Pointer}}{{.Name}}) bool) (result {{.Pointer}}{{.Name}}, err error) {
	for _, v := range rcv {
		if fn(v) {
			result = v
			return
		}
	}
	err = errors.New("no {{.Plural}} elements return true for passed func")
	return
}
`},

	"MaxBy": &Template{
		Text: `
// MaxBy returns an element of {{.Plural}} containing the maximum value, when compared to other elements using a passed func defining ‘less’. In the case of multiple items being equally maximal, the last such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#MaxBy
//
// (Note: this is implemented by negating the passed ‘less’ func, effectively testing ‘greater than or equal to’.)
func (rcv {{.Plural}}) MaxBy(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) (result {{.Pointer}}{{.Name}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine the MaxBy of an empty slice")
		return
	}
	m := 0
	for i := 1; i < l; i++ {
		if !less(rcv[i], rcv[m]) {
			m = i
		}
	}
	result = rcv[m]
	return
}
`},

	"MinBy": &Template{
		Text: `
// MinBy returns an element of {{.Plural}} containing the minimum value, when compared to other elements using a passed func defining ‘less’. In the case of multiple items being equally minimal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#MinBy
func (rcv {{.Plural}}) MinBy(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) (result {{.Pointer}}{{.Name}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine the Min of an empty slice")
		return
	}
	m := 0
	for i := 1; i < l; i++ {
		if less(rcv[i], rcv[m]) {
			m = i
		}
	}
	result = rcv[m]
	return
}
`},

	"Single": &Template{
		Text: `
// Single returns exactly one element of {{.Plural}} that returns true for the passed func. Returns error if no or multiple elements return true. See: http://clipperhouse.github.io/gen/#Single
func (rcv {{.Plural}}) Single(fn func({{.Pointer}}{{.Name}}) bool) (result {{.Pointer}}{{.Name}}, err error) {
	var candidate {{.Pointer}}{{.Name}}
	found := false
	for _, v := range rcv {
		if fn(v) {
			if found {
				err = errors.New("multiple {{.Plural}} elements return true for passed func")
				return
			}
			candidate = v
			found = true
		}
	}
	if found {
		result = candidate
	} else {
		err = errors.New("no {{.Plural}} elements return true for passed func")
	}
	return
}
`},

	"Where": &Template{
		Text: `
// Where returns a new {{.Plural}} slice whose elements return true for func. See: http://clipperhouse.github.io/gen/#Where
func (rcv {{.Plural}}) Where(fn func({{.Pointer}}{{.Name}}) bool) (result {{.Plural}}) {
	for _, v := range rcv {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
`},

	"SortBy": &Template{
		Text: `
// SortBy returns a new ordered {{.Plural}} slice, determined by a func defining ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv {{.Plural}}) SortBy(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) {{.Plural}} {
	result := make({{.Plural}}, len(rcv))
	copy(result, rcv)
	// Switch to heapsort if depth of 2*ceil(lg(n+1)) is reached.
	n := len(result)
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSort{{.Plural}}(result, less, 0, n, maxDepth)
	return result
}
`},

	"IsSortedBy": &Template{
		Text: `
// IsSortedBy reports whether an instance of {{.Plural}} is sorted, using the pass func to define ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv {{.Plural}}) IsSortedBy(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) bool {
	n := len(rcv)
	for i := n - 1; i > 0; i-- {
		if less(rcv[i], rcv[i-1]) {
			return false
		}
	}
	return true
}
`},

	"SortByDesc": &Template{
		Text: `
// SortByDesc returns a new, descending-ordered {{.Plural}} slice, determined by a func defining ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
//
// (Note: this is implemented by negating the passed ‘less’ func, effectively testing ‘greater than or equal to’.)
func (rcv {{.Plural}}) SortByDesc(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) {{.Plural}} {
	greaterOrEqual := func(a, b {{.Pointer}}{{.Name}}) bool {
		return !less(a, b)
	}
	return rcv.SortBy(greaterOrEqual)
}
`},

	"IsSortedDesc": &Template{
		Text: `
// IsSortedDesc reports whether an instance of {{.Plural}} is sorted in descending order, using the pass func to define ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv {{.Plural}}) IsSortedByDesc(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) bool {
	greaterOrEqual := func(a, b {{.Pointer}}{{.Name}}) bool {
		return !less(a, b)
	}
	return rcv.IsSortedBy(greaterOrEqual)
}
`},
}

const sortSupport = `
// Sort support methods

func swap{{.Plural}}(rcv {{.Plural}}, a, b int) {
	rcv[a], rcv[b] = rcv[b], rcv[a]
}

// Insertion sort
func insertionSort{{.Plural}}(rcv {{.Plural}}, less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && less(rcv[j], rcv[j-1]); j-- {
			swap{{.Plural}}(rcv, j, j-1)
		}
	}
}

// siftDown implements the heap property on rcv[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDown{{.Plural}}(rcv {{.Plural}}, less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && less(rcv[first+child], rcv[first+child+1]) {
			child++
		}
		if !less(rcv[first+root], rcv[first+child]) {
			return
		}
		swap{{.Plural}}(rcv, first+root, first+child)
		root = child
	}
}

func heapSort{{.Plural}}(rcv {{.Plural}}, less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown{{.Plural}}(rcv, less, i, hi, first)
	}

	// Pop elements, largest first, into end of rcv.
	for i := hi - 1; i >= 0; i-- {
		swap{{.Plural}}(rcv, first, first+i)
		siftDown{{.Plural}}(rcv, less, lo, i, first)
	}
}

// Quicksort, following Bentley and McIlroy,
// Engineering a Sort Function, SP&E November 1993.

// medianOfThree moves the median of the three values rcv[a], rcv[b], rcv[c] into rcv[a].
func medianOfThree{{.Plural}}(rcv {{.Plural}}, less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if less(rcv[m1], rcv[m0]) {
		swap{{.Plural}}(rcv, m1, m0)
	}
	if less(rcv[m2], rcv[m1]) {
		swap{{.Plural}}(rcv, m2, m1)
	}
	if less(rcv[m1], rcv[m0]) {
		swap{{.Plural}}(rcv, m1, m0)
	}
	// now rcv[m0] <= rcv[m1] <= rcv[m2]
}

func swapRange{{.Plural}}(rcv {{.Plural}}, a, b, n int) {
	for i := 0; i < n; i++ {
		swap{{.Plural}}(rcv, a+i, b+i)
	}
}

func doPivot{{.Plural}}(rcv {{.Plural}}, less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's Ninther, median of three medians of three.
		s := (hi - lo) / 8
		medianOfThree{{.Plural}}(rcv, less, lo, lo+s, lo+2*s)
		medianOfThree{{.Plural}}(rcv, less, m, m-s, m+s)
		medianOfThree{{.Plural}}(rcv, less, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThree{{.Plural}}(rcv, less, lo, m, hi-1)

	// Invariants are:
	//	rcv[lo] = pivot (set up by ChoosePivot)
	//	rcv[lo <= i < a] = pivot
	//	rcv[a <= i < b] < pivot
	//	rcv[b <= i < c] is unexamined
	//	rcv[c <= i < d] > pivot
	//	rcv[d <= i < hi] = pivot
	//
	// Once b meets c, can swap the "= pivot" sections
	// into the middle of the slice.
	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if less(rcv[b], rcv[pivot]) { // rcv[b] < pivot
				b++
			} else if !less(rcv[pivot], rcv[b]) { // rcv[b] = pivot
				swap{{.Plural}}(rcv, a, b)
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if less(rcv[pivot], rcv[c-1]) { // rcv[c-1] > pivot
				c--
			} else if !less(rcv[c-1], rcv[pivot]) { // rcv[c-1] = pivot
				swap{{.Plural}}(rcv, c-1, d-1)
				c--
				d--
			} else {
				break
			}
		}
		if b >= c {
			break
		}
		// rcv[b] > pivot; rcv[c-1] < pivot
		swap{{.Plural}}(rcv, b, c-1)
		b++
		c--
	}

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	n := min(b-a, a-lo)
	swapRange{{.Plural}}(rcv, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRange{{.Plural}}(rcv, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSort{{.Plural}}(rcv {{.Plural}}, less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSort{{.Plural}}(rcv, less, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot{{.Plural}}(rcv, less, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort{{.Plural}}(rcv, less, a, mlo, maxDepth)
			a = mhi // i.e., quickSort{{.Plural}}(rcv, mhi, b)
		} else {
			quickSort{{.Plural}}(rcv, less, mhi, b, maxDepth)
			b = mlo // i.e., quickSort{{.Plural}}(rcv, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSort{{.Plural}}(rcv, less, a, b)
	}
}
`

var ProjectionMethods = map[string]*Template{
	"Aggregate": &Template{
		Text: `
// {{.MethodName}} iterates over {{.Parent.Plural}}, operating on each element while maintaining ‘state’. See: http://clipperhouse.github.io/gen/#Aggregate
func (rcv {{.Parent.Plural}}) {{.MethodName}}(fn func({{.Type}}, {{.Parent.Pointer}}{{.Parent.Name}}) {{.Type}}) (result {{.Type}}) {
	for _, v := range rcv {
		result = fn(result, v)
	}
	return
}
`},

	"Average": &Template{
		Text: `
// {{.MethodName}} sums {{.Type}} over all elements and divides by len({{.Parent.Plural}}). See: http://clipperhouse.github.io/gen/#Average
func (rcv {{.Parent.Plural}}) {{.MethodName}}(fn func({{.Parent.Pointer}}{{.Parent.Name}}) {{.Type}}) (result {{.Type}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine {{.MethodName}} of zero-length {{.Parent.Plural}}")
		return
	}
	for _, v := range rcv {
		result += fn(v)
	}
	result = result / {{.Type}}(l)
	return
}
`,
		RequiresNumeric: true,
	},

	"GroupBy": &Template{
		Text: `
// {{.MethodName}} groups elements into a map keyed by {{.Type}}. See: http://clipperhouse.github.io/gen/#GroupBy
func (rcv {{.Parent.Plural}}) {{.MethodName}}(fn func({{.Parent.Pointer}}{{.Parent.Name}}) {{.Type}}) map[{{.Type}}]{{.Parent.Plural}} {
	result := make(map[{{.Type}}]{{.Parent.Plural}})
	for _, v := range rcv {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}
`,
		RequiresComparable: true,
	},

	"Max": &Template{
		Text: `
// {{.MethodName}} selects the largest value of {{.Type}} in {{.Parent.Plural}}. Returns error on {{.Parent.Plural}} with no elements. See: http://clipperhouse.github.io/gen/#MaxCustom
func (rcv {{.Parent.Plural}}) {{.MethodName}}(fn func({{.Parent.Pointer}}{{.Parent.Name}}) {{.Type}}) (result {{.Type}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine {{.MethodName}} of zero-length {{.Parent.Plural}}")
		return
	}
	result = fn(rcv[0])
	if l > 1 {
		for _, v := range rcv[1:] {
			f := fn(v)
			if f > result {
				result = f
			}
		}
	}
	return
}
`,
		RequiresOrdered: true,
	},

	"Min": &Template{
		Text: `
// {{.MethodName}} selects the least value of {{.Type}} in {{.Parent.Plural}}. Returns error on {{.Parent.Plural}} with no elements. See: http://clipperhouse.github.io/gen/#MinCustom
func (rcv {{.Parent.Plural}}) {{.MethodName}}(fn func({{.Parent.Pointer}}{{.Parent.Name}}) {{.Type}}) (result {{.Type}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine {{.MethodName}} of zero-length {{.Parent.Plural}}")
		return
	}
	result = fn(rcv[0])
	if l > 1 {
		for _, v := range rcv[1:] {
			f := fn(v)
			if f < result {
				result = f
			}
		}
	}
	return
}
`,
		RequiresOrdered: true,
	},

	"Select": &Template{
		Text: `
// {{.MethodName}} returns a slice of {{.Type}} in {{.Parent.Plural}}, projected by passed func. See: http://clipperhouse.github.io/gen/#Select
func (rcv {{.Parent.Plural}}) {{.MethodName}}(fn func({{.Parent.Pointer}}{{.Parent.Name}}) {{.Type}}) (result []{{.Type}}) {
	for _, v := range rcv {
		result = append(result, fn(v))
	}
	return
}
`,
	},

	"Sum": &Template{
		Text: `
// {{.MethodName}} sums {{.Type}} over elements in {{.Parent.Plural}}. See: http://clipperhouse.github.io/gen/#Sum
func (rcv {{.Parent.Plural}}) {{.MethodName}}(fn func({{.Parent.Pointer}}{{.Parent.Name}}) {{.Type}}) (result {{.Type}}) {
	for _, v := range rcv {
		result += fn(v)
	}
	return
}
`,
		RequiresNumeric: true,
	},
}
