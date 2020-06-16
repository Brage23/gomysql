package gomysql

import(
)

/*
StringMerge(["A","B","C"]) -> "A,B,C"
*/
func StringMerge(str []string) string{
	var ret string
	if len(str) == 0{
		panic("str_merge failed")
	}

	for _,s := range str{
		if len(s) == 0{
			ret += "NULL,"
		} else{
			ret += (s + ",")
		}
		
	}
	ret = ret[:len(ret)-1]
	return ret
}

/*
ParenPackage("A") -> "(A)"
*/
func ParenPackage(str string) string{
	var ret string
	ret = "(" + str + ")"
	return ret
}