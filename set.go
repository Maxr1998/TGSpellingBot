package main

import (
	"sort"
)

// AddToSortedStringSet add element to the set (which must be sorted) at the right position
func AddToSortedStringSet(set []string, element string) (bool, []string) {
	i := sort.SearchStrings(set, element)
	if i < len(set) && set[i] == element {
		return false, set
	}
	set = append(set, "_")
	if i < len(set)-1 {
		copy(set[i+1:], set[i:])
	}
	set[i] = element
	return true, set
}

// RemoveFromSortedStringSet remove element from the set (which must be sorted) if possible
func RemoveFromSortedStringSet(set []string, element string) (bool, []string) {
	i := sort.SearchStrings(set, element)
	if i < len(set) && set[i] == element {
		copy(set[i:], set[i+1:])
		return true, set[:len(set)-1]
	}
	return false, set
}
