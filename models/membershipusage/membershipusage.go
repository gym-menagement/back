package membershipusage

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnMembership
    ColumnUser
    ColumnType
    ColumnTotaldays
    ColumnUseddays
    ColumnRemainingdays
    ColumnTotalcount
    ColumnUsedcount
    ColumnRemainingcount
    ColumnStartdate
    ColumnEnddate
    ColumnStatus
    ColumnPausedays
    ColumnLastuseddate
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}


type Type int

const (
    _ Type  = iota

    Type0
    Type1
)

var Types = []string{ "", "기간제", "횟수제" }

type Status int

const (
    _ Status  = iota

    Status0
    Status1
    Status2
    Status3
)

var Statuss = []string{ "", "사용중", "일시정지", "만료", "환불" }



func GetType(value Type) string {
    i := int(value)
    if i <= 0 || i >= len(Types) {
        return ""
    }
     
    return Types[i]
}

func FindType(value string) Type {
    for i := 1; i < len(Types); i++ {
        if Types[i] == value {
            return Type(i)
        }
    }
     
    return 0
}

func ConvertType(value []int) []Type {
     items := make([]Type, 0)

     for item := range value {
         items = append(items, Type(item))
     }
     
     return items
}

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

