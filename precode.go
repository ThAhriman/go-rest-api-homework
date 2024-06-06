package main

import (
	"encoding/json"
	//"fmt"
	"net/http"

	//"bytes"

	"log"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// ...

func getTasks(w http.ResponseWriter, r *http.Request) {

	resp, err := json.Marshal(tasks) // сбор данных в json формат
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print("Ошибка при сериализации данных в JSON формат")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(resp) // запись в тело ответа полученного ранее JSON
}

func addTask(w http.ResponseWriter, r *http.Request) {

	var task Task
	//var buf bytes.Buffer

	//	_, err := buf.ReadFrom(r.Body) //читаем запрос, добавляем в буфер
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusBadRequest)
	//		log.Print("Ошибка при чтении запроса")
	//		return
	//	}

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil { // читаем тело запроса и сразу десереализируем
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print("Ошибка при попытке десериализации")
		return
	}

	for id, _ := range tasks {
		if id == task.ID {
			http.Error(w, "Задача с указанным id уже существует", http.StatusBadRequest)
			log.Println("Задача с указанным id уже существует")
			return
		}
	}
	if task.ID == "" {
		task.ID = "-"
	}
	if task.Applications == nil {
		task.Applications = make([]string, 0)
		task.Applications = append(task.Applications, r.UserAgent())
	}
	tasks[task.ID] = task //вносим в мапу данные по ключу равному значению поля ID структуры task (добавляем ключ и его значение)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func getTask(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id") // возврат параметра из запроса (Request)

	task, ok := tasks[id] //инициализация (+ проверка наличия в мапе) структуры по ключу(параметр полученный с помощью URLParam)
	if !ok {
		http.Error(w, "Задача с указанным id не обнаружена", http.StatusBadRequest)
		log.Printf("Задача с id = %s не обнаружена", id)
		return
	}

	resp, err := json.Marshal(task) //сериализация запрашиваемой структуры
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print("Ошибка при сериализации данных в JSON формат")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id") // возврат параметра из запроса (Request)

	_, ok := tasks[id] //инициализация (+ проверка наличия в мапе) структуры по ключу(параметр полученный с помощью URLParam)
	if !ok {
		http.Error(w, "Задача с указанным id не обнаружена", http.StatusBadRequest)
		log.Printf("Задача с id = %s не обнаружена", id)
		return
	}

	delete(tasks, id) // удаление пары ключ - значение из мапы по id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// ...

	r.Get("/tasks", getTasks) // регистрируем эндпоинт "/tasks" с методом Get

	r.Post("/tasks", addTask) // регистрируем эндпоинт "/tasks" с методом Post

	r.Get("/tasks/{id}", getTask) // регистрируем эндпоинт "/tasks{id}" с методом Get

	r.Delete("/tasks/{id}", deleteTask) //регистрируем эндпоинт "/tasks{id}" с методом Delete

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
