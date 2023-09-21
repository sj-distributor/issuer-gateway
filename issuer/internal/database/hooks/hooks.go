package hooks

import (
	"context"
	"github.com/pygzfei/issuer-gateway/utils"
	"gorm.io/gorm"
	"reflect"
)

type GenerateSnowflakeId struct {
	Name string
}

func (g GenerateSnowflakeId) Initialize(db *gorm.DB) {
	if db.Statement.Schema != nil {
		for _, field := range db.Statement.Schema.PrimaryFields {
			if field.Name == "Id" {
				if reflect.TypeOf(field).Kind() == reflect.Ptr {
					typeOf := reflect.TypeOf(db.Statement.Dest)
					if typeOf.Elem().Kind() == reflect.Slice {
						valueOf := reflect.ValueOf(db.Statement.Dest)
						for i := 0; i < valueOf.Elem().Len(); i++ {
							valueOf.Elem().Index(i).FieldByName("Id").Set(reflect.ValueOf(utils.Id()))
						}
					} else {
						dbErr := field.Set(context.Background(), db.Statement.ReflectValue, utils.Id())
						if dbErr != nil {
							panic(dbErr)
						}
					}
				}
			}
		}
	}
}

func GenerateId() *GenerateSnowflakeId {
	return &GenerateSnowflakeId{
		Name: "GenerateSnowflakeId",
	}
}
