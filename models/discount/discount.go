package discount

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnName
    ColumnDiscount
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




