package alarm

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnTitle
    ColumnContent
    ColumnType
    ColumnStatus
    ColumnUser
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}


type Type int

const (
    _ Type  = iota

    TypeNotice
    TypeWarning
    TypeError
    TypeInfo
)

var Types = []string{ "", "공지", "경고", "에러", "정보" }

type Status int

const (
    _ Status  = iota

    StatusSuccess
    StatusFail
    StatusPending
)

var Statuss = []string{ "", "성공", "실패", "대기" }



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

