package gymadmin

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnGym
    ColumnUser
    ColumnLevel
    ColumnStatus
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




