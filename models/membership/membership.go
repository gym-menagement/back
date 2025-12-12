package membership

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnUser
    ColumnGym
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




