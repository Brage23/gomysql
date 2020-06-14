package gomysql

import (
	//"fmt"
)

type ConnOpt struct{
	Charset string
}

type ConnFunc func (c *ConnOpt)

func WithCharset(charset string) ConnFunc{
	return func(c *ConnOpt){
		c.Charset = charset
	}
}


type SearchOpt struct{
	Where string
	Order string
}

type SearchFunc func (s *SearchOpt)

func WithWhere(where string) SearchFunc{
	return func(s *SearchOpt){
		s.Where = where
	}
}

func WithOrder(order string) SearchFunc{
	return func(s *SearchOpt){
		s.Order = order
	}
}