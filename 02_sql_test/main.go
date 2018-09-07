package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:password@tcp(192.168.1.104:3036)/mytest.db")
	checkErr(err)
	defer db.Close()

	insert, err := db.Query("INSERT INTO users VALUES('ELIOT')")
	checkErr(err)
	defer insert.Close()
	// // insert
	// stmt, err := db.Prepare("INSERT userinfo SET username=?,departname=?,created=?")
	// checkErr(err)

	// res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
	// checkErr(err)

	// id, err := res.LastInsertId()
	// checkErr(err)

	// fmt.Println(id)
	// // update
	// stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	// checkErr(err)

	// res, err = stmt.Exec("astaxieupdate", id)
	// checkErr(err)

	// affect, err := res.RowsAffected()
	// checkErr(err)

	// fmt.Println(affect)

	// // query
	// rows, err := db.Query("SELECT * FROM userinfo")
	// checkErr(err)

	// for rows.Next() {
	// 	var uid int
	// 	var username string
	// 	var department string
	// 	var created string
	// 	err = rows.Scan(&uid, &username, &department, &created)
	// 	checkErr(err)
	// 	fmt.Println(uid)
	// 	fmt.Println(username)
	// 	fmt.Println(department)
	// 	fmt.Println(created)
	// }

	// // delete
	// stmt, err = db.Prepare("delete from userinfo where uid=?")
	// checkErr(err)

	// res, err = stmt.Exec(id)
	// checkErr(err)

	// affect, err = res.RowsAffected()
	// checkErr(err)

	// fmt.Println(affect)

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
