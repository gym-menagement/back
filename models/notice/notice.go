package notice

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnGym
    ColumnTitle
    ColumnContent
    ColumnType
    ColumnIspopup
    ColumnIspush
    ColumnTarget
    ColumnViewcount
    ColumnStartdate
    ColumnEnddate
    ColumnStatus
    ColumnCreatedby
    ColumnCreateddate
    ColumnUpdateddate
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}


type Type int

const (
    _ Type  = iota

    TypeGeneral
    TypeImportant
    TypeEvent
)

var Types = []string{ "", "일반", "중요", "이벤트" }

type Ispopup int

const (
    _ Ispopup  = iota

    IspopupNo
    IspopupYes
)

var Ispopups = []string{ "", "아니오", "예" }

type Ispush int

const (
    _ Ispush  = iota

    IspushNo
    IspushYes
)

var Ispushs = []string{ "", "아니오", "예" }

type Target int

const (
    _ Target  = iota

    TargetAll
    TargetMembers_only
    TargetSpecific_members
)

var Targets = []string{ "", "전체", "회원만", "특정회원" }

type Status int

const (
    _ Status  = iota

    StatusPrivate
    StatusPublic
)

var Statuss = []string{ "", "비공개", "공개" }



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

func GetIspopup(value Ispopup) string {
    i := int(value)
    if i <= 0 || i >= len(Ispopups) {
        return ""
    }
     
    return Ispopups[i]
}

func FindIspopup(value string) Ispopup {
    for i := 1; i < len(Ispopups); i++ {
        if Ispopups[i] == value {
            return Ispopup(i)
        }
    }
     
    return 0
}

func ConvertIspopup(value []int) []Ispopup {
     items := make([]Ispopup, 0)

     for item := range value {
         items = append(items, Ispopup(item))
     }
     
     return items
}

func GetIspush(value Ispush) string {
    i := int(value)
    if i <= 0 || i >= len(Ispushs) {
        return ""
    }
     
    return Ispushs[i]
}

func FindIspush(value string) Ispush {
    for i := 1; i < len(Ispushs); i++ {
        if Ispushs[i] == value {
            return Ispush(i)
        }
    }
     
    return 0
}

func ConvertIspush(value []int) []Ispush {
     items := make([]Ispush, 0)

     for item := range value {
         items = append(items, Ispush(item))
     }
     
     return items
}

func GetTarget(value Target) string {
    i := int(value)
    if i <= 0 || i >= len(Targets) {
        return ""
    }
     
    return Targets[i]
}

func FindTarget(value string) Target {
    for i := 1; i < len(Targets); i++ {
        if Targets[i] == value {
            return Target(i)
        }
    }
     
    return 0
}

func ConvertTarget(value []int) []Target {
     items := make([]Target, 0)

     for item := range value {
         items = append(items, Target(item))
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

