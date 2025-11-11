package trainermember

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnTrainer
    ColumnMember
    ColumnGym
    ColumnStartdate
    ColumnEnddate
    ColumnStatus
    ColumnNote
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
)

var Statuss = []string{ "", "종료", "진행중" }



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

