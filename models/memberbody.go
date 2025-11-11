package models

import (
    "gym/global/config"
    "gym/models/memberbody"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Memberbody struct {
            
    Id                int64 `json:"id"`         
    User                int64 `json:"user"`         
    Height                int `json:"height"`         
    Weight                int `json:"weight"`         
    Bodyfat                int `json:"bodyfat"`         
    Musclemass                int `json:"musclemass"`         
    Bmi                int `json:"bmi"`         
    Skeletalmuscle                int `json:"skeletalmuscle"`         
    Bodywater                int `json:"bodywater"`         
    Chest                int `json:"chest"`         
    Waist                int `json:"waist"`         
    Hip                int `json:"hip"`         
    Arm                int `json:"arm"`         
    Thigh                int `json:"thigh"`         
    Note                string `json:"note"`         
    Measureddate                string `json:"measureddate"`         
    Measuredby                int64 `json:"measuredby"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type MemberbodyManager struct {
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

func (c *Memberbody) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewMemberbodyManager(conn *Connection) *MemberbodyManager {
    var item MemberbodyManager


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

func (p *MemberbodyManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *MemberbodyManager) SetIndex(index string) {
    p.Index = index
}

func (p *MemberbodyManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *MemberbodyManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *MemberbodyManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *MemberbodyManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
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

func (p *MemberbodyManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select mb_id, mb_user, mb_height, mb_weight, mb_bodyfat, mb_musclemass, mb_bmi, mb_skeletalmuscle, mb_bodywater, mb_chest, mb_waist, mb_hip, mb_arm, mb_thigh, mb_note, mb_measureddate, mb_measuredby, mb_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date from memberbody_tb, user_tb, user_tb")

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
    
    ret.WriteString("and mb_user = u_id ")
    
    ret.WriteString("and mb_measuredby = u_id ")
    

    return ret.String()
}

func (p *MemberbodyManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from memberbody_tb")

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
    
    ret.WriteString("and mb_user = u_id ")
    
    ret.WriteString("and mb_measuredby = u_id ")
    

    return ret.String()
}

func (p *MemberbodyManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select mb_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from memberbody_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and mb_user = u_id ")
    
    ret.WriteString("and mb_measuredby = u_id ")
    

    return ret.String()
}

func (p *MemberbodyManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate memberbody_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *MemberbodyManager) Insert(item *Memberbody) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    if item.Date == "" {
        t := time.Now().UTC().Add(time.Hour * 9)
        //t := time.Now()
        item.Date = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    }

    
    if item.Measureddate == "" {
       item.Measureddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

    query := ""
    var res sql.Result
    var err error
    if item.Id > 0 {
        query = "insert into memberbody_tb (mb_id, mb_user, mb_height, mb_weight, mb_bodyfat, mb_musclemass, mb_bmi, mb_skeletalmuscle, mb_bodywater, mb_chest, mb_waist, mb_hip, mb_arm, mb_thigh, mb_note, mb_measureddate, mb_measuredby, mb_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.User, item.Height, item.Weight, item.Bodyfat, item.Musclemass, item.Bmi, item.Skeletalmuscle, item.Bodywater, item.Chest, item.Waist, item.Hip, item.Arm, item.Thigh, item.Note, item.Measureddate, item.Measuredby, item.Date)
    } else {
        query = "insert into memberbody_tb (mb_user, mb_height, mb_weight, mb_bodyfat, mb_musclemass, mb_bmi, mb_skeletalmuscle, mb_bodywater, mb_chest, mb_waist, mb_hip, mb_arm, mb_thigh, mb_note, mb_measureddate, mb_measuredby, mb_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.User, item.Height, item.Weight, item.Bodyfat, item.Musclemass, item.Bmi, item.Skeletalmuscle, item.Bodywater, item.Chest, item.Waist, item.Hip, item.Arm, item.Thigh, item.Note, item.Measureddate, item.Measuredby, item.Date)
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

func (p *MemberbodyManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from memberbody_tb where mb_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *MemberbodyManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from memberbody_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *MemberbodyManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
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
                query.WriteString(" and mb_")
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

func (p *MemberbodyManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from memberbody_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *MemberbodyManager) Update(item *Memberbody) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    
    if item.Measureddate == "" {
       item.Measureddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

	query := "update memberbody_tb set mb_user = ?, mb_height = ?, mb_weight = ?, mb_bodyfat = ?, mb_musclemass = ?, mb_bmi = ?, mb_skeletalmuscle = ?, mb_bodywater = ?, mb_chest = ?, mb_waist = ?, mb_hip = ?, mb_arm = ?, mb_thigh = ?, mb_note = ?, mb_measureddate = ?, mb_measuredby = ?, mb_date = ? where mb_id = ?"
	_, err := p.Exec(query, item.User, item.Height, item.Weight, item.Bodyfat, item.Musclemass, item.Bmi, item.Skeletalmuscle, item.Bodywater, item.Chest, item.Waist, item.Hip, item.Arm, item.Thigh, item.Note, item.Measureddate, item.Measuredby, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *MemberbodyManager) UpdateWhere(columns []memberbody.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update memberbody_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == memberbody.ColumnId {
        initQuery.WriteString("mb_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnUser {
        initQuery.WriteString("mb_user = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnHeight {
        initQuery.WriteString("mb_height = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnWeight {
        initQuery.WriteString("mb_weight = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnBodyfat {
        initQuery.WriteString("mb_bodyfat = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnMusclemass {
        initQuery.WriteString("mb_musclemass = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnBmi {
        initQuery.WriteString("mb_bmi = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnSkeletalmuscle {
        initQuery.WriteString("mb_skeletalmuscle = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnBodywater {
        initQuery.WriteString("mb_bodywater = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnChest {
        initQuery.WriteString("mb_chest = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnWaist {
        initQuery.WriteString("mb_waist = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnHip {
        initQuery.WriteString("mb_hip = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnArm {
        initQuery.WriteString("mb_arm = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnThigh {
        initQuery.WriteString("mb_thigh = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnNote {
        initQuery.WriteString("mb_note = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnMeasureddate {
        initQuery.WriteString("mb_measureddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnMeasuredby {
        initQuery.WriteString("mb_measuredby = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberbody.ColumnDate {
        initQuery.WriteString("mb_date = ?")
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


func (p *MemberbodyManager) UpdateUser(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_user = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberbodyManager) UpdateHeight(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_height = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberbodyManager) UpdateWeight(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_weight = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberbodyManager) UpdateBodyfat(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_bodyfat = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberbodyManager) UpdateMusclemass(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_musclemass = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberbodyManager) UpdateBmi(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_bmi = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberbodyManager) UpdateSkeletalmuscle(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_skeletalmuscle = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberbodyManager) UpdateBodywater(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_bodywater = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberbodyManager) UpdateChest(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_chest = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberbodyManager) UpdateWaist(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_waist = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberbodyManager) UpdateHip(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_hip = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberbodyManager) UpdateArm(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_arm = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberbodyManager) UpdateThigh(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_thigh = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberbodyManager) UpdateNote(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_note = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberbodyManager) UpdateMeasureddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_measureddate = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberbodyManager) UpdateMeasuredby(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_measuredby = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberbodyManager) UpdateDate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberbody_tb set mb_date = ? where mb_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *MemberbodyManager) GetIdentity() int64 {
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

func (p *Memberbody) InitExtra() {
    p.Extra = map[string]interface{}{

    }
}

func (p *MemberbodyManager) ReadRow(rows *sql.Rows) *Memberbody {
    var item Memberbody
    var err error

    var _memberuser User
    var _measuredbyuser User
    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.User, &item.Height, &item.Weight, &item.Bodyfat, &item.Musclemass, &item.Bmi, &item.Skeletalmuscle, &item.Bodywater, &item.Chest, &item.Waist, &item.Hip, &item.Arm, &item.Thigh, &item.Note, &item.Measureddate, &item.Measuredby, &item.Date, &_memberuser.Id, &_memberuser.Loginid, &_memberuser.Passwd, &_memberuser.Email, &_memberuser.Name, &_memberuser.Tel, &_memberuser.Address, &_memberuser.Image, &_memberuser.Sex, &_memberuser.Birth, &_memberuser.Type, &_memberuser.Connectid, &_memberuser.Level, &_memberuser.Role, &_memberuser.Use, &_memberuser.Logindate, &_memberuser.Lastchangepasswddate, &_memberuser.Date, &_measuredbyuser.Id, &_measuredbyuser.Loginid, &_measuredbyuser.Passwd, &_measuredbyuser.Email, &_measuredbyuser.Name, &_measuredbyuser.Tel, &_measuredbyuser.Address, &_measuredbyuser.Image, &_measuredbyuser.Sex, &_measuredbyuser.Birth, &_measuredbyuser.Type, &_measuredbyuser.Connectid, &_measuredbyuser.Level, &_measuredbyuser.Role, &_measuredbyuser.Use, &_measuredbyuser.Logindate, &_measuredbyuser.Lastchangepasswddate, &_measuredbyuser.Date)
        
        if item.Measureddate == "0000-00-00 00:00:00" || item.Measureddate == "1000-01-01 00:00:00" || item.Measureddate == "9999-01-01 00:00:00" {
            item.Measureddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Measureddate = strings.ReplaceAll(strings.ReplaceAll(item.Measureddate, "T", " "), "Z", "")
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
        _memberuser.InitExtra()
        item.AddExtra("memberuser",  _memberuser)
_measuredbyuser.InitExtra()
        item.AddExtra("measuredbyuser",  _measuredbyuser)

        return &item
    }
}

func (p *MemberbodyManager) ReadRows(rows *sql.Rows) []Memberbody {
    var items []Memberbody

    for rows.Next() {
        var item Memberbody
        var _memberuser User
        var _measuredbyuser User
        

        err := rows.Scan(&item.Id, &item.User, &item.Height, &item.Weight, &item.Bodyfat, &item.Musclemass, &item.Bmi, &item.Skeletalmuscle, &item.Bodywater, &item.Chest, &item.Waist, &item.Hip, &item.Arm, &item.Thigh, &item.Note, &item.Measureddate, &item.Measuredby, &item.Date, &_memberuser.Id, &_memberuser.Loginid, &_memberuser.Passwd, &_memberuser.Email, &_memberuser.Name, &_memberuser.Tel, &_memberuser.Address, &_memberuser.Image, &_memberuser.Sex, &_memberuser.Birth, &_memberuser.Type, &_memberuser.Connectid, &_memberuser.Level, &_memberuser.Role, &_memberuser.Use, &_memberuser.Logindate, &_memberuser.Lastchangepasswddate, &_memberuser.Date, &_measuredbyuser.Id, &_measuredbyuser.Loginid, &_measuredbyuser.Passwd, &_measuredbyuser.Email, &_measuredbyuser.Name, &_measuredbyuser.Tel, &_measuredbyuser.Address, &_measuredbyuser.Image, &_measuredbyuser.Sex, &_measuredbyuser.Birth, &_measuredbyuser.Type, &_measuredbyuser.Connectid, &_measuredbyuser.Level, &_measuredbyuser.Role, &_measuredbyuser.Use, &_measuredbyuser.Logindate, &_measuredbyuser.Lastchangepasswddate, &_measuredbyuser.Date)
        if err != nil {
           if p.Log {
             log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
        }

        
        if item.Measureddate == "0000-00-00 00:00:00" || item.Measureddate == "1000-01-01 00:00:00" || item.Measureddate == "9999-01-01 00:00:00" {
            item.Measureddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Measureddate = strings.ReplaceAll(strings.ReplaceAll(item.Measureddate, "T", " "), "Z", "")
        }
		
        if item.Date == "0000-00-00 00:00:00" || item.Date == "1000-01-01 00:00:00" || item.Date == "9999-01-01 00:00:00" {
            item.Date = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Date = strings.ReplaceAll(strings.ReplaceAll(item.Date, "T", " "), "Z", "")
        }
		

        item.InitExtra()
        _memberuser.InitExtra()
        item.AddExtra("memberuser",  _memberuser)
_measuredbyuser.InitExtra()
        item.AddExtra("measuredbyuser",  _measuredbyuser)

        items = append(items, item)
    }


     return items
}

func (p *MemberbodyManager) Get(id int64) *Memberbody {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and mb_id = ?")

    
    query.WriteString(" and mb_user = u_id")
    
    query.WriteString(" and mb_measuredby = u_id")
    
    
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

func (p *MemberbodyManager) GetWhere(args []interface{}) *Memberbody {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *MemberbodyManager) Count(args []interface{}) int {
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

func (p *MemberbodyManager) FindAll() []Memberbody {
    return p.Find(nil)
}

func (p *MemberbodyManager) Find(args []interface{}) []Memberbody {
    if !p.Conn.IsConnect() {
        var items []Memberbody
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
                query.WriteString(" and mb_")
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
            orderby = "mb_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "mb_" + orderby
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
            orderby = "mb_id"
        } else {
            if !strings.Contains(orderby, "_") {
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "mb_" + orderby
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
        items := make([]Memberbody, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}





func (p *MemberbodyManager) GroupBy(name string, args []interface{}) []Groupby {
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
                query.WriteString(" and mb_")
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
    
    query.WriteString(" group by mb_")
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
