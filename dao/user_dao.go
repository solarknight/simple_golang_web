package dao

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	id    int
	name  string
	sex   int
	descp string
}

func OpenDB() {
	db, err := sql.Open("mysql", "root:solarknight@tcp(127.0.0.1:3306)/user_info")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// user, _ := GetByID(db, 1)
	user, _ := GetByIDSimpler(db, 3)
	log.Println(user)
}

func GetByIDSimpler(db *sql.DB, id int) (*User, error) {
	user := new(User)
	err := db.QueryRow("select id, name, sex, descp from user where id = ?", id).Scan(&user.id, &user.name, &user.sex, &user.descp)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return user, nil
}

func GetByID(db *sql.DB, id int) (*User, error) {
	stmt, err := db.Prepare("select id, name, sex, descp from user where id = ?")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(1)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	user := new(User)
	for rows.Next() {
		err := rows.Scan(&user.id, &user.name, &user.sex, &user.descp)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		log.Println(user)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return user, nil
}
