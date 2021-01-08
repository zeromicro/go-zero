package spec

func (t PrimitiveType) Name() string {
	return t.RawName
}

func (t DefineStruct) Name() string {
	return t.RawName
}

func (t MapType) Name() string {
	return t.RawName
}

func (t ArrayType) Name() string {
	return t.RawName
}

func (t PointerType) Name() string {
	return t.RawName
}

func (t InterfaceType) Name() string {
	return t.RawName
}
