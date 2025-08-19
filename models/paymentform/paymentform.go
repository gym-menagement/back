package paymentform

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnGym
    ColumnPayment
    ColumnType
    ColumnCost
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




