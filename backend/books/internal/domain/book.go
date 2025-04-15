package domain

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Price  int    `json:"price"`
}

type DeleteRequest struct {
	ID int `json:"id"`
}

type ListBooks []Book
