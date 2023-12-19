package itypes

type Role struct {
	ID     uint
	Name   string
	Detail string
}

// Keep it in this order
var Roles = []Role{
	{ID: 0, Name: "user", Detail: "Normal User, which can comment"},
	{ID: 1, Name: "contributor", Detail: "The one who can create post!"},
	{ID: 2, Name: "publisher", Detail: "The one who can create and publish posts"},
	{ID: 3, Name: "manager", Detail: "The one who can manage posts, tags...."},
	{ID: 4, Name: "admin", Detail: "The one who can do pretty much anything!"},
}
