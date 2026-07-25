package usehealthusage

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnGym
    ColumnUsehealth
    ColumnMembership
    ColumnUser
    ColumnAttendance
    ColumnType
    ColumnUsedcount
    ColumnRemainingcount
    ColumnCheckintime
    ColumnCheckouttime
    ColumnDuration
    ColumnNote
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}


type Type int

const (
    _ Type  = iota

    TypeEntry
    TypePt
    TypeClass
)

var Types = []string{ "", "입장", "PT수업", "그룹수업" }



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

