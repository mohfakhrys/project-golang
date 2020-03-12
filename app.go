package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "bri"
	dbname   = "app-user"
)

//struct untuk db
type Info struct {
	Username   string
	Department string
	Created    string
}

//struct untuk json
type User struct {
	Username   string `json:username`
	Department string `json:department`
}

func getInfo(Username string) (err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("error connect db", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println("error ping db", err)
		return
	}

	var i Info

	err = db.QueryRow("select username, department, created from userinfo where username = 'fakhry'").Scan(&i.Username, &i.Department, &i.Created)
	if err != nil {
		log.Println("error query get", err)
		return
	}
	fmt.Printf("username: %s, department: %s, create:%s\n", i.Username, i.Department, i.Created)

	//untuk ngerubah ke json
	user := &User{Username: i.Username, Department: i.Department}
	u, err := json.Marshal(user)
	if err != nil {
		fmt.Println("error ke json", err)
		return
	}
	fmt.Println(string(u))

	mc, err := json.Marshal(i)
	fmt.Println(i)
	if err != nil {
		log.Println(err)
		return
	}

	err = ioutil.WriteFile("./file/var/json/home.json", mc, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	return err
}

func JsonHome(w http.ResponseWriter, r *http.Request) {
	jsonFile, err := os.Open("./file/var/json/home.json")
	if err != nil {
		log.Println("error open json", err)
		return
	}
	fmt.Println("success open file json")
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println("[cuaca] error open file json", err)
		return
	}

	var data User

	err = json.Unmarshal([]byte(byteValue), &data)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(data)

	t, err := template.ParseFiles("./file/var/json/home.json")

	err = t.Execute(w, data)
}

func insetData() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("error connect db", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println("error ping db", err)
		return
	}
	fmt.Println("sukses")

	smt, err := db.Prepare("insert into userinfo (username, department, created) values ($1,$2,$3)")
	if err != nil {
		log.Println("error prepare", err)
		return
	}

	resp, err := smt.Exec("fakhry", "administrator", "2020-03-08 15:15:00")
	if err != nil {
		log.Println("error exec", err)
		return
	}

	_, ra := resp.RowsAffected()
	log.Println(ra)

}

func main() {
	getInfo("fakhry")

	http.HandleFunc("/Home", JsonHome)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Println("main func error listening")
	}
}
