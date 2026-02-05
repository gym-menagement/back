package notificationhistory

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnSender
    ColumnReceiver
    ColumnGym
    ColumnType
    ColumnTitle
    ColumnBody
    ColumnData
    ColumnStatus
    ColumnErrormessage
    ColumnSentdate
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
    TypeMembership_expiry
    TypeMembership_near_expiry
    TypeAttendance_encourage
    TypeGym_announcement
    TypeSystem_notice
    TypePayment_confirm
    TypePause_expiry
    TypeWeekly_goal_achieved
    TypePersonal_record
)

var Types = []string{ "", "일반", "이용권만료", "이용권임박", "출석독려", "체육관공지", "시스템공지", "결제확인", "일시정지만료", "주간목표달성", "개인기록갱신" }

type Status int

const (
    _ Status  = iota

    StatusPending
    StatusSuccess
    StatusFailed
)

var Statuss = []string{ "", "대기중", "성공", "실패" }



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

