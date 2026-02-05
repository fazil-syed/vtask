package schemas

// create a schema for creating a task
type CreateTaskInput struct {
	Title    string  `json:"title" binding:"required"`
	Content  string  `json:"content"`
	DueDate  *string `json:"due_date"`
	Timezone string  `json:"timezone"`
}

// create a schema for updating a task
type UpdateTaskInput struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

type TaskResponse struct {
	ID        uint    `json:"id"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	DueAt     *string `json:"due_at"`
	CreatedAt string  `json:"created_at"`
}
