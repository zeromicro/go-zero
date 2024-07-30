package spec

// Name returns a basic string, such as int32,int64
func (t PrimitiveType) Name() string {
	return t.RawName
}

// Comments returns the comments of struct
func (t PrimitiveType) Comments() []string {
	return nil
}

// Documents returns the documents of struct
func (t PrimitiveType) Documents() []string {
	return nil
}

// Name returns a structure string, such as User
func (t DefineStruct) Name() string {
	return t.RawName
}

// Comments returns the comments of struct
func (t DefineStruct) Comments() []string {
	return nil
}

// Documents returns the documents of struct
func (t DefineStruct) Documents() []string {
	return t.Docs
}

// Name returns a structure string, such as User
func (t NestedStruct) Name() string {
	return t.RawName
}

// Comments returns the comments of struct
func (t NestedStruct) Comments() []string {
	return nil
}

// Documents returns the documents of struct
func (t NestedStruct) Documents() []string {
	return t.Docs
}

// Name returns a map string, such as map[string]int
func (t MapType) Name() string {
	return t.RawName
}

// Comments returns the comments of struct
func (t MapType) Comments() []string {
	return nil
}

// Documents returns the documents of struct
func (t MapType) Documents() []string {
	return nil
}

// Name returns a slice string, such as []int
func (t ArrayType) Name() string {
	return t.RawName
}

// Comments returns the comments of struct
func (t ArrayType) Comments() []string {
	return nil
}

// Documents returns the documents of struct
func (t ArrayType) Documents() []string {
	return nil
}

// Name returns a pointer string, such as *User
func (t PointerType) Name() string {
	return t.RawName
}

// Comments returns the comments of struct
func (t PointerType) Comments() []string {
	return nil
}

// Documents returns the documents of struct
func (t PointerType) Documents() []string {
	return nil
}

// Name returns an interface string, Its fixed value is any
func (t InterfaceType) Name() string {
	return t.RawName
}

// Comments returns the comments of struct
func (t InterfaceType) Comments() []string {
	return nil
}

// Documents returns the documents of struct
func (t InterfaceType) Documents() []string {
	return nil
}
