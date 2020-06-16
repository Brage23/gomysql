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

func NewDB() *Database{
	db := &Database{}
	return db
}

func (d *Database) Connect(database string,opts ...ConnFunc) error{
	var err error
	d.Mutex.Lock()
	defer d.Mutex.Unlock()

	opt := ConnOpt{
		Charset : "",
	}
	for _,fun := range opts{
		fun(&opt)
	}

	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", USER, PWD, HOST, PORT, database)

	if len(opt.Charset) != 0{
		dbDSN = dbDSN + "?charset=" + opt.Charset
	}
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

	value_str := StringMerge(values)
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

	key_str := ParenPackage(StringMerge(keys))

	var values_list []string
	for _,avalue := range values{
		if len(avalue) > len(keys){
			fmt.Println("MYSQL too many values")
			panic(1)
		}
		value_str := ParenPackage(StringMerge(avalue))
		values_list = append(values_list,value_str)
	}
	value_s_str := StringMerge(values_list)
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

func (d *Database) Search(table string,st reflect.Value,opts ...SearchFunc) (error,[]interface{}){
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()

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
	ret := StringMerge(search)

	opt := SearchOpt{
		Where : "",
		Order : "",
	}
	for _,fun := range opts{
		fun(&opt)
	}
	cmd := "select " + ret + " from " + table
	if len(opt.Where) != 0 {
		cmd = cmd + " where " + opt.Where
	} 

	if len(opt.Order) != 0{
		cmd = cmd + " ORDER BY " + opt.Order
	}
	cmd = cmd + ";"

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