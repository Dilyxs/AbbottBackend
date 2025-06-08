package utils

import (
	"database/sql"
	"time"
)

type Client struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Message     string `json:"message"`
	CallWeek    string `json:"callweek"`
	HighValue   bool   `json:"highvalue"`
	BookingDone bool   `json:"bookingdone"`
	BookingDate string `json:"bookingdate"` // use time.Time if you parse it
	Signed      bool   `json:"signed"`
	SignedPrice string `json:"signedprice"`
	WorkDate    string `json:"workdate"` // use time.Time if you parse it
	LeadUser    int    `json:"leaduser"`
}

type RegisterIngo struct {
	Email    string `json:"email"`
	Password string `json:password`
}

type ClientAllInfo struct {
	Id          int            `json:"id"`
	Name        sql.NullString `json:"name"`
	Phone       sql.NullString `json:"phone"`
	Address     sql.NullString `json:"address"`
	Message     sql.NullString `json:"message"`
	CallWeek    sql.NullTime   `json:"callweek"`
	HighValue   sql.NullBool   `json:"highvalue"`
	BookingDone sql.NullBool   `json:"bookingdone"`
	BookingDate sql.NullTime   `json:"bookingdate"`
	Signed      sql.NullBool   `json:"signed"`
	SignedPrice sql.NullString `json:"signedprice"`
	WorkDate    sql.NullTime   `json:"workdate"`
	LeadUser    sql.NullInt64  `json:"leaduser"`
}

type MapsData struct {
	Lat  string `json:"Lat"`
	Long string `json:"Long"`
}
type UserInformation struct {
	Id      int    `json:"id"`
	Isadmin bool   `json:"isadmin"`
	Email   string `json:"email"`
	Cut     int    `json:"cut"`
}

type SendUserID struct {
	Id int `json:"id"`
}

type CalenderInfo struct {
	Id             int       `json:"id"`
	StartTime      time.Time `json:"starttime"`
	EndTime        time.Time `json:"endtime"`
	Message        string    `json:"message"`
	ClientRelation int       `json:"clientrelation"`
}
type CalenderInfoReceived struct {
	StartTime      string `json:"starttime"`
	EndTime        string `json:"endtime"`
	Message        string `json:"message"`
	ClientRelation int    `json:"clientrelation"`
}

type NotesDetail struct {
	Id             int           `json:"id"`
	Title          string        `json:"title"`
	Detail         string        `json:"detail"`
	ClientRelation sql.NullInt64 `json:"clientrelation"`
}
type NotesDetailSender struct {
	Title          string `json:"title"`
	Detail         string `json:"detail"`
	ClientRelation int    `json:"clientrelation"`
}

type NotesDetailUpdater struct {
	Id             int    `json:"id"`
	Title          string `json:"title"`
	Detail         string `json:"detail"`
	ClientRelation int    `json:"clientrelation"`
}

type FinanceDetailFull struct {
	Id           int     `json:"id"`
	UserID       int     `json:"userid"`
	Cost         float64 `json:"cost"`
	Context      string  `json:"context"`
	HasBeenTaken bool    `json:"hasbeentaken"`
}

type FinanceDetailSend struct {
	UserID  int     `json:"userid"`
	Cost    float64 `json:"cost"`
	Context string  `json:"context"`
}

type FinanceDetailChangeHasBeenTaken struct {
	Ids []int `json:"ids"`
}

type ImageDataSent struct {
	Url         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ImageDataReceived struct {
	Id          int            `json:"id"`
	Datetime    time.Time      `json:"time"`
	Url         string         `json:"url"`
	Title       sql.NullString `json:"title"`
	Description sql.NullString `json:"description"`
}
