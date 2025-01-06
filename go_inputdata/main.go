package main

import (
	"fmt"
	"net/http"

	"go_inputdata/controllers"
	authcontroller "go_inputdata/controllers"
)

func main() {
	http.HandleFunc("/", authcontroller.Index)
	http.HandleFunc("/login", authcontroller.Login)
	http.HandleFunc("/logout", authcontroller.Logout)
	http.HandleFunc("/register", authcontroller.Register)
	http.HandleFunc("/save", authcontroller.Register)
	http.HandleFunc("/add_data/get_form", controllers.GetForm)
	http.HandleFunc("/add_data/store", controllers.Store)
	http.HandleFunc("/data/delete", controllers.Delete)
	

	fmt.Println("Server jalan di: http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
