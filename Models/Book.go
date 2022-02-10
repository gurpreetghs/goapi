package Models

type Book struct {
	ID     uint64   `json:"id" gorm:"primary_key"`
	Title  string `json:"title" `
	Author string `json:"author"`
	Username string `json:"username"`
 	Password string `json:"password"`
}
type API_SECRET struct {
	Username string `json:"username"`
 	Password string `json:"password"`
}