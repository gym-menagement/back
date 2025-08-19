package user

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnLoginid
    ColumnPasswd
    ColumnEmail
    ColumnName
    ColumnTel
    ColumnAddress
    ColumnImage
    ColumnSex
    ColumnBirth
    ColumnType
    ColumnConnectid
    ColumnLevel
    ColumnRole
    ColumnUse
    ColumnLogindate
    ColumnLastchangepasswddate
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}


type Level int

const (
    _ Level  = iota

    LevelNormal
    LevelManager
    LevelAdmin
    LevelSuperadmin
    LevelRootadmin
)

var Levels = []string{ "", "일반", "팀장", "관리자", "승인관리자", "전체관리자" }

type Use int

const (
    _ Use  = iota

    UseUse
    UseNotuse
)

var Uses = []string{ "", "사용", "사용안함" }

type Type int

const (
    _ Type  = iota

    TypeNormal
    TypeKakao
    TypeNaver
    TypeGoogle
    TypeApple
)

var Types = []string{ "", "일반", "카카오", "네이버", "구글", "애플" }

type Role int

const (
    _ Role  = iota

    RoleSupervisor
    RoleCoach
    RoleParent
    RolePlayer
    RoleUse
    RoleNormal
)

var Roles = []string{ "", "감독", "코치", "학부모", "현역선수", "동호회", "일반인" }



func GetLevel(value Level) string {
    i := int(value)
    if i <= 0 || i >= len(Levels) {
        return ""
    }
     
    return Levels[i]
}

func FindLevel(value string) Level {
    for i := 1; i < len(Levels); i++ {
        if Levels[i] == value {
            return Level(i)
        }
    }
     
    return 0
}

func ConvertLevel(value []int) []Level {
     items := make([]Level, 0)

     for item := range value {
         items = append(items, Level(item))
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

