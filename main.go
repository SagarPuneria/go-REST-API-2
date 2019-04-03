package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	sqlInterface "go-REST-API-2/sqlinterface"

	mx "github.com/gorilla/mux"
)

type app struct {
	sqlDB  *sqlInterface.MySqldb
	router *mx.Router
}

type album struct {
	AlbumId int    `json:"id"`
	UserId  int    `json:"userId"`
	Title   string `json:"title"`
}

type photo struct {
	AlbumId      int    `json:"albumId"`
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Url          string `json:"url"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

func main() {
	a := &app{}
	user := os.Args[1]
	pass := os.Args[2]
	host := os.Args[3]
	port := os.Args[4]
	DNS := user + ":" + pass + "@tcp(" + host + ":" + port + ")/"
	//DNS: user:passwd@tcp(127.0.0.1:3306)/
	a.sqlDB = a.initialiseDB(DNS)
	defer a.sqlDB.Close()

	go a.insertAlbumPhotoTable("https://jsonplaceholder.typicode.com/albums")

	a.router = mx.NewRouter()
	a.initialiseRouter()
	addr := "localhost:8080"
	a.run(addr)
}

func (a *app) insertAlbumPhotoTable(url string) {
	jsondata := getHttResponse(url)
	var albums []album
	err := json.Unmarshal(jsondata, &albums)
	if err != nil {
		log.Fatal(err)
	}
	err = a.sqlDB.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	for _, alb := range albums {
		strQuery := fmt.Sprintf("INSERT INTO album (id, userid, title) VALUES('%d', '%d', '%s');", alb.AlbumId, alb.UserId, alb.Title)
		err = a.sqlDB.ExecuteQuery(strQuery)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println("alb.AlbumId:", alb.AlbumId)
		url = "https://jsonplaceholder.typicode.com/photos?albumId=" + strconv.Itoa(alb.AlbumId)
		jsondata = getHttResponse(url)
		var photos []photo
		err = json.Unmarshal(jsondata, &photos)
		if err != nil {
			log.Fatal("Error in insert into album table, err:", err)
		}
		for _, pt := range photos {
			strQuery := fmt.Sprintf("INSERT INTO photo (id, albumId, photoId, title, url, thumbnailUrl) VALUES('%d', '%d', '%d', '%s', '%s', '%s');", pt.Id, pt.AlbumId, pt.Id, pt.Title, pt.Url, pt.ThumbnailUrl)
			err := a.sqlDB.ExecuteQuery(strQuery)
			if err != nil {
				log.Fatal("Error in insert into photo table, err:", err)
			}
		}
	}
	err = a.sqlDB.CommitTx()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Transaction commit done")
}

func getHttResponse(url string) []byte {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Connection", "close")
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func (a *app) initialiseDB(DNS string) *sqlInterface.MySqldb {
	createDB := "CREATE DATABASE IF NOT EXISTS typicode;"
	useDB := "USE typicode;"
	table1 := "CREATE TABLE IF NOT EXISTS album (id int NOT NULL, userId int, title varchar(255), PRIMARY KEY(id));"
	table2 := "CREATE TABLE IF NOT EXISTS photo (id int NOT NULL, albumId int, photoId int, title varchar(255), url varchar(255), thumbnailUrl varchar(255), FOREIGN KEY (albumId) REFERENCES album(id));"
	sqlDB, err := sqlInterface.CreateDataBase(DNS, createDB, useDB, table1, table2)
	if err != nil {
		log.Fatal(err)
	}
	return sqlDB
}

func (a *app) run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.router))
}

func (a *app) initialiseRouter() {
	//http://localhost:8080/search
	a.router.HandleFunc("/search", a.search).Methods("GET")
}

func (a *app) search(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside search")
	fmt.Println(r.FormValue("type"))
	fmt.Println(r.FormValue("id"))
	if r.FormValue("type") == "album" {
		strQuery := fmt.Sprintf("SELECT * FROM %s WHERE id=%s", r.FormValue("type"), r.FormValue("id"))
		fmt.Println("strQuery:", strQuery)
		rows, err := a.sqlDB.SelectQuery(strQuery)
		if err == nil {
			for rows.Next() {
				var title string
				var userId, id int
				if err := rows.Scan(&userId, &id, &title); err != nil {
					fmt.Println("rows.Scan ,Error Info : ", err)
					break
				}

				fmt.Println(userId, id, title)
				type paylod struct {
					UserId int
					Id     int
					Title  string
				}
				var p paylod
				p.Id = id
				p.UserId = userId
				p.Title = title
				if response, err := json.Marshal(p); err == nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write(response)
				}
			}
			rows.Close()
		}
	} else if r.FormValue("type") == "photo" {
		strQuery := fmt.Sprintf("SELECT id, albumId, title, url, thumbnailUrl FROM %s WHERE id=%s", r.FormValue("type"), r.FormValue("id"))
		fmt.Println("strQuery:", strQuery)
		rows, err := a.sqlDB.SelectQuery(strQuery)
		if err == nil {
			for rows.Next() {
				var title, url, thumbnailUrl string
				var albumId, id int
				if err := rows.Scan(&id, &albumId, &title, &url, &thumbnailUrl); err != nil {
					fmt.Println("rows.Scan ,Error Info : ", err)
					break
				}

				fmt.Println(id, albumId, title, url, thumbnailUrl)
				type paylod struct {
					AlbumId      int
					Id           int
					Title        string
					Url          string
					ThumbnailUrl string
				}
				var p paylod
				p.Id = id
				p.AlbumId = albumId
				p.Title = title
				p.Url = url
				p.ThumbnailUrl = thumbnailUrl
				if response, err := json.Marshal(p); err == nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write(response)
				}
			}
			rows.Close()
		}
	}
}
