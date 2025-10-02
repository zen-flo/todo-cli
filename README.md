# todo-cli

Простое командное приложение TODO на Go.  
Управляйте задачами, отмечайте их выполненными, фильтруйте, ищите и многое другое прямо из терминала.

---

## Возможности

- Добавление задач (`todo add "Название задачи"`)
- Отметка задач как выполненных (`todo done <ID>`)
- Удаление задач (`todo delete <ID>`)
- Обновление заголовка задачи и отметка как важной (`todo update <ID> "Новый заголовок" --important`)
- Просмотр всех задач (`todo list`) с возможностью сортировки (`--sort=name|date`) и фильтрации (`--filter=all|pending|completed`)
- Просмотр только невыполненных задач (`todo pending`)
- Просмотр только выполненных задач (`todo completed`)
- Отметить все задачи как выполненные (`todo complete-all`)
- Очистка выполненных задач (`todo clear`)
- Поиск задач по ключевому слову (`todo search "ключевое слово"`)

---

## Установка

1. Клонируйте репозиторий:

```bash
git clone https://github.com/zen-flo/todo-cli.git
cd todo-cli
```
2. Соберите CLI:

```bash
go build -o todo
```

3. Запустите CLI:

```bash
./todo add "Купить молоко"
./todo list
```

---

## Примеры использования

### Добавить задачу
```bash
todo add "Закончить проект на Go" --important
```

### Просмотр задач
```bash
todo list --filter=pending --sort=date
```

### Отметить задачу как выполненную
```bash
todo done 1
```

### Обновить заголовок задачи
```bash
todo update 1 "Закончить Go CLI проект" --important
```

### Удалить задачу
```bash
todo delete 1
```

### Поиск задач
```bash
todo search "Go"
```

### Очистить выполненные задачи
```bash
todo clear
```

### Отметить все задачи как выполненные
```bash
todo complete-all
```

---

## Тестирование

Проект содержит полный набор тестов.

### Запуск всех тестов с покрытием:

```bash
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

#### Текущее покрытие:

* `cmd`: ~70%

* `storage`: ~86%

* `task`: 100%

---