package memberqr

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnUser
    ColumnCode
    ColumnImageurl
    ColumnIsactive
    ColumnExpiredate
    ColumnGenerateddate
    ColumnLastuseddate
    ColumnUsecount
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}


type Isactive int

const (
    _ Isactive  = iota

    IsactiveInactive
    IsactiveActive
)

var Isactives = []string{ "", "비활성", "활성" }



func GetIsactive(value Isactive) string {
    i := int(value)
    if i <= 0 || i >= len(Isactives) {
        return ""
    }
     
    return Isactives[i]
}

func FindIsactive(value string) Isactive {
    for i := 1; i < len(Isactives); i++ {
        if Isactives[i] == value {
            return Isactive(i)
        }
    }
     
    return 0
}

func ConvertIsactive(value []int) []Isactive {
     items := make([]Isactive, 0)

     for item := range value {
         items = append(items, Isactive(item))
     }
     
     return items
}

