package ipblock

type Column int

const (
    _ Column = iota
    ColumnId
    ColumnAddress
    ColumnType
    ColumnPolicy
    ColumnUse
    ColumnOrder
    ColumnDate

)

type Params struct {
    Column Column
    Value interface{}
}


type Type int

const (
    _ Type  = iota

    TypeAdmin
    TypeNormal
)

var Types = []string{ "", "관리자 접근", "일반 접근" }

type Policy int

const (
    _ Policy  = iota

    PolicyGrant
    PolicyDeny
)

var Policys = []string{ "", "허용", "거부" }

type Use int

const (
    _ Use  = iota

    UseUse
    UseNotuse
)

var Uses = []string{ "", "사용", "사용안함" }



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

func GetPolicy(value Policy) string {
    i := int(value)
    if i <= 0 || i >= len(Policys) {
        return ""
    }
     
    return Policys[i]
}

func FindPolicy(value string) Policy {
    for i := 1; i < len(Policys); i++ {
        if Policys[i] == value {
            return Policy(i)
        }
    }
     
    return 0
}

func ConvertPolicy(value []int) []Policy {
     items := make([]Policy, 0)

     for item := range value {
         items = append(items, Policy(item))
     }
     
     return items
}

func GetUse(value Use) string {
    i := int(value)
    if i <= 0 || i >= len(Uses) {
        return ""
    }
     
    return Uses[i]
}

func FindUse(value string) Use {
    for i := 1; i < len(Uses); i++ {
        if Uses[i] == value {
            return Use(i)
        }
    }
     
    return 0
}

func ConvertUse(value []int) []Use {
     items := make([]Use, 0)

     for item := range value {
         items = append(items, Use(item))
     }
     
     return items
}

