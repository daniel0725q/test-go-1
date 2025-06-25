package dto

// PaginatedResponse represents a paginated response with metadata
type PaginatedResponse struct {
	Data       []*StockRatingResponse `json:"data"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalCount int64                  `json:"total_count"`
	TotalPages int                    `json:"total_pages"`
	HasNext    bool                   `json:"has_next"`
	HasPrev    bool                   `json:"has_prev"`
}
