package gorm

import "strconv"


type Var struct {
	d []byte
}

type Value = *Var

func (v *Var)Int() int64 {
	if v.d == nil {
		return 0
	}
	n, _ := strconv.ParseInt(string(v.d), 10, 64)
	return n
}



func (v *Var)Uint() uint64 {
	if v.d == nil {
		return 0
	}
	n, _ := strconv.ParseInt(string(v.d), 10, 64)
	return uint64(n)
}

func (v *Var) String() string {
	return string(v.d)
}

func (v *Var) Bytes() []byte {
	return v.d
}