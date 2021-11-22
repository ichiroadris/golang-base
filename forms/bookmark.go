package forms

type BookmarkPayload struct {
	Name string `json:"name" binding:"required"`
	Link string `json:"link" binding:"required"`
}
