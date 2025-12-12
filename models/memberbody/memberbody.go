package memberbody

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnGym
    ColumnUser
    ColumnHeight
    ColumnWeight
    ColumnBodyfat
    ColumnMusclemass
    ColumnBmi
    ColumnSkeletalmuscle
    ColumnBodywater
    ColumnChest
    ColumnWaist
    ColumnHip
    ColumnArm
    ColumnThigh
    ColumnNote
    ColumnMeasureddate
    ColumnMeasuredby
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




