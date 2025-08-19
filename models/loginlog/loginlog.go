package loginlog

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnIp
    ColumnIpvalue
    ColumnUser
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




