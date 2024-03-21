package main

import (
    "testing"
)

func matchSlice(a, b []int) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}
func matchAnySlice(a []int, b [][]int) bool {
    for _, option := range b {
        if res := matchSlice(a, option); res {
            return true
        }
    }
    return false
}

func TestMoveElement(t *testing.T) {
    s := []int{1, 2, 3, 4, 5}
    s = MoveElement(s, 1, 3)
    expected := []int{1, 3, 4, 2, 5}
    if !matchSlice(s, expected) {
        t.Errorf("Expected %v, got %v, with indices %d, %d", expected, s, 1, 3)
    }

    s = MoveElement(s, 1, 1)
    expected = []int{1, 3, 4, 2, 5}
    if !matchSlice(s, expected) {
        t.Errorf("Expected %v, got %v, with indices %d, %d", expected, s, 1, 1)
    }
    s = MoveElement(s, 4, 0)
    expected = []int{5, 1, 3, 4, 2}
    if !matchSlice(s, expected) {
        t.Errorf("Expected %v, got %v, with indices %d, %d", expected, s, 4, 0)
    }

    s = MoveElement(s, 1, 0)
    expected = []int{1, 5, 3, 4, 2}

    if !matchSlice(s, expected) {
        t.Errorf("Expected %v, got %v, with indices %d, %d", expected, s, 1, 0)
    }
}
