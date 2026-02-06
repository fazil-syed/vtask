package schemas

// create a schema for creating a task
type CreateTaskInput struct {
	Title    string  `json:"title" binding:"required"`
	Content  string  `json:"content"`
	DueDate  *string `json:"due_date"`
	Timezone string  `json:"timezone"`
}
type EditTaskInput struct {
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	DueDate  *string `json:"due_date"`
	Timezone string  `json:"timezone"`
}

type TaskResponse struct {
	ID        uint    `json:"id"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	DueAt     *string `json:"due_at"`
	CreatedAt string  `json:"created_at"`
	Completed bool    `json:"completed"`
}
