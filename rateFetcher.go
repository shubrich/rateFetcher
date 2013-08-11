package main

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"
import "database/sql"
import _ "github.com/mattn/go-sqlite3"
import "os"

func perror(err error) {
	if err != nil {
		panic(err)
	}
}

func get_rates() {
	url := "http://api.apirates.com/api/update"

	res, err := http.Get(url)
	perror(err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	perror(err)

	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		switch v := err.(type) {
		case *json.SyntaxError:
			fmt.Println(string(body[v.Offset-40 : v.Offset]))
		}
	}

	mappedData := data.(map[string]interface{})
	for k, v := range mappedData {
		fmt.Println("Iterating over ", k)
		mappedInnerData := v.(map[string]interface{})
		for l, w := range mappedInnerData {
			fmt.Println(l, " - ", w)
		}
	}
}

func main() {
	os.Remove("./rates.db")
	db, err := sql.Open("sqlite3", "./rates.db")
	perror(err)
	defer db.Close()

	sqls := []string{
		"create table rates (id integer not null primary key, fxpair text, fxrate real)",
		"delete from rates",
	}
	for _, sql := range sqls {
		_, err = db.Exec(sql)
		perror(err)
	}

	tx, err := db.Begin()
	perror(err)

	stmt, err := tx.Prepare("insert into rates(fxpair, fxrate) values(?, ?)")
	perror(err)

	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err := stmt.Exec(fmt.Sprintf("Test %03d", i), 1.234)
		perror(err)
	}
	tx.Commit()

	// get_rates()
}
