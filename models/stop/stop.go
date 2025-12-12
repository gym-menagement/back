package stop

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnUsehealth
    ColumnStartday
    ColumnEndday
    ColumnCount
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




