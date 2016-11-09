package db

import "github.com/zzkkff/lib4go/db/tpl"

//IDB 数据库操作接口
type IDB interface {
	Query(string, ...interface{}) ([]map[string]interface{}, []string, error)
	Execute(string, ...interface{}) (int64, error)
	Begin() (IDBTrans, error)
}

//IDBTrans 数据库事务接口
type IDBTrans interface {
	Query(string, ...interface{}) ([]map[string]interface{}, []string, error)
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
func NewDB(provider string, connString string, maxIdle int, maxOpen int) (obj *DB, err error) {
	obj = &DB{}
	obj.tpl, err = tpl.GetDBContext(provider)
	if err != nil {
		return
	}
	obj.db, err = NewSysDB(provider, connString, maxIdle, maxOpen)
	return
}

//Query 查询数据
func (db *DB) Query(sql string, input map[string]interface{}) (data []map[string]interface{}, query string, args []interface{}, err error) {
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
	row, err = db.db.Execute(query, args)
	return
}

//ExecuteSP 根据包含@名称占位符的语句执行查询语句
func (db *DB) ExecuteSP(sql string, input map[string]interface{}) (row int64, query string, args []interface{}, err error) {
	query, args = db.tpl.GetSPContext(sql, input)
	row, err = db.db.Execute(query, args)
	return
}

//Replace 替换SQL语句中的参数
func (db *DB) Replace(sql string, args []interface{}) string {
	return db.tpl.Replace(sql, args)
}

//Begin 创建事务
func (db *DB) Begin() (t IDBTrans, err error) {
	return db.db.Begin()
}
