package health

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnCategory
    ColumnTerm
    ColumnName
    ColumnCount
    ColumnCost
    ColumnDiscount
    ColumnCostdiscount
    ColumnContent
    ColumnGym
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




