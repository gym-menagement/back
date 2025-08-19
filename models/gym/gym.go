package gym

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnName
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




