package payment

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnGym
    ColumnOrder
    ColumnMembership
    ColumnCost
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




