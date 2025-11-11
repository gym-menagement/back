package ptreservation

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnTrainer
    ColumnMember
    ColumnGym
    ColumnReservationdate
    ColumnStarttime
    ColumnEndtime
    ColumnDuration
    ColumnStatus
    ColumnNote
    ColumnCancelreason
    ColumnCreateddate
    ColumnUpdateddate
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}


type Status int

const (
    _ Status  = iota

    Status0
    Status1
    Status2
    Status3
)

var Statuss = []string{ "", "예약", "완료", "취소", "노쇼" }



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

