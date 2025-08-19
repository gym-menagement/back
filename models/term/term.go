package term

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnGym
    ColumnDaytype
    ColumnName
    ColumnTerm
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




