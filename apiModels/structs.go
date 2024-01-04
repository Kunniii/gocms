package apiModels

type Like struct {
	UserID    string `json:"userID"`
	CreatedAt string `json:"createdAt"`
}

type Post struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"userID"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type Comment struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"userID"`
	Body   string `json:"body"`
}

type Tag struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
