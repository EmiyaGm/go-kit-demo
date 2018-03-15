package alarmservice

import (
	"context"
	"errors"
	"fmt"
	"time"
	"gopkg.in/mgo.v2/bson"
	"log"
)

// Service describes a service that adds things together.
type Service interface {
	Create(ctx context.Context, ID string, FlowID uint32, Source string, Type string, Strategy string, Target string, SourceID string) (string, error)
	Add(ctx context.Context, ID string, FlowID uint32, Source string, Type string, Strategy string, Target string, SourceID string) (string, error)
	End(ctx context.Context, ID string, FlowID uint32, Source string, Type string, Strategy string, Target string, SourceID string) (string, error)
}

// New returns a basic Service with all of the expected middlewares wired in.
func New() Service {
	var svc Service
	{
		svc = NewBasicService()
	}
	return svc
}

var (
	// ErrTwoZeroes is an arbitrary business rule for the Add method.
	ErrTwoZeroes = errors.New("can't sum two zeroes")

	// ErrIntOverflow protects the Add method. We've decided that this error
	// indicates a misbehaving service and should count against e.g. circuit
	// breakers. So, we return it directly in endpoints, to illustrate the
	// difference. In a real service, this probably wouldn't be the case.
	ErrIntOverflow = errors.New("integer overflow")

	// ErrMaxSizeExceeded protects the Concat method.
	ErrMaxSizeExceeded = errors.New("result exceeds maximum size")
)

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService() Service {
	return basicService{}
}

type basicService struct{}


func (s basicService) Create(_ context.Context, ID string, FlowID uint32, Source string, Type string, Strategy string, Target string, SourceID string) (string, error){
	now := time.Now()
	data := bson.M{
		"_id":          ID,
		"source":       Source,
		//"location":     m.Location,
		//"trigger_time": m.Time,
		"target":       Target,
		//"status":       defaultStatus,
		"type":         Type,

		//"_p_vehicle_team":   m.Data["_p_vehicle_team"],
		//"speed":             m.Data["speed"],
		//"address":           m.Data["address"],
		//"_p_vehicle":        m.Data["_p_vehicle"],
		"_created_at":       now,
		"server_receive_at": now,
		//"_p_vehicle_models": m.Data["_p_vehicle_models"],
		"_updated_at": now,
		"ids":         []string{SourceID},
	}
	if err := dbc.Insert(data); err != nil {
		log.Println(err)
	}
	return "create alarm data",nil
}
func (s basicService) Add(_ context.Context, ID string, FlowID uint32, Source string, Type string, Strategy string, Target string, SourceID string) (string, error){
	now := time.Now()
	data := bson.M{
		//"end_time":    m.Time,
		"_updated_at": now,
	}

	if _, err := dbc.Upsert(bson.M{"_id": ID}, bson.M{
		"$set": data,
		"$push": bson.M{
			"ids": SourceID,
		},
	}); err != nil {
		log.Println(err)
	}
	fmt.Print("add alarm data")
	return "add alarm data",nil
}
func (s basicService) End(_ context.Context, ID string, FlowID uint32, Source string, Type string, Strategy string, Target string, SourceID string) (string,error){
	now := time.Now()
	data := bson.M{
		//"end_time":    m.Time,
		"_updated_at": now,
	}

	if _, err := dbc.Upsert(bson.M{"_id": ID}, bson.M{"$set": data}); err != nil {
		log.Println(err)
	}
	fmt.Print("end alarm data")
	return "end alarm data",nil
}
