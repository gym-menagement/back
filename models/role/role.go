package role

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnGym
    ColumnRole
    ColumnName
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




