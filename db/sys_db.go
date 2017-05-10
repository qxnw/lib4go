package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"time"

	_ "github.com/mattn/go-oci8"
	_ "github.com/mattn/go-sqlite3"
	_ "gopkg.in/rana/ora.v4"
)

/*
github.com/mattn/go-oci8

http://www.simonzhang.net/?p=2890
http://blog.sina.com.cn/s/blog_48c95a190102w2ln.html
http://www.tudou.com/programs/view/yet9OngrV_4/
https://github.com/wendal/go-oci8/downloads
https://github.com/wendal/go-oci8

安装方法
1. 下载：http://www.oracle.com/technetwork/database/features/instant-client/index.html
2. 解压文件 unzip instantclient-basic-linux.x64-12.1.0.1.0.zip -d /usr/local/
3. 配置环境变量
vi .bash_profile
export ora_home=/usr/local/instantclient_12_1
export PATH=$PATH:$ora_home
export LD_LIBRARY_PATH=$ora_home


*/

const (
	//SQLITE3 Sqlite3数据库
	SQLITE3 = "sqlite3"
	//OCI8 oralce数据库
	OCI8 = "oci8"
	//ora  oralce数据库
	ORA = "ora"
)

//SysDB 数据库实体
type SysDB struct {
	provider   string
	connString string
	db         *sql.DB
	maxIdle    int
	maxOpen    int
}

//NewSysDB 创建DB实例
func NewSysDB(provider string, connString string, max int) (obj *SysDB, err error) {
	/*add by champly 2016年11月14日17:03:46*/
	if provider == "" || connString == "" {
		err = errors.New("provider or connString not allow nil")
		return
	}
	/*end*/

	obj = &SysDB{provider: provider, connString: connString}
	switch strings.ToLower(provider) {
	case "ora":
		obj.db, err = sql.Open(ORA, connString)
	case "oracle":
		obj.db, err = sql.Open(OCI8, connString)
	case "sqlite":
		obj.db, err = sql.Open(SQLITE3, connString)
	default:
		return nil, errors.New("数据库类型不支持:" + provider)
	}
	if err != nil {
		return
	}
	if max > 0 {
		obj.db.SetMaxIdleConns(max)
		obj.db.SetMaxOpenConns(max)
	}
	obj.db.SetConnMaxLifetime(time.Second * 300)
	return
}

//Query 执行SQL查询语句
func (db *SysDB) Query(query string, args ...interface{}) (dataRows []map[string]interface{}, colus []string, err error) {

	rows, err := db.db.Query(query, args...)
	if err != nil {
		if rows != nil {
			rows.Close()
		}
		return
	}
	dataRows, colus, err = resolveRows(rows, 0)
	rows.Close()
	return

}

func resolveRows(rows *sql.Rows, col int) (dataRows []map[string]interface{}, columns []string, err error) {
	dataRows = make([]map[string]interface{}, 0)
	colus, err := rows.Columns()
	if err != nil {
		return
	}
	columns = make([]string, 0, len(colus))
	for _, v := range colus {
		columns = append(columns, strings.ToLower(v))
	}

	for rows.Next() {
		row := make(map[string]interface{})
		dataRows = append(dataRows, row)
		var buffer []interface{}
		for index := 0; index < len(columns); index++ {
			var va []byte
			buffer = append(buffer, &va)
		}
		err = rows.Scan(buffer...)
		if err != nil {
			return
		}
		for index := 0; index < len(columns) && (index < col || col == 0); index++ {
			key := columns[index]
			value := buffer[index]
			if value == nil {
				continue
			} else {
				row[key] = strings.TrimPrefix(fmt.Sprintf("%s", value), "&")
			}
		}
	}
	return
}

//Execute 执行SQL操作语句
func (db *SysDB) Execute(query string, args ...interface{}) (affectedRow int64, err error) {
	result, err := db.db.Exec(query, args...)
	if err != nil {
		return
	}
	affectedRow, err = result.RowsAffected()
	return
}

//Begin 创建一个事务请求
func (db *SysDB) Begin() (r IDBTrans, err error) {
	t := &SysDBTransaction{}
	t.tx, err = db.db.Begin()
	return t, err
}

func (db *SysDB) Print() {
	fmt.Printf("maxIdle: %+v\n", db.db.Stats())
	fmt.Println("maxOpen: ", db.maxOpen)
}
func (db *SysDB) Close() {
	db.db.Close()
}
