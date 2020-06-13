package gomysql

import (
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"database/sql"
	"sync"
	"reflect"
)

type Database struct{
	DB *sql.DB
	Mutex sync.RWMutex
}

/*private function for gomysql*/
func str_merge(str []string) string{
	var ret string
	if len(str) == 0{
		panic("str_merge failed")
	}

	for _,s := range str{
		ret += (s + ",")
	}
	ret = ret[:len(ret)-1]
	return ret
}

func NewDB() *Database{
	db := &Database{}
	return db
}

func (d *Database) Connect(database string) error{
	var err error
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", USER, PWD, HOST, PORT, database)
	d.DB, err = sql.Open("mysql", dbDSN)
	if err != nil {
        fmt.Println("MYSQL CONNECT: " + dbDSN)
        return err
	}
	err = d.DB.Ping()
	if err != nil{
        fmt.Println("MYSQL Ping error")
        return err		
	}
	return nil
}

func (d *Database) Create(table string,values []string) error{
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	if len(table) == 0{
		return fmt.Errorf("table name is none")
	}
	if len(values) == 0{
		return fmt.Errorf("values is none")
	}

	value_str := str_merge(values)
	cmd := "create table " + table + "(" + value_str + ") CHARSET=utf8;"
	_,err := d.DB.Query(cmd)
	if err != nil{
		fmt.Println("MYSQL Create:",err)
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

func (d *Database) Search(table string,st reflect.Value,where string) (error,[]interface{}){
	Val := st.Elem()
	if Val.Kind() != reflect.Struct{
		return fmt.Errorf("search type is not point"),nil
	}
	var search []string
	Type := st.Elem().Type()
	for i := 0;i<Type.NumField();i++{
		v := Type.Field(i)
		tag,valid := v.Tag.Lookup("item")
		if valid == true{
			search = append(search,tag)
		}
	}
	if len(search) == 0{
		return fmt.Errorf("cannot find any tag"),nil
	}
	ret := str_merge(search)

	cmd := "select " + ret + " from " + table + " " + where + " ;"
	rows,err := d.DB.Query(cmd)
	if err != nil{
		fmt.Println("MYSQL Search:",err)
		return err,nil
	}	

	var ans []interface{}
	
	for rows.Next(){
		var sqlPtr []interface{}
		for i := 0;i<Type.NumField();i++{
			v := Type.Field(i)
			_,valid := v.Tag.Lookup("item")
			if valid == true{
				sqlPtr = append(sqlPtr,Val.Field(i).Addr().Interface())
			}
			
		}
		err := rows.Scan(sqlPtr...)
		if err != nil{
			fmt.Println(err)
			return err,nil
		}

		item := Val.Interface()
		ans = append(ans,item)
	}
	return nil,ans
}

func (d *Database) Drop(table string){
	cmd := "drop table " + table + ";"
	_,err := d.DB.Query(cmd)
	if err != nil{
		fmt.Println(err)
	}
}