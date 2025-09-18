package task

import "time"

// Task представляет задачу
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// MarkDone отмечает задачу как выполненную
func (t *Task) MarkDone() {
	t.Completed = true
}
