package membership

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnGym
    ColumnUser
    ColumnName
    ColumnSex
    ColumnBirth
    ColumnPhonenum
    ColumnAddress
    ColumnImage
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}


type Sex int

const (
    _ Sex  = iota

    Sex0
    Sex1
)

var Sexs = []string{ "", "남성", "여성" }



func GetSex(value Sex) string {
    i := int(value)
    if i <= 0 || i >= len(Sexs) {
        return ""
    }
     
    return Sexs[i]
}

func FindSex(value string) Sex {
    for i := 1; i < len(Sexs); i++ {
        if Sexs[i] == value {
            return Sex(i)
        }
    }
     
    return 0
}

func ConvertSex(value []int) []Sex {
     items := make([]Sex, 0)

     for item := range value {
         items = append(items, Sex(item))
     }
     
     return items
}

