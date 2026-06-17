package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name         string
		page         int32
		pageSize     int32
		wantPage     int32
		wantPageSize int32
		wantOffset   int32
	}{
		{"default values", 0, 0, 1, 20, 0},
		{"normal values", 2, 10, 2, 10, 10},
		{"negative page", -1, 10, 1, 10, 0},
		{"zero page size defaults", 1, 0, 1, 20, 0},
		{"exceeds max page size", 1, 200, 1, 100, 0},
		{"large page", 100, 20, 100, 20, 1980},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Parse(tt.page, tt.pageSize)
			assert.Equal(t, tt.wantPage, p.Page)
			assert.Equal(t, tt.wantPageSize, p.PageSize)
			assert.Equal(t, tt.wantOffset, p.Offset())
		})
	}
}

func TestPage_Offset(t *testing.T) {
	p := Page{Page: 3, PageSize: 25}
	assert.Equal(t, int32(50), p.Offset())
}

func TestPage_TotalPages(t *testing.T) {
	tests := []struct {
		name       string
		total      int64
		pageSize   int32
		totalPages int32
	}{
		{"exact", 100, 20, 5},
		{"remainder", 101, 20, 6},
		{"zero items", 0, 20, 0},
		{"one item", 1, 20, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Page{PageSize: tt.pageSize}
			assert.Equal(t, tt.totalPages, p.TotalPages(tt.total))
		})
	}
}

func TestMaxPageSize(t *testing.T) {
	assert.Equal(t, int32(100), MaxPageSize)
}

func TestDefaultPageSize(t *testing.T) {
	assert.Equal(t, int32(20), DefaultPageSize)
}
