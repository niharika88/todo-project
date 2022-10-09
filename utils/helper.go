package utils

import "go.mongodb.org/mongo-driver/bson"

func ToMongoBson(v interface{}) (bsonD *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &bsonD)
	return
}

func Pointer(str string) *string {
	return &str
}
