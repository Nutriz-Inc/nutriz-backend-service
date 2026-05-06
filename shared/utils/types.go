package utils

type GenericReq struct {
	ActionBy string `validate:"required,id"`
}

func (g *GenericReq) SetActionBy(actionBy string) {
	g.ActionBy = actionBy
}

type PaginationReq struct {
	PageSize int `query:"page_size" default:"25" validate:"omitempty,min=1,max=50"`
	Page     int `query:"page" default:"1" validate:"omitempty,min=1"`
}

type PaginationInternalReq struct {
	PaginationReq
	GenericReq
}

type PaginationRes struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type GetReq struct {
	Id string `params:"id" validate:"required,id"`
	GenericReq
}

type DeleteRes struct {
	Success bool `json:"success"`
}
