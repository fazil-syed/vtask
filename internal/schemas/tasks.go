package schemas

type ReminderInput struct {
	OffsetMinutes *int `json:"offset_minutes,omitempty"`
	// RemindAt      *string `json:"remind_at,omitempty"`
}

// create a schema for creating a task
type CreateTaskInput struct {
	Title     string          `json:"title" binding:"required"`
	Content   string          `json:"content"`
	DueDate   *string         `json:"due_date"`
	Timezone  string          `json:"timezone"`
	Reminders []ReminderInput `json:"reminders"`
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
