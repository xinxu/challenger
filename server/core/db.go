package core

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"time"
)

var _ = log.Printf

type MatchAnswerType int

const (
	MatchNotAnswer MatchAnswerType = 0
	MatchAnswering MatchAnswerType = 1
	MatchAnswered  MatchAnswerType = 2
)

type PlayerData struct {
	ID           uint      `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	MatchID      int       `json:"-"`
	Name         string    `json:"name"`
	Gold         int       `json:"gold"`
	LostGold     int       `json:"lostGold"`
	Energy       float64   `json:"energy"`
	Combo        int       `json:"combo"`
	Grade        string    `json:"grade"`
	Level        int       `json:"level"`
	LevelData    string    `json:"levelData"`
	HitCount     int       `json:"hitCount"`
	ControllerID string    `json:"cid"`
}

func (PlayerData) TableName() string {
	return "players"
}

type MatchData struct {
	ID           uint            `json:"id"`
	CreatedAt    time.Time       `json:"createdAt"`
	Mode         string          `json:"mode"`
	Elasped      float64         `json:"elasped"`
	Gold         int             `json:"gold"`
	Member       []PlayerData    `gorm:"ForeignKey:MatchID" json:"member"`
	RampageCount int             `json:"rampageCount"`
	AnswerType   MatchAnswerType `json:"answerType"`
	TeamID       string          `json:"teamID"`
}

func (MatchData) TableName() string {
	return "matches"
}

type DB struct {
	conn *gorm.DB
}

func NewDb() *DB {
	return &DB{}
}

func (db *DB) connect(path string) error {
	conn, err := gorm.Open("sqlite3", path)
	conn.AutoMigrate(&MatchData{}, &PlayerData{})
	if err != nil {
		return err
	}
	db.conn = conn
	return nil
}

func (db *DB) newMatch() *MatchData {
	var m = MatchData{}
	db.conn.Create(&m)
	return &m
}

func (db *DB) saveMatchData(m *MatchData) {
	db.conn.Save(m)
}

func (db *DB) getHistory(count int) []MatchData {
	var matches []MatchData
	db.conn.Limit(count).Preload("Member").Find(&matches)
	return matches
}

func (db *DB) startAnswer(mid int) *MatchData {
	var match MatchData
	db.conn.Where("id = ?", mid).First(&match)
	match.AnswerType = MatchAnswering
	db.conn.Save(&match)
	return &match
}
