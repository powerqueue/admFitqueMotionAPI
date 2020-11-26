package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//MotionCollectionName - constant definition
const MotionCollectionName = "motion"

//MotionDefinition - struc definition
type MotionDefinition struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id" validate:"required"`
	LocationID    string             `bson:"LocationID" json:"LocationID" validate:"required"`
	SensorID      string             `bson:"SensorID" json:"SensorID" validate:"required"`
	MotionStartDt *time.Time         `bson:"MotionStartDt" json:"MotionStartDt" validate:"required"`
	MotionEndDt   *time.Time         `bson:"MotionEndDt" json:"MotionEndDt"`
	Measurements  []Measurement      `bson:"Measurements" json:"Measurements"`
}

//Measurement - struc definition
type Measurement struct {
	Report   string     `bson:"Report" json:"Report" validate:"required"`
	ReportDt *time.Time `bson:"ReportDt" json:"ReportDt" validate:"required"`
}
