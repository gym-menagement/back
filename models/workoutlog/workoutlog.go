package workoutlog

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnUser
    ColumnAttendance
    ColumnHealth
    ColumnExercisename
    ColumnSets
    ColumnReps
    ColumnWeight
    ColumnDuration
    ColumnCalories
    ColumnNote
    ColumnDate
)

type Params struct {
    Column Column
    Value interface{}
}




