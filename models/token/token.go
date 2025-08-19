package token

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnUser
    ColumnToken
    ColumnStatus
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




