# features
The follows need more discussion.

1. Support struct field \<array\> type?
```api
type Foo {
    Bar [2]int `json:""`
}
```

2. Support struct field type \<StructDataType\> or not?
```api
type Foo {
    Bar{
        Baz int `json:""`
    }`json:""`
}
```

3. Support alias?
```api
type Alias int
```
