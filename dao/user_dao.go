package dao

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type User struct {
	id    int
	name  string
	sex   int
	descp string
}

type Person struct {
	ID    int
	Name  string
	Sex   int
	Descp string
}

func ConnectDB() {
	db, err := sqlx.Connect("mysql", "root:solarknight@tcp(127.0.0.1:3306)/user_info")
	if err != nil {
		log.Fatalln(err)
	}

	persons := []Person{}
	err = db.Select(&persons, "select * from user where id < 5")
	if err != nil {
		log.Fatalln(err)
		return
	}
	tom, lilei := persons[0], persons[1]
	fmt.Printf("Tom: %v, Lilei: %v\n", tom, lilei)
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
	// user, _ := GetByIDSimpler(db, 3)
	// log.Println(user)

	// user := &User{name: "Tom Hanks", sex: 1, descp: "football"}
	// InsertUser(db, user)

	err = UpdateUser(db, 4, "bird")
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

func InsertUser(db *sql.DB, user *User) error {
	stmt, err := db.Prepare("insert into user (name, sex, descp) values (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.name, user.sex, user.descp)
	if err != nil {
		return err
	}

	lastID, _ := res.LastInsertId()
	rowCnt, _ := res.RowsAffected()
	log.Printf("ID = %d, affected = %d\n", lastID, rowCnt)
	return nil
}

func UpdateUser(db *sql.DB, id int, descp string) error {
	stmt, err := db.Prepare("update user set descp = ? where id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(descp, id)
	if err != nil {
		return err
	}
	return nil
}
