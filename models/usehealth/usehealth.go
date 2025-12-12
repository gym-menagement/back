package usehealth

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnOrder
    ColumnHealth
    ColumnMembership
    ColumnUser
    ColumnTerm
    ColumnDiscount
    ColumnStartday
    ColumnEndday
    ColumnGym
    ColumnStatus
    ColumnTotalcount
    ColumnUsedcount
    ColumnRemainingcount
    ColumnQrcode
    ColumnLastuseddate
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}


type Status int

const (
    _ Status  = iota

    StatusTerminated
    StatusUse
    StatusPaused
    StatusExpired
)

var Statuss = []string{ "", "종료", "사용중", "일시정지", "만료" }



func GetStatus(value Status) string {
    i := int(value)
    if i <= 0 || i >= len(Statuss) {
        return ""
    }
     
    return Statuss[i]
}

func FindStatus(value string) Status {
    for i := 1; i < len(Statuss); i++ {
        if Statuss[i] == value {
            return Status(i)
        }
    }
     
    return 0
}

func ConvertStatus(value []int) []Status {
     items := make([]Status, 0)

     for item := range value {
         items = append(items, Status(item))
     }
     
     return items
}

