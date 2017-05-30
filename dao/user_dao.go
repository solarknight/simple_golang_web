package dao

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/solarknight/simple_golang_web/common"
)

func ConnectDB() {
	db, err := sqlx.Connect("mysql", "root:solarknight@tcp(127.0.0.1:3306)/user_info")
	if err != nil {
		log.Fatalln(err)
	}

	persons := []common.User{}
	err = db.Select(&persons, "select * from user where id < 5")
	if err != nil {
		log.Fatalln(err)
		return
	}
	tom, lilei := persons[0], persons[1]
	fmt.Printf("Tom: %v, Lilei: %v\n", tom, lilei)
}

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:solarknight@tcp(127.0.0.1:3306)/user_info")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetByIDSimpler(db *sql.DB, id int) (*common.User, error) {
	user := new(common.User)
	err := db.QueryRow("select id, name, sex, descp from user where id = ?", id).Scan(&user.Id, &user.Name, &user.Sex, &user.Descp)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return user, nil
}

func GetByID(db *sql.DB, id int) (*common.User, error) {
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

	user := new(common.User)
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Sex, &user.Descp)
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

func InsertUser(db *sql.DB, user *common.User) error {
	stmt, err := db.Prepare("insert into user (name, sex, descp) values (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.Name, user.Sex, user.Descp)
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
