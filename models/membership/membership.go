package membership

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnGym
    ColumnUser
    ColumnName
    ColumnSex
    ColumnBirth
    ColumnPhonenum
    ColumnAddress
    ColumnImage
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




