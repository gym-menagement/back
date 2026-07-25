package notificationsetting

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnUser
    ColumnEnabled
    ColumnMembershipexpiry
    ColumnMembershipnear
    ColumnAttendanceenc
    ColumnGymannounce
    ColumnSystemnotice
    ColumnPaymentconfirm
    ColumnPauseexpiry
    ColumnWeeklygoal
    ColumnPersonalrecord
    ColumnQuietenabled
    ColumnQuietstart
    ColumnQuietend
    ColumnCreateddate
    ColumnUpdateddate
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




