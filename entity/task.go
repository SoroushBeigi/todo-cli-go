package entity

type Task struct {
	ID         int
	UserID     int
	Title      string
	DueDate    string
	CategoryID int
	IsDone     bool
}
