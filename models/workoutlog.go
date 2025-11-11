package models

import (
    "gym/global/config"
    "gym/models/workoutlog"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Workoutlog struct {
            
    Id                int64 `json:"id"`         
    User                int64 `json:"user"`         
    Attendance                int64 `json:"attendance"`         
    Health                int64 `json:"health"`         
    Exercisename                string `json:"exercisename"`         
    Sets                int `json:"sets"`         
    Reps                int `json:"reps"`         
    Weight                int `json:"weight"`         
    Duration                int `json:"duration"`         
    Calories                int `json:"calories"`         
    Note                string `json:"note"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type WorkoutlogManager struct {
    Conn    *Connection
    Result  *sql.Result
    Index   string
    Isolation   bool
    SelectQuery  string
    JoinQuery string
    CountQuery   string
    GroupQuery string
    SelectLog bool
    Log bool
}

func (c *Workoutlog) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewWorkoutlogManager(conn *Connection) *WorkoutlogManager {
    var item WorkoutlogManager


    if conn == nil {
        item.Conn = NewConnection()
        item.Isolation = false
    } else {
        item.Conn = conn 
        item.Isolation = conn.Isolation
    }

    item.Index = ""
    item.SelectLog = config.Log.Database
    item.Log = config.Log.Database

    return &item
}

func (p *WorkoutlogManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *WorkoutlogManager) SetIndex(index string) {
    p.Index = index
}

func (p *WorkoutlogManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *WorkoutlogManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *WorkoutlogManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *WorkoutlogManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
    if p.Isolation {
        query += " for update"
    }

    if p.SelectLog {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Query(query, params...)
}

func (p *WorkoutlogManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select wl_id, wl_user, wl_attendance, wl_health, wl_exercisename, wl_sets, wl_reps, wl_weight, wl_duration, wl_calories, wl_note, wl_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date, at_id, at_user, at_membership, at_gym, at_type, at_method, at_checkintime, at_checkouttime, at_duration, at_status, at_note, at_ip, at_device, at_createdby, at_date, h_id, h_category, h_term, h_name, h_count, h_cost, h_discount, h_costdiscount, h_content, h_date from workoutlog_tb, user_tb, attendance_tb, health_tb")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    if p.JoinQuery != "" {
        ret.WriteString(", ")
        ret.WriteString(p.JoinQuery)
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and wl_user = u_id ")
    
    ret.WriteString("and wl_attendance = at_id ")
    
    ret.WriteString("and wl_health = h_id ")
    

    return ret.String()
}

func (p *WorkoutlogManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from workoutlog_tb")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    if p.JoinQuery != "" {
        ret.WriteString(", ")
        ret.WriteString(p.JoinQuery)
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and wl_user = u_id ")
    
    ret.WriteString("and wl_attendance = at_id ")
    
    ret.WriteString("and wl_health = h_id ")
    

    return ret.String()
}

func (p *WorkoutlogManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select wl_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from workoutlog_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and wl_user = u_id ")
    
    ret.WriteString("and wl_attendance = at_id ")
    
    ret.WriteString("and wl_health = h_id ")
    

    return ret.String()
}

func (p *WorkoutlogManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate workoutlog_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *WorkoutlogManager) Insert(item *Workoutlog) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    if item.Date == "" {
        t := time.Now().UTC().Add(time.Hour * 9)
        //t := time.Now()
        item.Date = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    }

    
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

    query := ""
    var res sql.Result
    var err error
    if item.Id > 0 {
        query = "insert into workoutlog_tb (wl_id, wl_user, wl_attendance, wl_health, wl_exercisename, wl_sets, wl_reps, wl_weight, wl_duration, wl_calories, wl_note, wl_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.User, item.Attendance, item.Health, item.Exercisename, item.Sets, item.Reps, item.Weight, item.Duration, item.Calories, item.Note, item.Date)
    } else {
        query = "insert into workoutlog_tb (wl_user, wl_attendance, wl_health, wl_exercisename, wl_sets, wl_reps, wl_weight, wl_duration, wl_calories, wl_note, wl_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.User, item.Attendance, item.Health, item.Exercisename, item.Sets, item.Reps, item.Weight, item.Duration, item.Calories, item.Note, item.Date)
    }
    
    if err == nil {
        p.Result = &res
        
    } else {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
        p.Result = nil
    }

    return err
}

func (p *WorkoutlogManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from workoutlog_tb where wl_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *WorkoutlogManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from workoutlog_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *WorkoutlogManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
    var params []interface{}
    if initParams != nil {
        params = append(params, initParams...)
    }

    pos := 1

    var query strings.Builder
	query.WriteString(initQuery)

    for _, arg := range args {
        switch v := arg.(type) {        
        case Where:
            item := v

            if strings.Contains(item.Column, "_") {
                query.WriteString(" and ")
            } else {
                query.WriteString(" and wl_")
            }
            query.WriteString(item.Column)

            if item.Compare == "in" {
                query.WriteString(" in (")
                query.WriteString(strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]"))
                query.WriteString(")")
            } else if item.Compare == "not in" {
                query.WriteString(" not in (")
                query.WriteString(strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]"))
                query.WriteString(")")
            } else if item.Compare == "between" {
                if config.Database.Type == config.Postgresql {
                    query.WriteString(fmt.Sprintf(" between $%v and $%v", pos, pos + 1))
                    pos += 2
                } else {
                    query.WriteString(" between ? and ?")
                }

                s := item.Value.([2]string)
                params = append(params, s[0])
                params = append(params, s[1])
            } else {
                if config.Database.Type == config.Postgresql {
                    query.WriteString(" ")
                    query.WriteString(item.Compare)
                    query.WriteString(fmt.Sprintf(" $%v", pos))
                    pos++
                } else {
                    query.WriteString(" ")
                    query.WriteString(item.Compare)
                    query.WriteString(" ?")
                }
                if item.Compare == "like" {
                    params = append(params, "%" + item.Value.(string) + "%")
                } else {
                    params = append(params, item.Value)                
                }
            }
        case Custom:
             item := v

            query.WriteString(" and ")
            query.WriteString(item.Query)
        }        
    }

	query.WriteString(postQuery)

    return query.String(), params
}

func (p *WorkoutlogManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from workoutlog_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *WorkoutlogManager) Update(item *Workoutlog) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

	query := "update workoutlog_tb set wl_user = ?, wl_attendance = ?, wl_health = ?, wl_exercisename = ?, wl_sets = ?, wl_reps = ?, wl_weight = ?, wl_duration = ?, wl_calories = ?, wl_note = ?, wl_date = ? where wl_id = ?"
	_, err := p.Exec(query, item.User, item.Attendance, item.Health, item.Exercisename, item.Sets, item.Reps, item.Weight, item.Duration, item.Calories, item.Note, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *WorkoutlogManager) UpdateWhere(columns []workoutlog.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update workoutlog_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == workoutlog.ColumnId {
        initQuery.WriteString("wl_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == workoutlog.ColumnUser {
        initQuery.WriteString("wl_user = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == workoutlog.ColumnAttendance {
        initQuery.WriteString("wl_attendance = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == workoutlog.ColumnHealth {
        initQuery.WriteString("wl_health = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == workoutlog.ColumnExercisename {
        initQuery.WriteString("wl_exercisename = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == workoutlog.ColumnSets {
        initQuery.WriteString("wl_sets = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == workoutlog.ColumnReps {
        initQuery.WriteString("wl_reps = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == workoutlog.ColumnWeight {
        initQuery.WriteString("wl_weight = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == workoutlog.ColumnDuration {
        initQuery.WriteString("wl_duration = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == workoutlog.ColumnCalories {
        initQuery.WriteString("wl_calories = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == workoutlog.ColumnNote {
        initQuery.WriteString("wl_note = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == workoutlog.ColumnDate {
        initQuery.WriteString("wl_date = ?")
        initParams = append(initParams, v.Value)
        } else {
        
        }
    }

    initQuery.WriteString(" where 1=1 ")

    query, params := p.MakeQuery(initQuery.String(), "", initParams, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

/*


func (p *WorkoutlogManager) UpdateUser(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update workoutlog_tb set wl_user = ? where wl_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *WorkoutlogManager) UpdateAttendance(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update workoutlog_tb set wl_attendance = ? where wl_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *WorkoutlogManager) UpdateHealth(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update workoutlog_tb set wl_health = ? where wl_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *WorkoutlogManager) UpdateExercisename(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update workoutlog_tb set wl_exercisename = ? where wl_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *WorkoutlogManager) UpdateSets(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update workoutlog_tb set wl_sets = ? where wl_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *WorkoutlogManager) UpdateReps(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update workoutlog_tb set wl_reps = ? where wl_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *WorkoutlogManager) UpdateWeight(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update workoutlog_tb set wl_weight = ? where wl_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *WorkoutlogManager) UpdateDuration(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update workoutlog_tb set wl_duration = ? where wl_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *WorkoutlogManager) UpdateCalories(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update workoutlog_tb set wl_calories = ? where wl_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *WorkoutlogManager) UpdateNote(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update workoutlog_tb set wl_note = ? where wl_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *WorkoutlogManager) UpdateDate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update workoutlog_tb set wl_date = ? where wl_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *WorkoutlogManager) GetIdentity() int64 {
    if !p.Conn.IsConnect() {
        return 0
    }

    id, err := (*p.Result).LastInsertId()

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
        return 0
    } else {
        return id
    }
}

func (p *Workoutlog) InitExtra() {
    p.Extra = map[string]interface{}{

    }
}

func (p *WorkoutlogManager) ReadRow(rows *sql.Rows) *Workoutlog {
    var item Workoutlog
    var err error

    var _user User
    var _attendance Attendance
    var _health Health
    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.User, &item.Attendance, &item.Health, &item.Exercisename, &item.Sets, &item.Reps, &item.Weight, &item.Duration, &item.Calories, &item.Note, &item.Date, &_user.Id, &_user.Loginid, &_user.Passwd, &_user.Email, &_user.Name, &_user.Tel, &_user.Address, &_user.Image, &_user.Sex, &_user.Birth, &_user.Type, &_user.Connectid, &_user.Level, &_user.Role, &_user.Use, &_user.Logindate, &_user.Lastchangepasswddate, &_user.Date, &_attendance.Id, &_attendance.User, &_attendance.Membership, &_attendance.Gym, &_attendance.Type, &_attendance.Method, &_attendance.Checkintime, &_attendance.Checkouttime, &_attendance.Duration, &_attendance.Status, &_attendance.Note, &_attendance.Ip, &_attendance.Device, &_attendance.Createdby, &_attendance.Date, &_health.Id, &_health.Category, &_health.Term, &_health.Name, &_health.Count, &_health.Cost, &_health.Discount, &_health.Costdiscount, &_health.Content, &_health.Date)
        
        if item.Date == "0000-00-00 00:00:00" || item.Date == "1000-01-01 00:00:00" || item.Date == "9999-01-01 00:00:00" {
            item.Date = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Date = strings.ReplaceAll(strings.ReplaceAll(item.Date, "T", " "), "Z", "")
        }
		

    } else {
        return nil
    }

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
        return nil
    } else {

        item.InitExtra()
        _user.InitExtra()
        item.AddExtra("user",  _user)
_attendance.InitExtra()
        item.AddExtra("attendance",  _attendance)
_health.InitExtra()
        item.AddExtra("health",  _health)

        return &item
    }
}

func (p *WorkoutlogManager) ReadRows(rows *sql.Rows) []Workoutlog {
    var items []Workoutlog

    for rows.Next() {
        var item Workoutlog
        var _user User
        var _attendance Attendance
        var _health Health
        

        err := rows.Scan(&item.Id, &item.User, &item.Attendance, &item.Health, &item.Exercisename, &item.Sets, &item.Reps, &item.Weight, &item.Duration, &item.Calories, &item.Note, &item.Date, &_user.Id, &_user.Loginid, &_user.Passwd, &_user.Email, &_user.Name, &_user.Tel, &_user.Address, &_user.Image, &_user.Sex, &_user.Birth, &_user.Type, &_user.Connectid, &_user.Level, &_user.Role, &_user.Use, &_user.Logindate, &_user.Lastchangepasswddate, &_user.Date, &_attendance.Id, &_attendance.User, &_attendance.Membership, &_attendance.Gym, &_attendance.Type, &_attendance.Method, &_attendance.Checkintime, &_attendance.Checkouttime, &_attendance.Duration, &_attendance.Status, &_attendance.Note, &_attendance.Ip, &_attendance.Device, &_attendance.Createdby, &_attendance.Date, &_health.Id, &_health.Category, &_health.Term, &_health.Name, &_health.Count, &_health.Cost, &_health.Discount, &_health.Costdiscount, &_health.Content, &_health.Date)
        if err != nil {
           if p.Log {
             log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
        }

        
        if item.Date == "0000-00-00 00:00:00" || item.Date == "1000-01-01 00:00:00" || item.Date == "9999-01-01 00:00:00" {
            item.Date = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Date = strings.ReplaceAll(strings.ReplaceAll(item.Date, "T", " "), "Z", "")
        }
		

        item.InitExtra()
        _user.InitExtra()
        item.AddExtra("user",  _user)
_attendance.InitExtra()
        item.AddExtra("attendance",  _attendance)
_health.InitExtra()
        item.AddExtra("health",  _health)

        items = append(items, item)
    }


     return items
}

func (p *WorkoutlogManager) Get(id int64) *Workoutlog {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and wl_id = ?")

    
    query.WriteString(" and wl_user = u_id")
    
    query.WriteString(" and wl_attendance = at_id")
    
    query.WriteString(" and wl_health = h_id")
    
    
    rows, err := p.Query(query.String(), id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
       return nil
    }

    defer rows.Close()

    return p.ReadRow(rows)
}

func (p *WorkoutlogManager) GetWhere(args []interface{}) *Workoutlog {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *WorkoutlogManager) Count(args []interface{}) int {
    if !p.Conn.IsConnect() {
        return 0
    }

    query, params := p.MakeQuery(p.GetQuerySelect(), p.GroupQuery, nil, args)
    rows, err := p.Query(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
       return 0
    }

    defer rows.Close()

    if !rows.Next() {
        return 0
    }

    cnt := 0
    err = rows.Scan(&cnt)

    if err != nil {
        return 0
    } else {
        return cnt
    }
}

func (p *WorkoutlogManager) FindAll() []Workoutlog {
    return p.Find(nil)
}

func (p *WorkoutlogManager) Find(args []interface{}) []Workoutlog {
    if !p.Conn.IsConnect() {
        var items []Workoutlog
        return items
    }

    var params []interface{}
    baseQuery := p.GetQuery()

    var query strings.Builder

    page := 0
    pagesize := 0
    orderby := ""

    pos := 1
    
    for _, arg := range args {
        switch v := arg.(type) {
        case PagingType:
            item := v
            page = item.Page
            pagesize = item.Pagesize            
        case OrderingType:
            item := v
            orderby = item.Order
        case LimitType:
            item := v
            page = 1
            pagesize = item.Limit
        case OptionType:
            item := v
            if item.Limit > 0 {
                page = 1
                pagesize = item.Limit
            } else {
                page = item.Page
                pagesize = item.Pagesize                
            }
            orderby = item.Order
        case Where:
            item := v

            if strings.Contains(item.Column, "_") {
                query.WriteString(" and ")
            } else {
                query.WriteString(" and wl_")
            }
            query.WriteString(item.Column)
            
            if item.Compare == "in" {
                query.WriteString(" in (")
                query.WriteString(strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]"))
                query.WriteString(")")
            } else if item.Compare == "not in" {
                query.WriteString(" not in (")
                query.WriteString(strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]"))
                query.WriteString(")")
            } else if item.Compare == "between" {
                if config.Database.Type == config.Postgresql {
                    query.WriteString(fmt.Sprintf(" between $%v and $%v", pos, pos + 1))
                    pos += 2
                } else {
                    query.WriteString(" between ? and ?")
                }

                s := item.Value.([2]string)
                params = append(params, s[0])
                params = append(params, s[1])
            } else {
                if config.Database.Type == config.Postgresql {
                    query.WriteString(" ")
                    query.WriteString(item.Compare)
                    query.WriteString(fmt.Sprintf(" $%v", pos))
                    pos++
                } else {
                    query.WriteString(" ")
                    query.WriteString(item.Compare)
                    query.WriteString(" ?")
                }
                if item.Compare == "like" {
                    params = append(params, "%" + item.Value.(string) + "%")
                } else {
                    params = append(params, item.Value)                
                }
            }
        case Custom:
             item := v

            query.WriteString(" and ")
            query.WriteString(item.Query)
        case Base:
             item := v

             baseQuery = item.Query
        }
    }

    query.WriteString(p.GroupQuery)
    
    startpage := (page - 1) * pagesize
    
    if page > 0 && pagesize > 0 {
        if orderby == "" {
            orderby = "wl_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "wl_" + orderby
                }
            }
            
        }
        query.WriteString(" order by ")
        query.WriteString(orderby)
        if config.Database.Type == config.Postgresql {
            query.WriteString(fmt.Sprintf(" limit $%v offset $%v", pos, pos + 1))
            params = append(params, pagesize)
            params = append(params, startpage)
        } else if config.Database.Type == config.Mysql {
            query.WriteString(" limit ? offset ?")
            params = append(params, pagesize)
            params = append(params, startpage)
        } else if config.Database.Type == config.Sqlserver {
            query.WriteString("OFFSET ? ROWS FETCH NEXT ? ROWS ONLY")
            params = append(params, startpage)
            params = append(params, pagesize)
        }
    } else {
        if orderby == "" {
            orderby = "wl_id"
        } else {
            if !strings.Contains(orderby, "_") {
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "wl_" + orderby
                }
            }
        }
        query.WriteString(" order by ")
        query.WriteString(orderby)
    }

    rows, err := p.Query(baseQuery + query.String(), params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
        items := make([]Workoutlog, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}





func (p *WorkoutlogManager) GroupBy(name string, args []interface{}) []Groupby {
    if !p.Conn.IsConnect() {
        var items []Groupby
        return items
    }

    var params []interface{}
    baseQuery := p.GetQueryGroup(name)
    var query strings.Builder
    pos := 1

    for _, arg := range args {
        switch v := arg.(type) {
        case Where:
            item := v

            if strings.Contains(item.Column, "_") {
                query.WriteString(" and ")
            } else {
                query.WriteString(" and wl_")
            }
            query.WriteString(item.Column)
            
            if item.Compare == "in" {
                query.WriteString(" in (")
                query.WriteString(strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]"))
                query.WriteString(")")
            } else if item.Compare == "not in" {
                query.WriteString(" not in (")
                query.WriteString(strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]"))
                query.WriteString(")")
            } else if item.Compare == "between" {
                if config.Database.Type == config.Postgresql {
                    query.WriteString(fmt.Sprintf(" between $%v and $%v", pos, pos + 1))
                    pos += 2
                } else {
                    query.WriteString(" between ? and ?")
                }

                s := item.Value.([2]string)
                params = append(params, s[0])
                params = append(params, s[1])
            } else {
                if config.Database.Type == config.Postgresql {
                    query.WriteString(" ")
                    query.WriteString(item.Compare)
                    query.WriteString(fmt.Sprintf(" $%v", pos))
                    pos++
                } else {
                    query.WriteString(" ")
                    query.WriteString(item.Compare)
                    query.WriteString(" ?")
                }
                if item.Compare == "like" {
                    params = append(params, "%" + item.Value.(string) + "%")
                } else {
                    params = append(params, item.Value)                
                }
            }
        case Custom:
             item := v

            query.WriteString(" and ")
            query.WriteString(item.Query)
        case Base:
             item := v

             baseQuery = item.Query
        }
    }
    
    query.WriteString(" group by wl_")
    query.WriteString(name)

    rows, err := p.Query(baseQuery + query.String(), params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
        var items []Groupby
        return items
    }

    defer rows.Close()

    var items []Groupby

    for rows.Next() {
        var item Groupby
        err := rows.Scan(&item.Value, &item.Count)
        if err != nil {
           if p.Log {
                log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
        }

        items = append(items, item)
    }

    return items
}
