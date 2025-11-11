package rocker

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnGroup
    ColumnName
    ColumnAvailable
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}


type Available int

const (
    _ Available  = iota

    Available0
    Available1
)

var Availables = []string{ "", "사용중", "사용가능" }



func GetAvailable(value Available) string {
    i := int(value)
    if i <= 0 || i >= len(Availables) {
        return ""
    }
     
    return Availables[i]
}

func FindAvailable(value string) Available {
    for i := 1; i < len(Availables); i++ {
        if Availables[i] == value {
            return Available(i)
        }
    }
     
    return 0
}

func ConvertAvailable(value []int) []Available {
     items := make([]Available, 0)

     for item := range value {
         items = append(items, Available(item))
     }
     
     return items
}

