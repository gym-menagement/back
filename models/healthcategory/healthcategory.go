package healthcategory

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnGym
    ColumnName
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




