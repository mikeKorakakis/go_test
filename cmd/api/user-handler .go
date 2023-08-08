package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mikekorakakis/g_test/internal/models"
)

func (app *application) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.DB.GetAllUsers()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJson(w, http.StatusOK, users)
}

func (app *application) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(id)

	user, err := app.DB.GetOneUser(userID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	app.writeJson(w, http.StatusOK, user)
}

func (app *application) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := app.readJSON(w, r, &user)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, err = app.DB.AddUser(user, "password")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJson(w, http.StatusOK, user)

}

func (app *application) EditUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := app.readJSON(w, r, &user)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, err = app.DB.EditUser(user)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJson(w, http.StatusOK, user)
}

func (app *application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userId, err := strconv.Atoi(id)

	if err != nil {
		app.errorLog.Println(err)
	}

	user, err := app.DB.DeleteUser(userId)

	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJson(w, http.StatusOK, user)
}
