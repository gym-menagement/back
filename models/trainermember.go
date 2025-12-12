package models

import (
    "gym/global/config"
    "gym/models/trainermember"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Trainermember struct {
            
    Id                int64 `json:"id"`         
    Trainer                int64 `json:"trainer"`         
    Member                int64 `json:"member"`         
    Gym                int64 `json:"gym"`         
    Startdate                string `json:"startdate"`         
    Enddate                string `json:"enddate"`         
    Status                trainermember.Status `json:"status"`         
    Note                string `json:"note"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type TrainermemberManager struct {
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

func (c *Trainermember) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewTrainermemberManager(conn *Connection) *TrainermemberManager {
    var item TrainermemberManager


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

func (p *TrainermemberManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *TrainermemberManager) SetIndex(index string) {
    p.Index = index
}

func (p *TrainermemberManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *TrainermemberManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *TrainermemberManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *TrainermemberManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
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

func (p *TrainermemberManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select tm_id, tm_trainer, tm_member, tm_gym, tm_startdate, tm_enddate, tm_status, tm_note, tm_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date, g_id, g_name, g_address, g_tel, g_user, g_date from trainermember_tb, user_tb, user_tb, gym_tb")

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
    
    ret.WriteString("and tm_trainer = u_id ")
    
    ret.WriteString("and tm_member = u_id ")
    
    ret.WriteString("and tm_gym = g_id ")
    

    return ret.String()
}

func (p *TrainermemberManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from trainermember_tb")

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
    
    ret.WriteString("and tm_trainer = u_id ")
    
    ret.WriteString("and tm_member = u_id ")
    
    ret.WriteString("and tm_gym = g_id ")
    

    return ret.String()
}

func (p *TrainermemberManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select tm_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from trainermember_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and tm_trainer = u_id ")
    
    ret.WriteString("and tm_member = u_id ")
    
    ret.WriteString("and tm_gym = g_id ")
    

    return ret.String()
}

func (p *TrainermemberManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate trainermember_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *TrainermemberManager) Insert(item *Trainermember) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    if item.Date == "" {
        t := time.Now().UTC().Add(time.Hour * 9)
        //t := time.Now()
        item.Date = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    }

    
    if item.Startdate == "" {
       item.Startdate = "1000-01-01 00:00:00"
    }
	
    if item.Enddate == "" {
       item.Enddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

    query := ""
    var res sql.Result
    var err error
    if item.Id > 0 {
        query = "insert into trainermember_tb (tm_id, tm_trainer, tm_member, tm_gym, tm_startdate, tm_enddate, tm_status, tm_note, tm_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.Trainer, item.Member, item.Gym, item.Startdate, item.Enddate, item.Status, item.Note, item.Date)
    } else {
        query = "insert into trainermember_tb (tm_trainer, tm_member, tm_gym, tm_startdate, tm_enddate, tm_status, tm_note, tm_date) values (?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Trainer, item.Member, item.Gym, item.Startdate, item.Enddate, item.Status, item.Note, item.Date)
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

func (p *TrainermemberManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from trainermember_tb where tm_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *TrainermemberManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from trainermember_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *TrainermemberManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
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
                query.WriteString(" and tm_")
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

func (p *TrainermemberManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from trainermember_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *TrainermemberManager) Update(item *Trainermember) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    
    if item.Startdate == "" {
       item.Startdate = "1000-01-01 00:00:00"
    }
	
    if item.Enddate == "" {
       item.Enddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

	query := "update trainermember_tb set tm_trainer = ?, tm_member = ?, tm_gym = ?, tm_startdate = ?, tm_enddate = ?, tm_status = ?, tm_note = ?, tm_date = ? where tm_id = ?"
	_, err := p.Exec(query, item.Trainer, item.Member, item.Gym, item.Startdate, item.Enddate, item.Status, item.Note, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *TrainermemberManager) UpdateWhere(columns []trainermember.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update trainermember_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == trainermember.ColumnId {
        initQuery.WriteString("tm_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == trainermember.ColumnTrainer {
        initQuery.WriteString("tm_trainer = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == trainermember.ColumnMember {
        initQuery.WriteString("tm_member = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == trainermember.ColumnGym {
        initQuery.WriteString("tm_gym = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == trainermember.ColumnStartdate {
        initQuery.WriteString("tm_startdate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == trainermember.ColumnEnddate {
        initQuery.WriteString("tm_enddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == trainermember.ColumnStatus {
        initQuery.WriteString("tm_status = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == trainermember.ColumnNote {
        initQuery.WriteString("tm_note = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == trainermember.ColumnDate {
        initQuery.WriteString("tm_date = ?")
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


func (p *TrainermemberManager) UpdateTrainer(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update trainermember_tb set tm_trainer = ? where tm_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *TrainermemberManager) UpdateMember(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update trainermember_tb set tm_member = ? where tm_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *TrainermemberManager) UpdateGym(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update trainermember_tb set tm_gym = ? where tm_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *TrainermemberManager) UpdateStartdate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update trainermember_tb set tm_startdate = ? where tm_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *TrainermemberManager) UpdateEnddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update trainermember_tb set tm_enddate = ? where tm_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *TrainermemberManager) UpdateStatus(value trainermember.Status, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update trainermember_tb set tm_status = ? where tm_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *TrainermemberManager) UpdateNote(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update trainermember_tb set tm_note = ? where tm_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *TrainermemberManager) UpdateDate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update trainermember_tb set tm_date = ? where tm_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *TrainermemberManager) GetIdentity() int64 {
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

func (p *Trainermember) InitExtra() {
    p.Extra = map[string]interface{}{
            "status":     trainermember.GetStatus(p.Status),

    }
}

func (p *TrainermemberManager) ReadRow(rows *sql.Rows) *Trainermember {
    var item Trainermember
    var err error

    var _traineruser User
    var _memberuser User
    var _gym Gym
    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.Trainer, &item.Member, &item.Gym, &item.Startdate, &item.Enddate, &item.Status, &item.Note, &item.Date, &_traineruser.Id, &_traineruser.Loginid, &_traineruser.Passwd, &_traineruser.Email, &_traineruser.Name, &_traineruser.Tel, &_traineruser.Address, &_traineruser.Image, &_traineruser.Sex, &_traineruser.Birth, &_traineruser.Type, &_traineruser.Connectid, &_traineruser.Level, &_traineruser.Role, &_traineruser.Use, &_traineruser.Logindate, &_traineruser.Lastchangepasswddate, &_traineruser.Date, &_memberuser.Id, &_memberuser.Loginid, &_memberuser.Passwd, &_memberuser.Email, &_memberuser.Name, &_memberuser.Tel, &_memberuser.Address, &_memberuser.Image, &_memberuser.Sex, &_memberuser.Birth, &_memberuser.Type, &_memberuser.Connectid, &_memberuser.Level, &_memberuser.Role, &_memberuser.Use, &_memberuser.Logindate, &_memberuser.Lastchangepasswddate, &_memberuser.Date, &_gym.Id, &_gym.Name, &_gym.Address, &_gym.Tel, &_gym.User, &_gym.Date)
        
        if item.Startdate == "0000-00-00 00:00:00" || item.Startdate == "1000-01-01 00:00:00" || item.Startdate == "9999-01-01 00:00:00" {
            item.Startdate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Startdate = strings.ReplaceAll(strings.ReplaceAll(item.Startdate, "T", " "), "Z", "")
        }
		
        if item.Enddate == "0000-00-00 00:00:00" || item.Enddate == "1000-01-01 00:00:00" || item.Enddate == "9999-01-01 00:00:00" {
            item.Enddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Enddate = strings.ReplaceAll(strings.ReplaceAll(item.Enddate, "T", " "), "Z", "")
        }
		
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
        _traineruser.InitExtra()
        item.AddExtra("traineruser",  _traineruser)
_memberuser.InitExtra()
        item.AddExtra("memberuser",  _memberuser)
_gym.InitExtra()
        item.AddExtra("gym",  _gym)

        return &item
    }
}

func (p *TrainermemberManager) ReadRows(rows *sql.Rows) []Trainermember {
    var items []Trainermember

    for rows.Next() {
        var item Trainermember
        var _traineruser User
        var _memberuser User
        var _gym Gym
        

        err := rows.Scan(&item.Id, &item.Trainer, &item.Member, &item.Gym, &item.Startdate, &item.Enddate, &item.Status, &item.Note, &item.Date, &_traineruser.Id, &_traineruser.Loginid, &_traineruser.Passwd, &_traineruser.Email, &_traineruser.Name, &_traineruser.Tel, &_traineruser.Address, &_traineruser.Image, &_traineruser.Sex, &_traineruser.Birth, &_traineruser.Type, &_traineruser.Connectid, &_traineruser.Level, &_traineruser.Role, &_traineruser.Use, &_traineruser.Logindate, &_traineruser.Lastchangepasswddate, &_traineruser.Date, &_memberuser.Id, &_memberuser.Loginid, &_memberuser.Passwd, &_memberuser.Email, &_memberuser.Name, &_memberuser.Tel, &_memberuser.Address, &_memberuser.Image, &_memberuser.Sex, &_memberuser.Birth, &_memberuser.Type, &_memberuser.Connectid, &_memberuser.Level, &_memberuser.Role, &_memberuser.Use, &_memberuser.Logindate, &_memberuser.Lastchangepasswddate, &_memberuser.Date, &_gym.Id, &_gym.Name, &_gym.Address, &_gym.Tel, &_gym.User, &_gym.Date)
        if err != nil {
           if p.Log {
             log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
        }

        
        if item.Startdate == "0000-00-00 00:00:00" || item.Startdate == "1000-01-01 00:00:00" || item.Startdate == "9999-01-01 00:00:00" {
            item.Startdate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Startdate = strings.ReplaceAll(strings.ReplaceAll(item.Startdate, "T", " "), "Z", "")
        }
		
        if item.Enddate == "0000-00-00 00:00:00" || item.Enddate == "1000-01-01 00:00:00" || item.Enddate == "9999-01-01 00:00:00" {
            item.Enddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Enddate = strings.ReplaceAll(strings.ReplaceAll(item.Enddate, "T", " "), "Z", "")
        }
		
        if item.Date == "0000-00-00 00:00:00" || item.Date == "1000-01-01 00:00:00" || item.Date == "9999-01-01 00:00:00" {
            item.Date = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Date = strings.ReplaceAll(strings.ReplaceAll(item.Date, "T", " "), "Z", "")
        }
		

        item.InitExtra()
        _traineruser.InitExtra()
        item.AddExtra("traineruser",  _traineruser)
_memberuser.InitExtra()
        item.AddExtra("memberuser",  _memberuser)
_gym.InitExtra()
        item.AddExtra("gym",  _gym)

        items = append(items, item)
    }


     return items
}

func (p *TrainermemberManager) Get(id int64) *Trainermember {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and tm_id = ?")

    
    query.WriteString(" and tm_trainer = u_id")
    
    query.WriteString(" and tm_member = u_id")
    
    query.WriteString(" and tm_gym = g_id")
    
    
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

func (p *TrainermemberManager) GetWhere(args []interface{}) *Trainermember {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *TrainermemberManager) Count(args []interface{}) int {
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

func (p *TrainermemberManager) FindAll() []Trainermember {
    return p.Find(nil)
}

func (p *TrainermemberManager) Find(args []interface{}) []Trainermember {
    if !p.Conn.IsConnect() {
        var items []Trainermember
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
                query.WriteString(" and tm_")
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
            orderby = "tm_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "tm_" + orderby
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
            orderby = "tm_id"
        } else {
            if !strings.Contains(orderby, "_") {
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "tm_" + orderby
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
        items := make([]Trainermember, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}





func (p *TrainermemberManager) GroupBy(name string, args []interface{}) []Groupby {
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
                query.WriteString(" and tm_")
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
    
    query.WriteString(" group by tm_")
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
