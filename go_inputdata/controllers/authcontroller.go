package controllers

import (
	"errors"
	"go_inputdata/config"
	"go_inputdata/entities"
	"go_inputdata/models"

	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

var usernameGlobal string

var userModel = models.NewUserModel()

func Index(w http.ResponseWriter, r *http.Request) {

	session, _ := config.Store.Get(r, config.SESSION_ID)

	if len(session.Values) == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {

		if session.Values["loggedIn"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {

			Home(w, r)

			// data := map[string]interface{}{
			// 	"username": session.Values["username"],
			// }

		}

	}
}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		temp, _ := template.ParseFiles("views/login.html")
		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		// proses login
		r.ParseForm()
		UserInput := &UserInput{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}

		// errorMessages := validation.Struct(UserInput)
		var errorMessages error

		if errorMessages != nil {

			data := map[string]interface{}{
				"validation": errorMessages,
			}

			temp, _ := template.ParseFiles("views/login.html")
			temp.Execute(w, data)

		} else {

			var user entities.User
			userModel.Where(&user, "username", UserInput.Username)

			var message error
			if user.Username == "" {
				message = errors.New("Incorrect username or password")
			} else {
				// pengecekan password
				errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(UserInput.Password))
				if errPassword != nil {
					message = errors.New("Incorrect username or password")
				}
			}

			if message != nil {

				data := map[string]interface{}{
					"error": message,
				}

				temp, _ := template.ParseFiles("views/login.html")
				temp.Execute(w, data)
			} else {
				usernameGlobal = user.Username
				// set session
				session, _ := config.Store.Get(r, config.SESSION_ID)

				session.Values["loggedIn"] = true
				session.Values["username"] = user.Username

				session.Save(r, w)

				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}

	}

}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, config.SESSION_ID)
	// delete session
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		temp, _ := template.ParseFiles("views/register.html")
		temp.Execute(w, nil)

	} else if r.Method == http.MethodPost {

		// melakukan proses registrasi

		// mengambil inputan form
		r.ParseForm()

		user := entities.User{

			Username:    r.Form.Get("username"),
			Password:    r.Form.Get("password"),
			Cpassword:   r.Form.Get("confirmpassword"),
		}

		erroMessages := make(map[string]interface{})

		if user.Username == ""{
			erroMessages["Username"] = "username cannot be empty"
		}
		if user.Password == ""{
			erroMessages["Password"] = "password cannot be empty"
		}
		if user.Cpassword == ""{
			erroMessages["Cpassword"] = "confirm password cannot be empty"
		} else {
			if user.Cpassword != user.Password{
			erroMessages["Cpassword"] = "confirm password does not match"
		}
	}

		if len(erroMessages) > 0 {
			// validasi gagal
			data := map[string]interface{}{
				"validation": erroMessages,
			}

			temp, _ := template.ParseFiles("views/register.html")
			temp.Execute(w, data)
		} else{
			// proses insert ke database
			//hashPassword
			hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			user.Password = string(hashPassword)

			//insert ke database
			_, err := userModel.Create(user)

			var message string
			if err != nil{
				message = "Register failed : username already use"
			} else {
				message = "Register success"
			}

			data := map[string]interface{}{
				"message": message,
			}

			temp, _ := template.ParseFiles("views/register.html")
			temp.Execute(w, data)
		}

}
}
