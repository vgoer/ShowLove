// Package pagination provides pagination utilities for list endpoints.
package pagination

const (
	DefaultPageSize int32 = 20
	MaxPageSize     int32 = 100
)

// Page holds pagination parameters.
type Page struct {
	Page     int32
	PageSize int32
}

// Parse creates a Page with sanitized values.
// page is 1-based; pageSize is clamped between 1 and MaxPageSize.
func Parse(page, pageSize int32) Page {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return Page{Page: page, PageSize: pageSize}
}

// Offset returns the SQL OFFSET value for the current page.
func (p Page) Offset() int32 {
	return (p.Page - 1) * p.PageSize
}

// TotalPages calculates the total number of pages given the total item count.
func (p Page) TotalPages(total int64) int32 {
	if total == 0 {
		return 0
	}
	pages := int32(total / int64(p.PageSize))
	if total%int64(p.PageSize) != 0 {
		pages++
	}
	return pages
}
