package notificationsetting

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnUser
    ColumnEnabled
    ColumnMembershipexpiry
    ColumnMembershipnear
    ColumnAttendanceenc
    ColumnGymannounce
    ColumnSystemnotice
    ColumnPaymentconfirm
    ColumnPauseexpiry
    ColumnWeeklygoal
    ColumnPersonalrecord
    ColumnQuietenabled
    ColumnQuietstart
    ColumnQuietend
    ColumnCreateddate
    ColumnUpdateddate
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}


type Enabled int

const (
    _ Enabled  = iota

    EnabledEnabled
    EnabledDisabled
)

var Enableds = []string{ "", "활성화", "비활성화" }

type Membershipexpiry int

const (
    _ Membershipexpiry  = iota

    MembershipexpiryEnabled
    MembershipexpiryDisabled
)

var Membershipexpirys = []string{ "", "활성화", "비활성화" }

type Membershipnear int

const (
    _ Membershipnear  = iota

    MembershipnearEnabled
    MembershipnearDisabled
)

var Membershipnears = []string{ "", "활성화", "비활성화" }

type Attendanceenc int

const (
    _ Attendanceenc  = iota

    AttendanceencEnabled
    AttendanceencDisabled
)

var Attendanceencs = []string{ "", "활성화", "비활성화" }

type Gymannounce int

const (
    _ Gymannounce  = iota

    GymannounceEnabled
    GymannounceDisabled
)

var Gymannounces = []string{ "", "활성화", "비활성화" }

type Systemnotice int

const (
    _ Systemnotice  = iota

    SystemnoticeEnabled
    SystemnoticeDisabled
)

var Systemnotices = []string{ "", "활성화", "비활성화" }

type Paymentconfirm int

const (
    _ Paymentconfirm  = iota

    PaymentconfirmEnabled
    PaymentconfirmDisabled
)

var Paymentconfirms = []string{ "", "활성화", "비활성화" }

type Pauseexpiry int

const (
    _ Pauseexpiry  = iota

    PauseexpiryEnabled
    PauseexpiryDisabled
)

var Pauseexpirys = []string{ "", "활성화", "비활성화" }

type Weeklygoal int

const (
    _ Weeklygoal  = iota

    WeeklygoalEnabled
    WeeklygoalDisabled
)

var Weeklygoals = []string{ "", "활성화", "비활성화" }

type Personalrecord int

const (
    _ Personalrecord  = iota

    PersonalrecordEnabled
    PersonalrecordDisabled
)

var Personalrecords = []string{ "", "활성화", "비활성화" }

type Quietenabled int

const (
    _ Quietenabled  = iota

    QuietenabledEnabled
    QuietenabledDisabled
)

var Quietenableds = []string{ "", "활성화", "비활성화" }



func GetEnabled(value Enabled) string {
    i := int(value)
    if i <= 0 || i >= len(Enableds) {
        return ""
    }
     
    return Enableds[i]
}

func FindEnabled(value string) Enabled {
    for i := 1; i < len(Enableds); i++ {
        if Enableds[i] == value {
            return Enabled(i)
        }
    }
     
    return 0
}

func ConvertEnabled(value []int) []Enabled {
     items := make([]Enabled, 0)

     for item := range value {
         items = append(items, Enabled(item))
     }
     
     return items
}

func GetMembershipexpiry(value Membershipexpiry) string {
    i := int(value)
    if i <= 0 || i >= len(Membershipexpirys) {
        return ""
    }
     
    return Membershipexpirys[i]
}

func FindMembershipexpiry(value string) Membershipexpiry {
    for i := 1; i < len(Membershipexpirys); i++ {
        if Membershipexpirys[i] == value {
            return Membershipexpiry(i)
        }
    }
     
    return 0
}

func ConvertMembershipexpiry(value []int) []Membershipexpiry {
     items := make([]Membershipexpiry, 0)

     for item := range value {
         items = append(items, Membershipexpiry(item))
     }
     
     return items
}

func GetMembershipnear(value Membershipnear) string {
    i := int(value)
    if i <= 0 || i >= len(Membershipnears) {
        return ""
    }
     
    return Membershipnears[i]
}

func FindMembershipnear(value string) Membershipnear {
    for i := 1; i < len(Membershipnears); i++ {
        if Membershipnears[i] == value {
            return Membershipnear(i)
        }
    }
     
    return 0
}

func ConvertMembershipnear(value []int) []Membershipnear {
     items := make([]Membershipnear, 0)

     for item := range value {
         items = append(items, Membershipnear(item))
     }
     
     return items
}

func GetAttendanceenc(value Attendanceenc) string {
    i := int(value)
    if i <= 0 || i >= len(Attendanceencs) {
        return ""
    }
     
    return Attendanceencs[i]
}

