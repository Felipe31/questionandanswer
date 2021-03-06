package routes

import (
	"encoding/json"
	"felipesoares/questionandanswer/internal/model"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

func HandleNewUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleNewUser")
	var newUser model.User

	reqBody, err := ioutil.ReadAll(r.Body)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error while processing user body!", err.Error())
		return
	}
	json.Unmarshal(reqBody, &newUser)

	_, err = model.UserByUsername(newUser.Username)
	if err == nil {
		WriteError(w, http.StatusConflict, "Username already registered", "")
		return
	}

	newUser.ID, err = model.AddUser(newUser)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error while creating user! ", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// Note: Deleting a user does not delete its questions
func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleDeleteUser")
	var user model.User
	userIdVar := mux.Vars(r)["id-username"]
	w.Header().Set("Content-Type", "application/json")
	userId, err := strconv.Atoi(userIdVar)
	if err != nil {
		// Check if the username is in the URL instead of the ID
		user, err = model.UserByUsername(userIdVar)
		if err != nil {
			WriteError(w, http.StatusNotFound, "User to be deleted not found!", err.Error())
			return
		}
		userId = int(user.ID)
	} else {
		user, err = model.UserByID(userId)
		if err != nil {
			WriteError(w, http.StatusNotFound, "User to be deleted not found!", err.Error())
			return
		}
	}
	err = model.RemoveUser(userId)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error while removing user!", err.Error())
		return
	}

	json.NewEncoder(w).Encode(user)

}

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleGetUser")
	var user model.User
	userIdVar := mux.Vars(r)["id-username"]
	w.Header().Set("Content-Type", "application/json")
	userId, err := strconv.Atoi(userIdVar)
	if err != nil {
		// Check if the username is in the URL instead of the ID
		user, err = model.UserByUsername(userIdVar)
		if err != nil {
			WriteError(w, http.StatusNotFound, "User not found!", err.Error())
			return
		}
	} else {
		user, err = model.UserByID(userId)
		if err != nil {
			WriteError(w, http.StatusNotFound, "User not found!", err.Error())
			return
		}
	}

	json.NewEncoder(w).Encode(user)

}

func HandleGetUserQuestions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleGetUserQuestions")

	userIdVar := mux.Vars(r)["id-username"]
	w.Header().Set("Content-Type", "application/json")
	userId, err := strconv.Atoi(userIdVar)
	if err != nil {
		// Check if the username is in the URL instead of the ID
		user, err := model.UserByUsername(userIdVar)
		if err != nil {
			WriteError(w, http.StatusBadGateway, "Invalid username!", err.Error())
			return
		}
		userId = int(user.ID)
	}
	_, err = model.UserByID(userId)
	if err != nil {
		WriteError(w, http.StatusBadGateway, "Invalid id!", err.Error())
		return
	}

	questions, err := model.QuestionsByUserId(userId)
	if err != nil {
		WriteError(w, http.StatusNotFound, "Error while retrieving questions of a user!", err.Error())
		return
	}

	if questions == nil {
		json.NewEncoder(w).Encode([]model.Question{})
		return
	}

	json.NewEncoder(w).Encode(questions)
}

func HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleGetAllUsers")
	w.Header().Set("Content-Type", "application/json")
	users, err := model.GetAllUsers()
	if err != nil {
		WriteError(w, http.StatusNotFound, "Error while retrieving users!", err.Error())
		return
	}
	if users == nil {
		json.NewEncoder(w).Encode([]model.User{})
		return
	}

	json.NewEncoder(w).Encode(users)

}
