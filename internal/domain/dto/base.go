package dto

const (
	RoleAdmin  = "admin"
	RoleStudent = "student"
	RoleMentor = "mentor"
)

type Response[T interface{}] struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
	Data    *T     `json:"data,omitempty"`
	Error   *T     `json:"error,omitempty"`
}

func NewResponse[T any](message string, status int, data *T, err *T) *Response[T] {
	return &Response[T]{
		Message: message,
		Status:  status,
		Data:    data,
		Error:   err,
	}
}


type PaginationResponse[T any] struct {
    Data       []T   `json:"paginate"`
    Total      int64 `json:"total"`
    Page       int   `json:"page"`
    PerPage    int   `json:"per_page"`
    TotalPages int   `json:"total_pages"`
    HasNext    bool  `json:"has_next"`
    HasPrev    bool  `json:"has_prev"`
}

func NewPaginationResponse[T any](data []T, total int64, page int, perPage int) *PaginationResponse[T] {
	totalPages := (total + int64(perPage) - 1) / int64(perPage)
	hasNext := page < int(totalPages)
	hasPrev := page > 1

	return &PaginationResponse[T]{
		Data:    data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: int(totalPages),
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}
}