package algorithms

import (
	"reflect"
	"testing"
)

func TestMergeList(t *testing.T) {
	tests := []struct {
		name   string
		list1  []int32
		list2  []int32
		expect []int32
	}{
		{
			name:   "Both empty",
			list1:  []int32{},
			list2:  []int32{},
			expect: []int32{},
		},
		{
			name:   "One empty",
			list1:  []int32{1, 3, 5},
			list2:  []int32{},
			expect: []int32{1, 3, 5},
		},
		{
			name:   "No duplicates",
			list1:  []int32{1, 2, 4},
			list2:  []int32{3, 5, 6},
			expect: []int32{1, 2, 3, 4, 5, 6},
		},
		{
			name:   "With duplicates between lists",
			list1:  []int32{1, 2, 3, 5},
			list2:  []int32{2, 4, 5, 6},
			expect: []int32{1, 2, 3, 4, 5, 6},
		},
		{
			name:   "With internal duplicates",
			list1:  []int32{1, 2, 2, 3},
			list2:  []int32{2, 3, 4, 4},
			expect: []int32{1, 2, 3, 4},
		},
		{
			name:   "All values equal",
			list1:  []int32{1, 1, 1},
			list2:  []int32{1, 1, 1},
			expect: []int32{1},
		},
		{
			name:   "Interleaved duplicates",
			list1:  []int32{1, 3, 5, 7},
			list2:  []int32{2, 3, 6, 7},
			expect: []int32{1, 2, 3, 5, 6, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := MergeList(tt.list1, tt.list2)
			if !reflect.DeepEqual(actual, tt.expect) {
				t.Errorf("MergeList(%v, %v) = %v; want %v", tt.list1, tt.list2, actual, tt.expect)
			}
		})
	}
}

func TestMissingElementsInSortedList(t *testing.T) {
	tests := []struct {
		name     string
		total    []int32
		subList  []int32
		expected []int32
	}{
		{
			name:     "No missing elements",
			total:    []int32{1, 2, 3},
			subList:  []int32{1, 2, 3},
			expected: []int32{},
		},
		{
			name:     "One missing element in middle",
			total:    []int32{1, 2, 3},
			subList:  []int32{1, 3},
			expected: []int32{2},
		},
		{
			name:     "Missing elements at end",
			total:    []int32{1, 2, 3, 4, 5},
			subList:  []int32{1, 2, 3},
			expected: []int32{4, 5},
		},
		{
			name:     "Missing elements at beginning",
			total:    []int32{1, 2, 3, 4, 5},
			subList:  []int32{3, 4, 5},
			expected: []int32{1, 2},
		},
		{
			name:     "Empty sublist",
			total:    []int32{1, 2, 3},
			subList:  []int32{},
			expected: []int32{1, 2, 3},
		},
		{
			name:     "Empty total list",
			total:    []int32{},
			subList:  []int32{},
			expected: []int32{},
		},
		{
			name:     "Sublist has element not in total",
			total:    []int32{1, 3, 5},
			subList:  []int32{1, 2},
			expected: []int32{}, // indicates invalid subList
		},
		{
			name:     "Sublist longer than total",
			total:    []int32{1, 2},
			subList:  []int32{1, 2, 3},
			expected: []int32{}, // invalid input
		},
		{
			name:     "Sublist is nil",
			total:    []int32{1, 2, 3},
			subList:  nil,
			expected: []int32{1, 2, 3},
		},
		{
			name:     "Total is nil",
			total:    nil,
			subList:  []int32{1, 2},
			expected: []int32{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MissingElementsInSortedList(tt.total, tt.subList)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MissingElementsInSortedList(%v, %v) = %v; want %v", tt.total, tt.subList, result, tt.expected)
			}
		})
	}
}
