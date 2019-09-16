package fs

import (
	"sort"
)

// FileInfoComparer returns true, when f1 should be displayed before f2.
type FileInfoComparer func(f1, f2 FileInfo) int

var (
	// OrderDefault sorts elements lexicographically ascending and moves directories to the top of the list.
	OrderDefault = func(f1, f2 FileInfo) int {
		if order := OrderDirectoriesFirst(f1, f2); order != 0 {
			return order
		}
		return OrderLexicographicAsc(f1, f2)
	}

	// OrderFilesFirst moves files to the top of the list.
	OrderFilesFirst = func(f1, f2 FileInfo) int {
		if !f1.IsDir() && f2.IsDir() {
			return -1
		} else if f1.IsDir() && !f2.IsDir() {
			return 1
		}
		return 0
	}

	// OrderDirectoriesFirst moves directories to the top of the list.
	OrderDirectoriesFirst = func(f1, f2 FileInfo) int {
		return -OrderFilesFirst(f1, f2)
	}

	// OrderLexicographicAsc moves elements starting with A to the top of the list.
	OrderLexicographicAsc = func(f1, f2 FileInfo) int {
		if f1.Name() < f2.Name() {
			return -1
		} else if f1.Name() > f2.Name() {
			return 1
		}
		return 0
	}

	// OrderLexicographicDesc moves elements starting with Z to the top of the list.
	OrderLexicographicDesc = func(f1, f2 FileInfo) int {
		return -OrderLexicographicAsc(f1, f2)
	}
)

// NewCompoundComparer returns a new comparer based on the prioritized list of compare functions. The first comparer has the highest priority.
func NewCompoundComparer(compareFuncs ...FileInfoComparer) FileInfoComparer {
	return func(f1, f2 FileInfo) int {
		for _, f := range compareFuncs {
			order := f(f1, f2)
			if order != 0 {
				return order
			}
		}
		return 0
	}
}

// Sort sorts an array of FileInfo objects using the given comparer.
func Sort(files []FileInfo, cmp FileInfoComparer) {
	sort.Slice(files, func(i, j int) bool {
		return cmp(files[i], files[j]) < 0
	})
}
