package controllers

import (
	"bytes"
	"encoding/json"
	"go_inputdata/entities"
	"go_inputdata/models"
	"html/template"
	"net/http"
	"strconv"
)

var inputModel = models.New()

func Home(w http.ResponseWriter, r *http.Request) {

		data := map[string]interface{}{
			"data": template.HTML(GetData()),
			"username" : usernameGlobal,
		}

		temp, _ := template.ParseFiles("views/index.html")
		temp.Execute(w, data)
}

func GetData()string{

	buffer := &bytes.Buffer{}

	temp, _ := template.New("data.html").Funcs(template.FuncMap{
		"increment": func (a, b int) int {
			return a + b
		},
	}).ParseFiles("views/data.html")

	var inputdata []entities.Inputdata

	err := inputModel.FindAll(&inputdata)

	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"inputdata": inputdata,
	}

	temp.ExecuteTemplate(buffer, "data.html", data)

	return buffer.String()
}

func GetForm(w http.ResponseWriter, r *http.Request){

	queryString := r.URL.Query()
	id, err := strconv.ParseInt(queryString.Get("id"), 10, 64)
	
	var data map[string]interface{}
	var inputdata entities.Inputdata

	if err != nil{
		data = map[string]interface{}{
			"title":"Add Data",
		}
	} else{

		err := inputModel.Find(id, &inputdata)
		if err != nil {
			panic(err)
		}

		data = map[string]interface{}{
			"title":"Edit Data",
			"inputdata": inputdata,
		}
	}


	temp, _ := template.ParseFiles("views/form.html")
	temp.Execute(w, data)
}

func Store(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost{
		r.ParseForm()

		var inputdata entities.Inputdata
		inputdata.Name_person = r.Form.Get("name_person")
		inputdata.Npm = r.Form.Get("npm")
		inputdata.Gender = r.Form.Get("gender")
		inputdata.Birth_date = r.Form.Get("birth_date")
		inputdata.Address = r.Form.Get("address")

		id, err := strconv.ParseInt(r.Form.Get("id"), 10, 64)

		var data map[string]interface{}
		if err!= nil{
			// insert data
			err := inputModel.Create(&inputdata)

		if err != nil{
			ResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}

			data = map[string]interface{}{
			"message":"Adding data success",
			"data": template.HTML(GetData()),
		}
		} else{
			// update data
			inputdata.Id = id
			err := inputModel.Update(inputdata)
			if err != nil{
				ResponseError(w, http.StatusInternalServerError, err.Error())
				return
			}
			data = map[string]interface{}{
				"message":"Data update success",
				"data": template.HTML(GetData()),
			}

		}
		

		

		ResponseJson(w, http.StatusOK, data)
	}
}

func Delete(w http.ResponseWriter, r *http.Request){

	r.ParseForm()

	id, err := strconv.ParseInt(r.Form.Get("id"), 10, 64)

	if err != nil{
		panic(err)
	}

	err = inputModel.Delete(id)

	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"message": "Data delete success",
		"data":template.HTML(GetData()),
	}

	ResponseJson(w, http.StatusOK, data)

}


func ResponseError(w http.ResponseWriter, code int, message string){
	ResponseJson(w, code, map[string]string{"error":message})
}

// mengirim respon dalam bentuk json
func ResponseJson(w http.ResponseWriter, code int, payload interface{}){
	 response, _ := json.Marshal(payload)

	 w.Header().Set("Content-Type", "application/json")
	 w.WriteHeader(code)
	 w.Write(response)
}