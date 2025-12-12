package models

import (
    "gym/global/config"
    "gym/models/gymtrainer"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Gymtrainer struct {
            
    Id                int64 `json:"id"`         
    Gym                int64 `json:"gym"`         
    Trainer                int64 `json:"trainer"`         
    Startdate                string `json:"startdate"`         
    Enddate                string `json:"enddate"`         
    Status                gymtrainer.Status `json:"status"`         
    Position                string `json:"position"`         
    Note                string `json:"note"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type GymtrainerManager struct {
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

func (c *Gymtrainer) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewGymtrainerManager(conn *Connection) *GymtrainerManager {
    var item GymtrainerManager


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

func (p *GymtrainerManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *GymtrainerManager) SetIndex(index string) {
    p.Index = index
}

func (p *GymtrainerManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *GymtrainerManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *GymtrainerManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *GymtrainerManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
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

func (p *GymtrainerManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select gt_id, gt_gym, gt_trainer, gt_startdate, gt_enddate, gt_status, gt_position, gt_note, gt_date, g_id, g_name, g_address, g_tel, g_user, g_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date from gymtrainer_tb, gym_tb, user_tb")

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
    
    ret.WriteString("and gt_gym = g_id ")
    
    ret.WriteString("and gt_trainer = u_id ")
    

    return ret.String()
}

func (p *GymtrainerManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from gymtrainer_tb")

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
    
    ret.WriteString("and gt_gym = g_id ")
    
    ret.WriteString("and gt_trainer = u_id ")
    

    return ret.String()
}

func (p *GymtrainerManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select gt_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from gymtrainer_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and gt_gym = g_id ")
    
    ret.WriteString("and gt_trainer = u_id ")
    

    return ret.String()
}

func (p *GymtrainerManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate gymtrainer_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *GymtrainerManager) Insert(item *Gymtrainer) error {
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
        query = "insert into gymtrainer_tb (gt_id, gt_gym, gt_trainer, gt_startdate, gt_enddate, gt_status, gt_position, gt_note, gt_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.Gym, item.Trainer, item.Startdate, item.Enddate, item.Status, item.Position, item.Note, item.Date)
    } else {
        query = "insert into gymtrainer_tb (gt_gym, gt_trainer, gt_startdate, gt_enddate, gt_status, gt_position, gt_note, gt_date) values (?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Gym, item.Trainer, item.Startdate, item.Enddate, item.Status, item.Position, item.Note, item.Date)
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

func (p *GymtrainerManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from gymtrainer_tb where gt_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *GymtrainerManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from gymtrainer_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *GymtrainerManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
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
                query.WriteString(" and gt_")
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

func (p *GymtrainerManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from gymtrainer_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *GymtrainerManager) Update(item *Gymtrainer) error {
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
	

	query := "update gymtrainer_tb set gt_gym = ?, gt_trainer = ?, gt_startdate = ?, gt_enddate = ?, gt_status = ?, gt_position = ?, gt_note = ?, gt_date = ? where gt_id = ?"
	_, err := p.Exec(query, item.Gym, item.Trainer, item.Startdate, item.Enddate, item.Status, item.Position, item.Note, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *GymtrainerManager) UpdateWhere(columns []gymtrainer.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update gymtrainer_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == gymtrainer.ColumnId {
        initQuery.WriteString("gt_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == gymtrainer.ColumnGym {
        initQuery.WriteString("gt_gym = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == gymtrainer.ColumnTrainer {
        initQuery.WriteString("gt_trainer = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == gymtrainer.ColumnStartdate {
        initQuery.WriteString("gt_startdate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == gymtrainer.ColumnEnddate {
        initQuery.WriteString("gt_enddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == gymtrainer.ColumnStatus {
        initQuery.WriteString("gt_status = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == gymtrainer.ColumnPosition {
        initQuery.WriteString("gt_position = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == gymtrainer.ColumnNote {
        initQuery.WriteString("gt_note = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == gymtrainer.ColumnDate {
        initQuery.WriteString("gt_date = ?")
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


func (p *GymtrainerManager) UpdateGym(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update gymtrainer_tb set gt_gym = ? where gt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *GymtrainerManager) UpdateTrainer(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update gymtrainer_tb set gt_trainer = ? where gt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *GymtrainerManager) UpdateStartdate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update gymtrainer_tb set gt_startdate = ? where gt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *GymtrainerManager) UpdateEnddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update gymtrainer_tb set gt_enddate = ? where gt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *GymtrainerManager) UpdateStatus(value gymtrainer.Status, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update gymtrainer_tb set gt_status = ? where gt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *GymtrainerManager) UpdatePosition(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update gymtrainer_tb set gt_position = ? where gt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *GymtrainerManager) UpdateNote(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update gymtrainer_tb set gt_note = ? where gt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *GymtrainerManager) UpdateDate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update gymtrainer_tb set gt_date = ? where gt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *GymtrainerManager) GetIdentity() int64 {
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

func (p *Gymtrainer) InitExtra() {
    p.Extra = map[string]interface{}{
            "status":     gymtrainer.GetStatus(p.Status),

    }
}

func (p *GymtrainerManager) ReadRow(rows *sql.Rows) *Gymtrainer {
    var item Gymtrainer
    var err error

    var _gym Gym
    var _traineruser User
    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.Gym, &item.Trainer, &item.Startdate, &item.Enddate, &item.Status, &item.Position, &item.Note, &item.Date, &_gym.Id, &_gym.Name, &_gym.Address, &_gym.Tel, &_gym.User, &_gym.Date, &_traineruser.Id, &_traineruser.Loginid, &_traineruser.Passwd, &_traineruser.Email, &_traineruser.Name, &_traineruser.Tel, &_traineruser.Address, &_traineruser.Image, &_traineruser.Sex, &_traineruser.Birth, &_traineruser.Type, &_traineruser.Connectid, &_traineruser.Level, &_traineruser.Role, &_traineruser.Use, &_traineruser.Logindate, &_traineruser.Lastchangepasswddate, &_traineruser.Date)
        
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
        _gym.InitExtra()
        item.AddExtra("gym",  _gym)
_traineruser.InitExtra()
        item.AddExtra("traineruser",  _traineruser)

        return &item
    }
}

func (p *GymtrainerManager) ReadRows(rows *sql.Rows) []Gymtrainer {
    var items []Gymtrainer

    for rows.Next() {
        var item Gymtrainer
        var _gym Gym
        var _traineruser User
        

        err := rows.Scan(&item.Id, &item.Gym, &item.Trainer, &item.Startdate, &item.Enddate, &item.Status, &item.Position, &item.Note, &item.Date, &_gym.Id, &_gym.Name, &_gym.Address, &_gym.Tel, &_gym.User, &_gym.Date, &_traineruser.Id, &_traineruser.Loginid, &_traineruser.Passwd, &_traineruser.Email, &_traineruser.Name, &_traineruser.Tel, &_traineruser.Address, &_traineruser.Image, &_traineruser.Sex, &_traineruser.Birth, &_traineruser.Type, &_traineruser.Connectid, &_traineruser.Level, &_traineruser.Role, &_traineruser.Use, &_traineruser.Logindate, &_traineruser.Lastchangepasswddate, &_traineruser.Date)
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
        _gym.InitExtra()
        item.AddExtra("gym",  _gym)
_traineruser.InitExtra()
        item.AddExtra("traineruser",  _traineruser)

        items = append(items, item)
    }


     return items
}

func (p *GymtrainerManager) Get(id int64) *Gymtrainer {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and gt_id = ?")

    
    query.WriteString(" and gt_gym = g_id")
    
    query.WriteString(" and gt_trainer = u_id")
    
    
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

func (p *GymtrainerManager) GetWhere(args []interface{}) *Gymtrainer {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *GymtrainerManager) Count(args []interface{}) int {
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

func (p *GymtrainerManager) FindAll() []Gymtrainer {
    return p.Find(nil)
}

func (p *GymtrainerManager) Find(args []interface{}) []Gymtrainer {
    if !p.Conn.IsConnect() {
        var items []Gymtrainer
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
                query.WriteString(" and gt_")
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
            orderby = "gt_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "gt_" + orderby
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
            orderby = "gt_id"
        } else {
            if !strings.Contains(orderby, "_") {
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "gt_" + orderby
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
        items := make([]Gymtrainer, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}





func (p *GymtrainerManager) GroupBy(name string, args []interface{}) []Groupby {
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
                query.WriteString(" and gt_")
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
    
    query.WriteString(" group by gt_")
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
