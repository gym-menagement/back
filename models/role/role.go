package role

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnGym
    ColumnRole
    ColumnName
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}


type Role int

const (
    _ Role  = iota

    RoleMember
    RoleTrainer
    RoleStaff
    RoleGym_admin
    RolePlatform_admin
)

var Roles = []string{ "", "회원", "트레이너", "직원", "헬스장관리자", "플랫폼관리자" }



func GetRole(value Role) string {
    i := int(value)
    if i <= 0 || i >= len(Roles) {
        return ""
    }
     
    return Roles[i]
}

func FindRole(value string) Role {
    for i := 1; i < len(Roles); i++ {
        if Roles[i] == value {
            return Role(i)
        }
    }
     
    return 0
}

func ConvertRole(value []int) []Role {
     items := make([]Role, 0)

     for item := range value {
         items = append(items, Role(item))
     }
     
     return items
}

