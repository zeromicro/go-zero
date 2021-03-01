package spec

// Name returns a basic string, such as int32,int64
func (t PrimitiveType) Name() string {
	return t.RawName
}

// Name returns a structure string, such as User
func (t DefineStruct) Name() string {
	return t.RawName
}

// Name returns a map string, such as map[string]int
func (t MapType) Name() string {
	return t.RawName
}

// Name returns a slice string, such as []int
func (t ArrayType) Name() string {
	return t.RawName
}

// Name returns a pointer string, such as *User
func (t PointerType) Name() string {
	return t.RawName
}

// Name returns a interface string, Its fixed value is interface{}
func (t InterfaceType) Name() string {
	return t.RawName
}
