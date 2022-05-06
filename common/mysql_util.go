package common


import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type MysqlClient struct {
	Host string `json:"username"`
	Port int `json:"username"`
	Dbname string `json:"username"`
	Username string `json:"username"`
	Password string `json:"username"`
	MysqlConn sql.DB `json:"mysql_conn"`
}


func (m *MysqlClient) GetConn()  {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",m.Username,m.Password,m.Host,m.Port,m.Dbname)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetConnMaxIdleTime(5)
	conn.SetMaxOpenConns(5)
	err = conn.Ping()
	if err != nil {
		Error.Println("Ping Error! Mysql!")
		return
	}
	m.MysqlConn = *conn
	//fmt.Println("connect mysql success!")
}

func (m *MysqlClient) CloesConn() error {
	return m.MysqlConn.Close()
}

// 获取表数据
func (m *MysqlClient) Query(sql string) []map[string]interface{}{
	rows, err := m.MysqlConn.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	return getQueryResult(rows)
}

// DML操作
func (m *MysqlClient) Exec(sql string,arg...interface{}) error{
	tx, err := m.MysqlConn.Begin()
	if err != nil {
		Error.Println("open mysql database fail", err)
		return err
	}

	stmt, err := m.MysqlConn.Prepare(sql)
	if err != nil {
		Error.Println(err)
		tx.Rollback()
	}
	res, err := stmt.Exec()
	if err != nil {
		Error.Println(err)
		tx.Rollback()
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		Error.Println(err)
		tx.Rollback()
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		Error.Println(err)
		tx.Rollback()
	}
	fmt.Printf("ID=%d, affected=%d\n", lastId, rowCnt)
	return tx.Commit()
}




// 插入数据
func (m *MysqlClient) Insert(sql string) {
	//stmt, err := mysqlConn.Prepare("INSERT INTO user(name, age) VALUES(?, ?);")
	stmt, err := m.MysqlConn.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID=%d, affected=%d\n", lastId, rowCnt)
}

// 删除数据
func (m *MysqlClient) Delete(sql string) {
	stmt, err := m.MysqlConn.Prepare("DELETE FROM user WHERE name='python'")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec()
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID=%d, affected=%d\n", lastId, rowCnt)
}

// 更新数据
func (m *MysqlClient) Update(sql string) {
	//stmt, err := mysqlConn.Prepare("UPDATE user SET age=27 WHERE name='python'")
	stmt, err := m.MysqlConn.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec()
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID=%d, affected=%d\n", lastId, rowCnt)
}

//查询结果转换
func getQueryResult(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据

	for index, _ := range cache { //为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}

	var list []map[string]interface{} //返回的切片
	for rows.Next() {
		_ = rows.Scan(cache...)

		item := make(map[string]interface{})
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{}) //取实际类型
		}
		for k, v := range item {
			if v == nil {
				continue
			}
			item[k] = string(v.([]uint8))
		}
		list = append(list, item)
	}
	_ = rows.Close()
	return list
}
