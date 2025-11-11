package models

import (
    "gym/global/config"
    "gym/models/ptreservation"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Ptreservation struct {
            
    Id                int64 `json:"id"`         
    Trainer                int64 `json:"trainer"`         
    Member                int64 `json:"member"`         
    Gym                int64 `json:"gym"`         
    Reservationdate                string `json:"reservationdate"`         
    Starttime                string `json:"starttime"`         
    Endtime                string `json:"endtime"`         
    Duration                int `json:"duration"`         
    Status                ptreservation.Status `json:"status"`         
    Note                string `json:"note"`         
    Cancelreason                string `json:"cancelreason"`         
    Createddate                string `json:"createddate"`         
    Updateddate                string `json:"updateddate"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type PtreservationManager struct {
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

func (c *Ptreservation) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewPtreservationManager(conn *Connection) *PtreservationManager {
    var item PtreservationManager


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

func (p *PtreservationManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *PtreservationManager) SetIndex(index string) {
    p.Index = index
}

func (p *PtreservationManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *PtreservationManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *PtreservationManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *PtreservationManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
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

func (p *PtreservationManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select pr_id, pr_trainer, pr_member, pr_gym, pr_reservationdate, pr_starttime, pr_endtime, pr_duration, pr_status, pr_note, pr_cancelreason, pr_createddate, pr_updateddate, pr_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date, g_id, g_name, g_date from ptreservation_tb, user_tb, user_tb, gym_tb")

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
    
    ret.WriteString("and pr_trainer = u_id ")
    
    ret.WriteString("and pr_member = u_id ")
    
    ret.WriteString("and pr_gym = g_id ")
    

    return ret.String()
}

func (p *PtreservationManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from ptreservation_tb")

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
    
    ret.WriteString("and pr_trainer = u_id ")
    
    ret.WriteString("and pr_member = u_id ")
    
    ret.WriteString("and pr_gym = g_id ")
    

    return ret.String()
}

func (p *PtreservationManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select pr_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from ptreservation_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and pr_trainer = u_id ")
    
    ret.WriteString("and pr_member = u_id ")
    
    ret.WriteString("and pr_gym = g_id ")
    

    return ret.String()
}

func (p *PtreservationManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate ptreservation_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *PtreservationManager) Insert(item *Ptreservation) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    if item.Date == "" {
        t := time.Now().UTC().Add(time.Hour * 9)
        //t := time.Now()
        item.Date = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    }

    
    if item.Reservationdate == "" {
       item.Reservationdate = "1000-01-01 00:00:00"
    }
	
    if item.Createddate == "" {
       item.Createddate = "1000-01-01 00:00:00"
    }
	
    if item.Updateddate == "" {
       item.Updateddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

    query := ""
    var res sql.Result
    var err error
    if item.Id > 0 {
        query = "insert into ptreservation_tb (pr_id, pr_trainer, pr_member, pr_gym, pr_reservationdate, pr_starttime, pr_endtime, pr_duration, pr_status, pr_note, pr_cancelreason, pr_createddate, pr_updateddate, pr_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.Trainer, item.Member, item.Gym, item.Reservationdate, item.Starttime, item.Endtime, item.Duration, item.Status, item.Note, item.Cancelreason, item.Createddate, item.Updateddate, item.Date)
    } else {
        query = "insert into ptreservation_tb (pr_trainer, pr_member, pr_gym, pr_reservationdate, pr_starttime, pr_endtime, pr_duration, pr_status, pr_note, pr_cancelreason, pr_createddate, pr_updateddate, pr_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Trainer, item.Member, item.Gym, item.Reservationdate, item.Starttime, item.Endtime, item.Duration, item.Status, item.Note, item.Cancelreason, item.Createddate, item.Updateddate, item.Date)
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

func (p *PtreservationManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from ptreservation_tb where pr_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *PtreservationManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from ptreservation_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *PtreservationManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
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
                query.WriteString(" and pr_")
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

func (p *PtreservationManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from ptreservation_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *PtreservationManager) Update(item *Ptreservation) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    
    if item.Reservationdate == "" {
       item.Reservationdate = "1000-01-01 00:00:00"
    }
	
    if item.Createddate == "" {
       item.Createddate = "1000-01-01 00:00:00"
    }
	
    if item.Updateddate == "" {
       item.Updateddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

	query := "update ptreservation_tb set pr_trainer = ?, pr_member = ?, pr_gym = ?, pr_reservationdate = ?, pr_starttime = ?, pr_endtime = ?, pr_duration = ?, pr_status = ?, pr_note = ?, pr_cancelreason = ?, pr_createddate = ?, pr_updateddate = ?, pr_date = ? where pr_id = ?"
	_, err := p.Exec(query, item.Trainer, item.Member, item.Gym, item.Reservationdate, item.Starttime, item.Endtime, item.Duration, item.Status, item.Note, item.Cancelreason, item.Createddate, item.Updateddate, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *PtreservationManager) UpdateWhere(columns []ptreservation.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update ptreservation_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == ptreservation.ColumnId {
        initQuery.WriteString("pr_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == ptreservation.ColumnTrainer {
        initQuery.WriteString("pr_trainer = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == ptreservation.ColumnMember {
        initQuery.WriteString("pr_member = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == ptreservation.ColumnGym {
        initQuery.WriteString("pr_gym = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == ptreservation.ColumnReservationdate {
        initQuery.WriteString("pr_reservationdate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == ptreservation.ColumnStarttime {
        initQuery.WriteString("pr_starttime = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == ptreservation.ColumnEndtime {
        initQuery.WriteString("pr_endtime = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == ptreservation.ColumnDuration {
        initQuery.WriteString("pr_duration = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == ptreservation.ColumnStatus {
        initQuery.WriteString("pr_status = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == ptreservation.ColumnNote {
        initQuery.WriteString("pr_note = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == ptreservation.ColumnCancelreason {
        initQuery.WriteString("pr_cancelreason = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == ptreservation.ColumnCreateddate {
        initQuery.WriteString("pr_createddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == ptreservation.ColumnUpdateddate {
        initQuery.WriteString("pr_updateddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == ptreservation.ColumnDate {
        initQuery.WriteString("pr_date = ?")
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


func (p *PtreservationManager) UpdateTrainer(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update ptreservation_tb set pr_trainer = ? where pr_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *PtreservationManager) UpdateMember(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update ptreservation_tb set pr_member = ? where pr_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *PtreservationManager) UpdateGym(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update ptreservation_tb set pr_gym = ? where pr_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *PtreservationManager) UpdateReservationdate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update ptreservation_tb set pr_reservationdate = ? where pr_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *PtreservationManager) UpdateStarttime(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update ptreservation_tb set pr_starttime = ? where pr_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *PtreservationManager) UpdateEndtime(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update ptreservation_tb set pr_endtime = ? where pr_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *PtreservationManager) UpdateDuration(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update ptreservation_tb set pr_duration = ? where pr_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *PtreservationManager) UpdateStatus(value ptreservation.Status, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update ptreservation_tb set pr_status = ? where pr_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *PtreservationManager) UpdateNote(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update ptreservation_tb set pr_note = ? where pr_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *PtreservationManager) UpdateCancelreason(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update ptreservation_tb set pr_cancelreason = ? where pr_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *PtreservationManager) UpdateCreateddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update ptreservation_tb set pr_createddate = ? where pr_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *PtreservationManager) UpdateUpdateddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update ptreservation_tb set pr_updateddate = ? where pr_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *PtreservationManager) UpdateDate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update ptreservation_tb set pr_date = ? where pr_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *PtreservationManager) GetIdentity() int64 {
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

func (p *Ptreservation) InitExtra() {
    p.Extra = map[string]interface{}{
            "status":     ptreservation.GetStatus(p.Status),

    }
}

func (p *PtreservationManager) ReadRow(rows *sql.Rows) *Ptreservation {
    var item Ptreservation
    var err error

    var _traineruser User
    var _memberuser User
    var _gym Gym
    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.Trainer, &item.Member, &item.Gym, &item.Reservationdate, &item.Starttime, &item.Endtime, &item.Duration, &item.Status, &item.Note, &item.Cancelreason, &item.Createddate, &item.Updateddate, &item.Date, &_traineruser.Id, &_traineruser.Loginid, &_traineruser.Passwd, &_traineruser.Email, &_traineruser.Name, &_traineruser.Tel, &_traineruser.Address, &_traineruser.Image, &_traineruser.Sex, &_traineruser.Birth, &_traineruser.Type, &_traineruser.Connectid, &_traineruser.Level, &_traineruser.Role, &_traineruser.Use, &_traineruser.Logindate, &_traineruser.Lastchangepasswddate, &_traineruser.Date, &_memberuser.Id, &_memberuser.Loginid, &_memberuser.Passwd, &_memberuser.Email, &_memberuser.Name, &_memberuser.Tel, &_memberuser.Address, &_memberuser.Image, &_memberuser.Sex, &_memberuser.Birth, &_memberuser.Type, &_memberuser.Connectid, &_memberuser.Level, &_memberuser.Role, &_memberuser.Use, &_memberuser.Logindate, &_memberuser.Lastchangepasswddate, &_memberuser.Date, &_gym.Id, &_gym.Name, &_gym.Date)
        
        if item.Reservationdate == "0000-00-00 00:00:00" || item.Reservationdate == "1000-01-01 00:00:00" || item.Reservationdate == "9999-01-01 00:00:00" {
            item.Reservationdate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Reservationdate = strings.ReplaceAll(strings.ReplaceAll(item.Reservationdate, "T", " "), "Z", "")
        }
		
        if item.Createddate == "0000-00-00 00:00:00" || item.Createddate == "1000-01-01 00:00:00" || item.Createddate == "9999-01-01 00:00:00" {
            item.Createddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Createddate = strings.ReplaceAll(strings.ReplaceAll(item.Createddate, "T", " "), "Z", "")
        }
		
        if item.Updateddate == "0000-00-00 00:00:00" || item.Updateddate == "1000-01-01 00:00:00" || item.Updateddate == "9999-01-01 00:00:00" {
            item.Updateddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Updateddate = strings.ReplaceAll(strings.ReplaceAll(item.Updateddate, "T", " "), "Z", "")
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

func (p *PtreservationManager) ReadRows(rows *sql.Rows) []Ptreservation {
    var items []Ptreservation

    for rows.Next() {
        var item Ptreservation
        var _traineruser User
        var _memberuser User
        var _gym Gym
        

        err := rows.Scan(&item.Id, &item.Trainer, &item.Member, &item.Gym, &item.Reservationdate, &item.Starttime, &item.Endtime, &item.Duration, &item.Status, &item.Note, &item.Cancelreason, &item.Createddate, &item.Updateddate, &item.Date, &_traineruser.Id, &_traineruser.Loginid, &_traineruser.Passwd, &_traineruser.Email, &_traineruser.Name, &_traineruser.Tel, &_traineruser.Address, &_traineruser.Image, &_traineruser.Sex, &_traineruser.Birth, &_traineruser.Type, &_traineruser.Connectid, &_traineruser.Level, &_traineruser.Role, &_traineruser.Use, &_traineruser.Logindate, &_traineruser.Lastchangepasswddate, &_traineruser.Date, &_memberuser.Id, &_memberuser.Loginid, &_memberuser.Passwd, &_memberuser.Email, &_memberuser.Name, &_memberuser.Tel, &_memberuser.Address, &_memberuser.Image, &_memberuser.Sex, &_memberuser.Birth, &_memberuser.Type, &_memberuser.Connectid, &_memberuser.Level, &_memberuser.Role, &_memberuser.Use, &_memberuser.Logindate, &_memberuser.Lastchangepasswddate, &_memberuser.Date, &_gym.Id, &_gym.Name, &_gym.Date)
        if err != nil {
           if p.Log {
             log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
        }

        
        if item.Reservationdate == "0000-00-00 00:00:00" || item.Reservationdate == "1000-01-01 00:00:00" || item.Reservationdate == "9999-01-01 00:00:00" {
            item.Reservationdate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Reservationdate = strings.ReplaceAll(strings.ReplaceAll(item.Reservationdate, "T", " "), "Z", "")
        }
		
        if item.Createddate == "0000-00-00 00:00:00" || item.Createddate == "1000-01-01 00:00:00" || item.Createddate == "9999-01-01 00:00:00" {
            item.Createddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Createddate = strings.ReplaceAll(strings.ReplaceAll(item.Createddate, "T", " "), "Z", "")
        }
		
        if item.Updateddate == "0000-00-00 00:00:00" || item.Updateddate == "1000-01-01 00:00:00" || item.Updateddate == "9999-01-01 00:00:00" {
            item.Updateddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Updateddate = strings.ReplaceAll(strings.ReplaceAll(item.Updateddate, "T", " "), "Z", "")
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

func (p *PtreservationManager) Get(id int64) *Ptreservation {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and pr_id = ?")

    
    query.WriteString(" and pr_trainer = u_id")
    
    query.WriteString(" and pr_member = u_id")
    
    query.WriteString(" and pr_gym = g_id")
    
    
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

func (p *PtreservationManager) GetWhere(args []interface{}) *Ptreservation {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *PtreservationManager) Count(args []interface{}) int {
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

func (p *PtreservationManager) FindAll() []Ptreservation {
    return p.Find(nil)
}

func (p *PtreservationManager) Find(args []interface{}) []Ptreservation {
    if !p.Conn.IsConnect() {
        var items []Ptreservation
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
                query.WriteString(" and pr_")
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
            orderby = "pr_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "pr_" + orderby
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
            orderby = "pr_id"
        } else {
            if !strings.Contains(orderby, "_") {
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "pr_" + orderby
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
        items := make([]Ptreservation, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}





func (p *PtreservationManager) GroupBy(name string, args []interface{}) []Groupby {
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
                query.WriteString(" and pr_")
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
    
    query.WriteString(" group by pr_")
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
