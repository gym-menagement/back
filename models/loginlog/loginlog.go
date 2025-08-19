package loginlog

type Column int

const (
    _ Column = iota
    ColumnId
    ColumnIp
    ColumnIpvalue
    ColumnDate
    ColumnUser

)

type Params struct {
    Column Column
    Value interface{}
}




