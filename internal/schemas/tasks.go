package schemas

// create a schema for creating a task
type CreateTaskInput struct {
	Title    string  `json:"title" binding:"required"`
	Content  string  `json:"content" binding:"required"`
	DueDate  *string `json:"due_date"`
	Timezone string  `json:"timezone"`
}

// create a schema for updating a task
type UpdateTaskInput struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}
