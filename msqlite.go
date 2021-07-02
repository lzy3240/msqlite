package msqlite

import (
	"database/sql"
	"fmt"

	"github.com/demdxx/gocast"
)

// Msqlite ...
type Msqlite struct {
	DB *sql.DB
}

// NewSqlite ...
func NewSqlite(filepath string) Msqlite {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic("open sqlite faild:" + err.Error())
	}

	return Msqlite{
		DB: db,
	}
}

// Queryby ...
func (m *Msqlite) Queryby(sqlstr string, args ...interface{}) *[]map[string]interface{} {
	// `SELECT * FROM user WHERE mobile=?`
	stmt, err := m.DB.Prepare(sqlstr)
	checkErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	checkErr(err)
	//遍历每一行
	colNames, _ := rows.Columns()
	var cols = make([]interface{}, len(colNames))
	for i := 0; i < len(colNames); i++ {
		cols[i] = new(interface{})
	}
	var maps = make([]map[string]interface{}, 0)
	for rows.Next() {
		err := rows.Scan(cols...)
		checkErr(err)
		var rowMap = make(map[string]interface{})
		for i := 0; i < len(colNames); i++ {
			rowMap[colNames[i]] = convertRow(*(cols[i].(*interface{})))
		}
		maps = append(maps, rowMap)
	}
	//fmt.Println(maps)
	return &maps //返回指针
}

//Modifyby 修改数据操作
func (m *Msqlite) Modifyby(sqlstr string, args ...interface{}) int64 {
	// `INSERT user (uname, age, mobile) VALUES (?, ?, ?)`
	// "update user set mobile=? where id=?"
	// "DELETE FROM user where id=?"
	stmt, err := m.DB.Prepare(sqlstr) // Exec、Prepare均可实现增删改查
	checkErr(err)
	defer stmt.Close()
	res, err := stmt.Exec(args...)
	checkErr(err)
	//判断执行结果
	num, err := res.RowsAffected()
	checkErr(err)
	return num
}

// CloseDB ...
func (m *Msqlite) CloseDB() {
	m.DB.Close()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// convertRow 行数据转换
func convertRow(row interface{}) interface{} {
	switch row.(type) {
	case int:
		return gocast.ToInt(row)
	case int32:
		return gocast.ToFloat32(row)
	case int64:
		return gocast.ToFloat64(row)
	case float32:
		return gocast.ToFloat32(row)
	case float64:
		return gocast.ToFloat64(row)
	case string:
		return gocast.ToString(row)
	case []byte:
		return gocast.ToString(row)
	case bool:
		return gocast.ToBool(row)
	}
	return row
}