func FindAttendanceenc(value string) Attendanceenc {
    for i := 1; i < len(Attendanceencs); i++ {
        if Attendanceencs[i] == value {
            return Attendanceenc(i)
        }
    }
     
    return 0
}

func ConvertAttendanceenc(value []int) []Attendanceenc {
     items := make([]Attendanceenc, 0)

     for item := range value {
         items = append(items, Attendanceenc(item))
     }
     
     return items
}

func GetGymannounce(value Gymannounce) string {
    i := int(value)
    if i <= 0 || i >= len(Gymannounces) {
        return ""
    }
     
    return Gymannounces[i]
}

func FindGymannounce(value string) Gymannounce {
    for i := 1; i < len(Gymannounces); i++ {
        if Gymannounces[i] == value {
            return Gymannounce(i)
        }
    }
     
    return 0
}

func ConvertGymannounce(value []int) []Gymannounce {
     items := make([]Gymannounce, 0)

     for item := range value {
         items = append(items, Gymannounce(item))
     }
     
     return items
}

func GetSystemnotice(value Systemnotice) string {
    i := int(value)
    if i <= 0 || i >= len(Systemnotices) {
        return ""
    }
     
    return Systemnotices[i]
}

func FindSystemnotice(value string) Systemnotice {
    for i := 1; i < len(Systemnotices); i++ {
        if Systemnotices[i] == value {
            return Systemnotice(i)
        }
    }
     
    return 0
}

func ConvertSystemnotice(value []int) []Systemnotice {
     items := make([]Systemnotice, 0)

     for item := range value {
         items = append(items, Systemnotice(item))
     }
     
     return items
}

func GetPaymentconfirm(value Paymentconfirm) string {
    i := int(value)
    if i <= 0 || i >= len(Paymentconfirms) {
        return ""
    }
     
    return Paymentconfirms[i]
}

func FindPaymentconfirm(value string) Paymentconfirm {
    for i := 1; i < len(Paymentconfirms); i++ {
        if Paymentconfirms[i] == value {
            return Paymentconfirm(i)
        }
    }
     
    return 0
}

func ConvertPaymentconfirm(value []int) []Paymentconfirm {
     items := make([]Paymentconfirm, 0)

     for item := range value {
         items = append(items, Paymentconfirm(item))
     }
     
     return items
}

func GetPauseexpiry(value Pauseexpiry) string {
    i := int(value)
    if i <= 0 || i >= len(Pauseexpirys) {
        return ""
    }
     
    return Pauseexpirys[i]
}

func FindPauseexpiry(value string) Pauseexpiry {
    for i := 1; i < len(Pauseexpirys); i++ {
        if Pauseexpirys[i] == value {
            return Pauseexpiry(i)
        }
    }
     
    return 0
}

func ConvertPauseexpiry(value []int) []Pauseexpiry {
     items := make([]Pauseexpiry, 0)

     for item := range value {
         items = append(items, Pauseexpiry(item))
     }
     
     return items
}

func GetWeeklygoal(value Weeklygoal) string {
    i := int(value)
    if i <= 0 || i >= len(Weeklygoals) {
        return ""
    }
     
    return Weeklygoals[i]
}

func FindWeeklygoal(value string) Weeklygoal {
    for i := 1; i < len(Weeklygoals); i++ {
        if Weeklygoals[i] == value {
            return Weeklygoal(i)
        }
    }
     
    return 0
}

func ConvertWeeklygoal(value []int) []Weeklygoal {
     items := make([]Weeklygoal, 0)

     for item := range value {
         items = append(items, Weeklygoal(item))
     }
     
     return items
}

func GetPersonalrecord(value Personalrecord) string {
    i := int(value)
    if i <= 0 || i >= len(Personalrecords) {
        return ""
    }
     
    return Personalrecords[i]
}

func FindPersonalrecord(value string) Personalrecord {
    for i := 1; i < len(Personalrecords); i++ {
        if Personalrecords[i] == value {
            return Personalrecord(i)
        }
    }
     
    return 0
}

func ConvertPersonalrecord(value []int) []Personalrecord {
     items := make([]Personalrecord, 0)

     for item := range value {
         items = append(items, Personalrecord(item))
     }
     
     return items
}

func GetQuietenabled(value Quietenabled) string {
    i := int(value)
    if i <= 0 || i >= len(Quietenableds) {
        return ""
    }
     
    return Quietenableds[i]
}

func FindQuietenabled(value string) Quietenabled {
    for i := 1; i < len(Quietenableds); i++ {
        if Quietenableds[i] == value {
            return Quietenabled(i)
        }
    }
     
    return 0
}

func ConvertQuietenabled(value []int) []Quietenabled {
     items := make([]Quietenabled, 0)

     for item := range value {
         items = append(items, Quietenabled(item))
     }
     
     return items
}

