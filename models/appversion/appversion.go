package appversion

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnPlatform
    ColumnVersion
    ColumnMinversion
    ColumnForceupdate
    ColumnUpdatemessage
    ColumnDownloadurl
    ColumnStatus
    ColumnReleasedate
    ColumnCreateddate
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}


type Forceupdate int

const (
    _ Forceupdate  = iota

    Forceupdate0
    Forceupdate1
)

var Forceupdates = []string{ "", "아니오", "예" }

type Status int

const (
    _ Status  = iota

    Status0
    Status1
)

var Statuss = []string{ "", "비활성", "활성" }



func GetForceupdate(value Forceupdate) string {
    i := int(value)
    if i <= 0 || i >= len(Forceupdates) {
        return ""
    }
     
    return Forceupdates[i]
}

func FindForceupdate(value string) Forceupdate {
    for i := 1; i < len(Forceupdates); i++ {
        if Forceupdates[i] == value {
            return Forceupdate(i)
        }
    }
     
    return 0
}

func ConvertForceupdate(value []int) []Forceupdate {
     items := make([]Forceupdate, 0)

     for item := range value {
         items = append(items, Forceupdate(item))
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

