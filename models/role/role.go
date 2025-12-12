package role

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnGym
    ColumnRoleid
    ColumnName
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}


type Roleid int

const (
    _ Roleid  = iota

    RoleidMember
    RoleidTrainer
    RoleidStaff
    RoleidGym_admin
    RoleidPlatform_admin
)

var Roleids = []string{ "", "회원", "트레이너", "직원", "헬스장관리자", "플랫폼관리자" }



func GetRoleid(value Roleid) string {
    i := int(value)
    if i <= 0 || i >= len(Roleids) {
        return ""
    }
     
    return Roleids[i]
}

func FindRoleid(value string) Roleid {
    for i := 1; i < len(Roleids); i++ {
        if Roleids[i] == value {
            return Roleid(i)
        }
    }
     
    return 0
}

func ConvertRoleid(value []int) []Roleid {
     items := make([]Roleid, 0)

     for item := range value {
         items = append(items, Roleid(item))
     }
     
     return items
}

