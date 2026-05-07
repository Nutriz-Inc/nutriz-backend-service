package utils

type PaginationReq struct {
	PageSize int `query:"page_size" default:"25" validate:"required,min=1,max=50"`
	Page     int `query:"page" default:"1" validate:"required,min=1"`
}

type PaginationInternalReq struct {
	PaginationReq
}

type PaginationRes struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type GetReq struct {
	Id string `params:"id" validate:"required,id"`
}

type DeleteRes struct {
	Success bool `json:"success"`
}
