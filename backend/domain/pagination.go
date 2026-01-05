package domain

// PaginationInput represents pagination parameters from query string
type PaginationInput struct {
	Page    int `query:"page"`
	PerPage int `query:"perPage"`
}

// Validate validates and sets default values for pagination
func (p *PaginationInput) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 {
		p.PerPage = 10
	}
	// Limit maximum per page to prevent abuse
	if p.PerPage > 100 {
		p.PerPage = 100
	}
}

// PaginationResponse represents pagination metadata in response
type PaginationResponse struct {
	Total   int64 `json:"total"`
	Page    int   `json:"page"`
	PerPage int   `json:"perPage"`
	Pages   int   `json:"pages"`
}

// CalculatePages calculates total number of pages
func (p *PaginationResponse) CalculatePages() {
	if p.PerPage > 0 {
		p.Pages = int((p.Total + int64(p.PerPage) - 1) / int64(p.PerPage))
	}
}

// PaginatedResult wraps data with pagination metadata
type PaginatedResult struct {
	Data       interface{}        `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}
