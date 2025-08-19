package rocker

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnGroup
    ColumnName
    ColumnAvailable
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




