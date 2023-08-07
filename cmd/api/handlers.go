package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mikekorakakis/g_test/internal/models"
)

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

func (app *application) TestHandler(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("In Test Handler")
	j := jsonResponse{
		OK:      true,
		Message: "OK",
		Content: "",
	}
	out, err := json.MarshalIndent(j, "", "   ")
	if err != nil {
		app.errorLog.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *application) GetToDos(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("hello")
	todos, err := app.DB.GetToDos()
	app.infoLog.Println("hello2")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJson(w, http.StatusOK, todos)

}

func (app *application) GetToDoByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	todoID, _ := strconv.Atoi(id)

	todo, err := app.DB.GetToDo(todoID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	app.writeJson(w, http.StatusOK, todo)
}

func (app *application) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.ToDo

	err := app.readJSON(w, r, &todo)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	todo, err = app.DB.CreateToDo(todo)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJson(w, http.StatusOK, todo)

}

func (app *application) EditToDo(w http.ResponseWriter, r *http.Request) {
	var todo models.ToDo

	err := app.readJSON(w, r, &todo)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	todo, err = app.DB.EditToDo(todo)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJson(w, http.StatusOK, todo)
}

func (app *application) DeleteToDo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	todoId, err := strconv.Atoi(id)

	if err != nil {
		app.errorLog.Println(err)
	}

	todo, err := app.DB.DeleteToDo(todoId)

	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJson(w, http.StatusOK, todo)
}
