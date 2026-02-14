package request

type CreateTagRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=20"`
	Color string `json:"color" binding:"required,hexcolor"`
}
type UpdateTagRequest struct {
	ID    uint
	Name  *string `json:"name" binding:"omitempty,min=1,max=20"`
	Color *string `json:"color" binding:"omitempty,hexcolor"`
}

type GetByIdTagRequest struct {
	ID uint
}
type DeleteTagRequest struct {
	ID uint
}
