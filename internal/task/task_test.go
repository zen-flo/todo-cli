package task

import (
	"encoding/json"
	"testing"
	"time"
)

// TestMarkDone проверяет, что метод MarkDone корректно отмечает задачу как выполненную.
func TestMarkDone(t *testing.T) {
	t1 := Task{
		ID:        1,
		Title:     "Test",
		Completed: false,
		CreatedAt: time.Now(),
		Important: true,
	}

	t1.MarkDone()

	if !t1.Completed {
		t.Errorf("ожидалось, что Completed=true, но получили false")
	}
}

// TestTaskJSON проверяет корректность сериализации и десериализации задачи.
func TestTaskJSON(t *testing.T) {
	now := time.Now().Truncate(time.Second) // округляем, чтобы не потерять точность при сравнении
	original := Task{
		ID:        42,
		Title:     "JSON Test",
		Completed: true,
		CreatedAt: now,
		Important: false,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("ошибка при маршалинге: %v", err)
	}

	var decoded Task
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("ошибка при анмаршалинге: %v", err)
	}

	if decoded.ID != original.ID ||
		decoded.Title != original.Title ||
		decoded.Completed != original.Completed ||
		!decoded.CreatedAt.Equal(original.CreatedAt) ||
		decoded.Important != original.Important {
		t.Errorf("данные после JSON-сериализации не совпадают:\nисходный: %+v\nполученный: %+v",
			original, decoded)
	}
}

// TestTaskDefaults проверяет значения по умолчанию для новой задачи.
func TestTaskDefaults(t *testing.T) {
	t1 := Task{Title: "Default test"}

	if t1.ID != 0 {
		t.Errorf("ожидался ID=0, получили %d", t1.ID)
	}
	if t1.Completed {
		t.Errorf("ожидалось Completed=false по умолчанию")
	}
	if t1.Important {
		t.Errorf("ожидалось Important=false по умолчанию")
	}
	if t1.Title != "Default test" {
		t.Errorf("ожидался Title=Default test, получили %q", t1.Title)
	}
}
