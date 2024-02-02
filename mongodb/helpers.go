package mongodb

import (
	"reflect"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Storage) getCollection(st any) (*mongo.Collection, error) {
	if st != nil {
		t := reflect.TypeOf(st)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		// todo: need to make word plural too.
		collName := strcase.ToSnake(t.Name())

		return s.DB.Collection(collName, s.collOptions), nil
	} else {
		//  todo: need to improve this error msg. also, need to add handling for struct validation.
		return nil, errors.Errorf("empty collection")
	}
}

// TODO: need to add struct validation
func (s *Storage) IsValidStruct(st any) error {
	return nil
}
