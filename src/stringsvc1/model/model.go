package model

import (
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	collection = "vehicle_warning"

	defaultStatus = "未处理"

	// StrategyEnd 结束
	StrategyEnd = "end"
	// StrategyAdd 增加点
	StrategyAdd = "add"
	// StrategyCreate 创建
	StrategyCreate = "create"
)

// Event 报警事件
type Event struct {
	start time.Time
	end   time.Time

	addons map[string]interface{}
}

// Message 报警消息
type Message struct {
	ID     string
	FlowID uint32
	Source string
	Type   string
	Time   time.Time
	// new add end
	Strategy string

	Target   string
	SourceID string
	Location []float64
	Data     map[string]interface{}
}

// Config 配置
type Config struct {
	Host string
	Name string
}

// new 更新报警记录
func create(m *Message) {
	now := time.Now()
	data := bson.M{
		"_id":          m.ID,
		"source":       m.Source,
		"location":     m.Location,
		"trigger_time": m.Time,
		"target":       m.Target,
		"status":       defaultStatus,
		"type":         m.Type,

		"_p_vehicle_team":   m.Data["_p_vehicle_team"],
		"speed":             m.Data["speed"],
		"address":           m.Data["address"],
		"_p_vehicle":        m.Data["_p_vehicle"],
		"_created_at":       now,
		"server_receive_at": now,
		"_p_vehicle_models": m.Data["_p_vehicle_models"],

		"_updated_at": now,
		"ids":         []string{m.SourceID},
	}

	if err := c.Insert(data); err != nil {
		log.Println(err)
	}
}

// add 增加一个报警点
func add(m *Message) {
	now := time.Now()
	data := bson.M{
		"end_time":    m.Time,
		"_updated_at": now,
	}

	if _, err := c.Upsert(bson.M{"_id": m.ID}, bson.M{
		"$set": data,
		"$push": bson.M{
			"ids": m.SourceID,
		},
	}); err != nil {
		log.Println(err)
	}
}

// end 结束一次报警
func end(m *Message) {
	now := time.Now()
	data := bson.M{
		"end_time":    m.Time,
		"_updated_at": now,
	}

	if _, err := c.Upsert(bson.M{"_id": m.ID}, bson.M{"$set": data}); err != nil {
		log.Println(err)
	}
}
