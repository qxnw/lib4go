package db

import (
	"fmt"
	"time"

	"github.com/qxnw/lib4go/db/tpl"
)

//IDB 数据库操作接口,安装可需能需要执行export LD_LIBRARY_PATH=/usr/local/lib
type IDB interface {
	Query(string, ...interface{}) ([]QueryRow, []string, error)
	Execute(string, ...interface{}) (int64, error)
	Begin() (IDBTrans, error)
	Close()
}

//IDBTrans 数据库事务接口
type IDBTrans interface {
	Query(string, ...interface{}) ([]QueryRow, []string, error)
	Execute(string, ...interface{}) (int64, error)
	Rollback() error
	Commit() error
}

//DB 数据库操作类
type DB struct {
	db  IDB
	tpl tpl.ITPLContext
}

//NewDB 创建DB实例
func NewDB(provider string, connString string, maxOpen int, maxIdle int, maxLifeTime int) (obj *DB, err error) {
	obj = &DB{}
	obj.tpl, err = tpl.GetDBContext(provider)
	if err != nil {
		return
	}
	fmt.Println(provider, connString, maxOpen, maxIdle, maxLifeTime)
	obj.db, err = NewSysDB(provider, connString, maxOpen, maxIdle, time.Duration(maxLifeTime)*time.Second)
	return
}

//GetTPL 获取模板翻译参数
func (db *DB) GetTPL() tpl.ITPLContext {
	return db.tpl
}

//Query 查询数据
func (db *DB) Query(sql string, input map[string]interface{}) (data []QueryRow, query string, args []interface{}, err error) {
	query, args = db.tpl.GetSQLContext(sql, input)
	data, _, err = db.db.Query(query, args...)
	return
}

//Scalar 根据包含@名称占位符的查询语句执行查询语句
func (db *DB) Scalar(sql string, input map[string]interface{}) (data interface{}, query string, args []interface{}, err error) {
	query, args = db.tpl.GetSQLContext(sql, input)
	result, colus, err := db.db.Query(query, args...)
	if err != nil || len(result) == 0 || len(result[0]) == 0 || len(colus) == 0 {
		return
	}
	data = result[0][colus[0]]
	return
}

//Execute 根据包含@名称占位符的语句执行查询语句
func (db *DB) Execute(sql string, input map[string]interface{}) (row int64, query string, args []interface{}, err error) {
	query, args = db.tpl.GetSQLContext(sql, input)
	row, err = db.db.Execute(query, args...)
	return
}

//ExecuteSP 根据包含@名称占位符的语句执行查询语句
func (db *DB) ExecuteSP(sql string, input map[string]interface{}) (row int64, query string, args []interface{}, err error) {
	query, args = db.tpl.GetSPContext(sql, input)
	row, err = db.db.Execute(query, args...)
	return
}

//Replace 替换SQL语句中的参数
func (db *DB) Replace(sql string, args []interface{}) string {
	return db.tpl.Replace(sql, args)
}

//Begin 创建事务
func (db *DB) Begin() (t *DBTrans, err error) {
	tt := &DBTrans{}
	tt.tx, err = db.db.Begin()
	if err != nil {
		return
	}
	tt.tpl = db.tpl
	return tt, nil
}
func (db *DB) Close() {
	db.db.Close()
}
