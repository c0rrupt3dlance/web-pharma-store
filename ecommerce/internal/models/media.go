package models

import (
	"time"
)

type FileDataType struct {
	FileName string
	Data     []byte
	DataType string
	Position int
}

type MediaUrl struct {
	ObjectId string `json:"object_id,omitempty"`
	Url      string `json:"url"`
	Position int    `json:"position"`
}

type ProductMedia struct {
	Id        int
	ProductId int
	MediaId   string
	CreatedAt time.Time
}
