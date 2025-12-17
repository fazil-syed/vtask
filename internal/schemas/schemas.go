package schemas

// create a schema for creating a task
type CreateTaskInput struct {
	Name string `json:"name" binding:"required"`
	Completed bool   `json:"completed"`
}
// create a schema for updating a task
type UpdateTaskInput struct {
	Name string `json:"name"`
	Completed bool   `json:"completed"`
}
// create a schema for uploading a voice note
type UploadVoiceNoteInput struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	Size        int64  `json:"size" binding:"required"`
}