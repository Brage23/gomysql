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

var UnitTest2 = []Search{
	Search{4,"小红"},
	Search{4,"小张"},
	Search{6,"小钱"},
}

func TestMysql(t *testing.T){
	//1.connect to the sql
	db := NewDB()
	err := db.Connect("test")
	if err != nil{
		t.Errorf("step1:Mysql connect failed!")
		return
	}

	//2.create new table
	table := "GolangTest"
	values := []string{"id int not null","name varchar(32) not null"}
	db.Drop(table)
	err = db.Create(table,values)
	if err != nil{
		t.Errorf("step2:Mysql table create failed!")
		return
	}

	//3.insert information to sql one-by-one
	insert := make(map[string]string)
	for _,unit := range UnitTest{
		insert["id"] = strconv.Itoa(unit.Id)
		insert["name"] = "'" + unit.Name + "'"
		db.Insert(table,insert)
	}

	//4.get information from sql 
	s := Search{}
	err,datas := db.Search(table,reflect.ValueOf(&s),"")

	if err != nil{
		t.Errorf("step4:Mysql search failed!")
		return		
	}
	for index,data := range datas{
		d,ok := data.(Search)
		if ok == false{
			t.Errorf("step4:Get Error Struct Type")
			return				
		} 
		if d.Id != UnitTest[index].Id || d.Name != UnitTest[index].Name {
			fmt.Println("real:",d.Id,d.Name," // expect:",UnitTest[index].Id,UnitTest[index].Name)
			t.Errorf("step4:Get Value Error")
			return		
		}	
	}

	//5.insert to sql with several items together
	keys := []string{"id","name"}
	var rows [][]string
	for _,unit := range UnitTest2{
		row := []string{strconv.Itoa(unit.Id),"'" + unit.Name + "'"}
		rows = append(rows,row)
	}

	err = db.InsertRows(table,keys,rows)
	if err != nil{
		t.Errorf("step5:Mysql insert rows failed!")
		return		
	}

	//6.search information by ID = 4
	err,datas = db.Search(table,reflect.ValueOf(&s),"where id = 4")

	if err != nil{
		t.Errorf("step6:Mysql search failed!")
		return		
	}

	for _,data := range datas{
		d,ok := data.(Search)
		if ok == false{
			t.Errorf("step6:Get Error Struct Type")
			return				
		} 
		if d.Id != 4{
			t.Errorf("step6:Get Error ID number")
			return	
		}
		if d.Name != "小明" && d.Name != "小红" && d.Name != "小张"{
			fmt.Println(d.Name)
			t.Errorf("step6:Get Error Number")
			return	
		}
	}
}