package usehelth

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnOrder
    ColumnHelth
    ColumnUser
    ColumnRocker
    ColumnTerm
    ColumnDiscount
    ColumnStartday
    ColumnEndday
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




