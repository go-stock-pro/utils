package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	env "github.com/joho/godotenv"
	"github.com/rohanraj7316/utils/postgres"
)

const (
	envFile = ".env.sample"
)

var loadEnv = env.Load

type classType int

const (
	one   classType = 1
	two   classType = 2
	three classType = 3
	four  classType = 4
)

type ParentsDetails struct {
}

type Student struct {
	ID             string         `gorm:"column:id"`
	Name           string         `gorm:"column:name"`
	Age            int            `gorm:"column:age"`
	Class          classType      `sql:"type:class" gorm:"column:class"`
	ParentsDetails ParentsDetails `gorm:"column:parentDetails"`
}

func setup(t *testing.T) (*postgres.Storage, error) {
	// load env config
	err := loadEnv(envFile)
	if err != nil {
		t.Error(err)
	}

	db, err := postgres.New()
	if err != nil {
		t.Fatalf("failed to create db connection: %v", err)
	}

	return db, nil
}

func TestMigration(t *testing.T) {
	db, err := setup(t)
	if err != nil {
		t.Error(err)
	}

	err = db.Migrate(context.Background(), Student{})
	if err != nil {
		t.Errorf("failed to run migration: %v", err)
	}
}

func TestCreateTransaction(t *testing.T) {
	query := "INSERT INTO Students (id, name, age, class) VALUES (?, ?, ?, ?)"

	u, _ := uuid.NewUUID()

	db, err := setup(t)
	if err != nil {
		t.Error(err)
	}

	st := Student{
		ID:   "10",
		Name: "a",
		Age:  27,
	}

	id := u.String()
	err = db.ScanWithCtx(context.Background(), &Student{}, query, id, st.Name, 12, one)
	if err != nil {
		t.Errorf("failed to insert data into db: %s", err)
	}
}

func TestSelect(t *testing.T) {
	query := "SELECT * FROM Students"

	db, err := setup(t)
	if err != nil {
		t.Error(err)
	}

	dst := &[]Student{}

	err = db.ScanWithCtx(context.Background(), dst, query)
	if err != nil {
		t.Errorf("failed to insert data into db: %s", err)
	}
}
