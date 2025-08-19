package order

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnMembership
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




