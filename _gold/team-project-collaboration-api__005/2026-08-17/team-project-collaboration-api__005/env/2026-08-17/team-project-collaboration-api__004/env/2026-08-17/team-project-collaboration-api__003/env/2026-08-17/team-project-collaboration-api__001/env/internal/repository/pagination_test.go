package repository

import "testing"

func TestNewPaginationDefaults(t *testing.T) {
	got := NewPagination(0, 0)
	if got.Page != 1 || got.PageSize != 20 || got.Offset != 0 {
		t.Fatalf("NewPagination(0, 0) = %+v", got)
	}
}

func TestNewPaginationCapsPageSize(t *testing.T) {
	got := NewPagination(3, 500)
	if got.Page != 3 || got.PageSize != maxPageSize || got.Offset != 2*maxPageSize {
		t.Fatalf("NewPagination(3, 500) = %+v", got)
	}
}

func TestTotalPages(t *testing.T) {
	if got := TotalPages(0, 20); got != 0 {
		t.Fatalf("TotalPages(0) = %d", got)
	}
	if got := TotalPages(21, 20); got != 2 {
		t.Fatalf("TotalPages(21) = %d", got)
	}
}
