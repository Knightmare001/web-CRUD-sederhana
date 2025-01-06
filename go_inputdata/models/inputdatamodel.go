package models

import (
	"database/sql"
	"go_inputdata/config"
	"go_inputdata/entities"
)

type InputModel struct {
	db *sql.DB
}


func New() *InputModel{
	db, err := config.DBConn()

	if err != nil{
		panic(err)
	}
	return &InputModel{db: db}
}


func (m *InputModel) FindAll(inputdata *[]entities.Inputdata) error{
	rows, err := m.db.Query("SELECT id, name_person, npm,gender, birth_date, address FROM biodata")

	if err != nil{
		return err
	}
	defer rows.Close() 

	for rows.Next(){
		var data entities.Inputdata
		rows.Scan(&data.Id, &data.Name_person, &data.Npm, &data.Gender, &data.Birth_date, &data.Address)


		*inputdata = append(*inputdata, data)
	}

	return nil
}

func (m *InputModel) Create(inputdata *entities.Inputdata) error {

	result, err := m.db.Exec("INSERT INTO biodata (name_person, npm, gender, birth_date, address) VALUES(?,?,?,?,?)",
	inputdata.Name_person,inputdata.Npm,inputdata.Gender,inputdata.Birth_date,inputdata.Address)


	if err != nil {
		return err
	}

	lastInsertId, _ := result.LastInsertId()
	inputdata.Id = int64(lastInsertId)
	return nil

}



func (m *InputModel) Find(id int64, inputdata *entities.Inputdata) error{
	return m.db.QueryRow("SELECT id, name_person, npm, gender, birth_date, address FROM biodata WHERE id = ?", id).Scan(
		&inputdata.Id,
		&inputdata.Name_person,
		&inputdata.Npm,
		&inputdata.Gender,
		&inputdata.Birth_date,
		&inputdata.Address,
	)

}

func (m *InputModel) Update(inputdata entities.Inputdata) error{

	_,err := m.db.Exec("UPDATE biodata SET name_person = ?, npm = ?, gender = ?, birth_date = ?, address = ? WHERE id = ?",
							inputdata.Name_person, inputdata.Npm, inputdata.Gender, inputdata.Birth_date, inputdata.Address, inputdata.Id)


	if err != nil{
		return err
	}

	return nil
}


func (m *InputModel) Delete(id int64) error{
	_, err := m.db.Exec("DELETE FROM biodata WHERE id = ?", id)

	if err != nil {
		return err
	}

	return nil
}