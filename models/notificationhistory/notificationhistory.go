package notificationhistory

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnSender
    ColumnReceiver
    ColumnGym
    ColumnType
    ColumnTitle
    ColumnBody
    ColumnData
    ColumnStatus
    ColumnErrormessage
    ColumnSentdate
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




