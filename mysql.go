package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"database/sql"
	"sync"
)

type Database struct{
	DB *sql.DB
	Mutex sync.RWMutex
}

func NewDB() *Database{
	db := &Database{}
	return db
}

func (d *Database) Connect(username string,pwd string,host string,port int,database string) error{
	var err error
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", username, pwd, host, port, database)
	d.DB, err = sql.Open("mysql", dbDSN)
	if err != nil {
        fmt.Println("MYSQL CONNECT: " + dbDSN)
        return err
	}
	return nil
}

func (d *Database) Insert(table string,values map[string]string) error{
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	if len(values) == 0{
		return fmt.Errorf("values is none")
	}
	if len(table) == 0{
		return fmt.Errorf("table name is none")
	}
	var key_list string
	var value_list string

	for key,value := range values{
		key_list += (key + ",")
		value_list += (value + ",")
	}
	key_list = key_list[:len(key_list)-1]
	value_list = value_list[:len(value_list)-1]

	query := "insert into " + table + "(" + key_list + ")" + " values(" + value_list + ");"
	_,err := d.DB.Query(query)
	if err != nil{
		fmt.Println("MYSQL INSERT:",err)
		return err
	}
	return nil
}

func(d *Database) InsertRows(table string,keys []string,values [][]string) error{
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	if (len(values) == 0 || len(keys) == 0){
		return fmt.Errorf("values is none")
	}
	if len(table) == 0{
		return fmt.Errorf("table name is none")
	}

	var key_str string
	var value_s_str string
	for _,key := range keys{
		key_str += (key + ",")
	}
	key_str = key_str[:len(key_str)-1]
	key_str = "(" + key_str + ")"

	for _,avalue := range values{
		var value_str string
		if len(avalue) > len(keys){
			fmt.Println("MYSQL too many values")
			panic(1)
		}
		for _,value := range avalue{
			if len(value) == 0{
				value_str += "NULL,"
			} else{
				value_str += (value + ",")
			}
		}
		value_str = value_str[:len(value_str)-1]
		value_str = "(" + value_str + ")"

		value_s_str += (value_str + ",")
	}

	value_s_str = value_s_str[:len(value_s_str)-1]
	cmd := "insert into " + table + " " + key_str + " values " + value_s_str + ";"
	_,err := d.DB.Query(cmd)
	if err != nil{
		fmt.Println("MYSQL INSERT_ROWS:",err)
		return err
	}
	return nil	
}

func (d *Database) Clear(table string) error{
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	if len(table) == 0{
		return fmt.Errorf("table name is none")
	}
	query := "truncate " + table
	_,err := d.DB.Query(query)
	if err != nil{
		fmt.Println("MYSQL CLEAR:",err)
		return err
	}
	return nil
}