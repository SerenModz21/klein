package models

type Url struct {
	ID   interface{} `json:"id,omitempty" bson:"_id,omitempty"`
	Slug string      `json:"slug" bson:"slug"`
	Long string      `json:"long" bson:"long"`
}
