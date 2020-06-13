package gomysql

import (
	"fmt"
	"testing"
	"reflect"
	"strconv"
)


type Search struct{
	Id int `item:"id"`
	Name string `item:"name"`
}

var UnitTest = []Search{
	Search{4,"小明"},
	Search{5,"小王"},
}
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

	insert := make(map[string]string)
	for _,unit := range UnitTest{
		insert["id"] = strconv.Itoa(unit.Id)
		insert["name"] = "'" + unit.Name + "'"
		db.Insert(table,insert)
	}

	s := Search{}
	err,datas := db.Search(table,reflect.ValueOf(&s))

	if err != nil{
		t.Errorf("Mysql search failed!")
		return		
	}
	for index,data := range datas{
		d,ok := data.(Search)
		if ok == false{
			t.Errorf("Get Error Struct Type")
			return				
		} 
		if d.Id != UnitTest[index].Id || d.Name != UnitTest[index].Name {
			fmt.Println("real:",d.Id,d.Name," // expect:",UnitTest[index].Id,UnitTest[index].Name)
			t.Errorf("Get Value Error")
			return		
		}	
	}
}