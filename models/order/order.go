package order

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnUser
    ColumnGym
    ColumnHealth
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




