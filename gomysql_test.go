package gomysql

import (
	//"fmt"
	"testing"
)


func TestMysql(t *testing.T){
	db := NewDB()
	err := db.Connect("test")
	if err != nil{
		t.Errorf("Mysql connect failed!")
		return
	}

	table := "GolangTest"
	values := []string{"id int not null","name varchar(32) not null"}
	db.Drop(table)
	err = db.Create(table,values)
	if err != nil{
		t.Errorf("Mysql table create failed!")
		return
	}
}