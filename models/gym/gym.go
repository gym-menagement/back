package gym

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnName
    ColumnAddress
    ColumnTel
    ColumnUser
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




