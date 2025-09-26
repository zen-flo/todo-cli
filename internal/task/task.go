package task

import "time"

// Task — основная модель задачи.
type Task struct {
	ID        int       `json:"id"`         // Уникальный идентификатор
	Title     string    `json:"title"`      // Заголовок задачи
	Completed bool      `json:"completed"`  // Статус выполнения (true = выполнено)
	CreatedAt time.Time `json:"created_at"` // Время создания задачи
}

// MarkDone — метод, который отмечает задачу как выполненную.
func (t *Task) MarkDone() {
	t.Completed = true
}
