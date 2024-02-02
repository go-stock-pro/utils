package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	Db    *gorm.DB
	SqlDb *sql.DB
}

// NewDBStringFromConfig build database connection string from config file.
func NewDBStringFromConfig(config Config) (string, error) {
	var dbParams []string
	dbParams = append(dbParams, fmt.Sprintf("user=%s", config.User))
	dbParams = append(dbParams, fmt.Sprintf("host=%s", config.Host))
	dbParams = append(dbParams, fmt.Sprintf("port=%s", config.Port))
	dbParams = append(dbParams, fmt.Sprintf("dbname=%s", config.DbName))
	if password := config.Password; password != "" {
		dbParams = append(dbParams, fmt.Sprintf("password=%s", password))
	}
	dbParams = append(dbParams, fmt.Sprintf("sslmode=%s",
		config.SslMode))
	dbParams = append(dbParams, fmt.Sprintf("connect_timeout=%d",
		config.ConnectionTimeout))
	dbParams = append(dbParams, fmt.Sprintf("statement_timeout=%d",
		config.StatementTimeout))
	dbParams = append(dbParams, fmt.Sprintf("idle_in_transaction_session_timeout=%d",
		config.IdleInTransactionSessionTimeout))
	dbParams = append(dbParams, "TimeZone=Asia/Kolkata")

	return strings.Join(dbParams, " "), nil
}

// New used to create new connection.
func New(config ...Config) (*Storage, error) {
	cfg, err := configDefault(config...)
	if err != nil {
		return nil, err
	}

	dbString, err := NewDBStringFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN: dbString,
			},
		), cfg.gormCfg)
	if err != nil {
		return nil, errors.Errorf("failed to initialize database: %v", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, errors.Errorf("failed to connect to database: %v", err)
	}

	sqlDb.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDb.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDb.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	sqlDb.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return &Storage{Db: db, SqlDb: sqlDb}, nil
}

// ExecWithCtx use it when you don't care abt the output. it doesn't returns error.
func (s *Storage) ExecWithCtx(ctx context.Context, query string, values ...interface{}) {
	s.Db.WithContext(ctx).Exec(query, values...)
}

// ScanWithCtx use it when you care abt the output. you gonna get scanned output in dst.
func (s *Storage) ScanWithCtx(ctx context.Context, dst interface{},
	query string, values ...interface{}) error {
	if err := s.Db.WithContext(ctx).Raw(query, values...).Scan(dst).Error; err != nil {
		return err
	}

	return nil
}

func (s *Storage) Migrate(ctx context.Context, dst ...interface{}) error {
	return s.Db.WithContext(ctx).AutoMigrate(dst...)
}

func (s *Storage) Create(ctx context.Context, dst interface{}) error {
	return s.Db.WithContext(ctx).Create(dst).Error
}

func (s *Storage) Read(ctx context.Context, projection []string, query, dst interface{}, limit, offset int) error {
	return s.Db.WithContext(ctx).Select(projection).Where(query).Find(dst).Limit(limit).Offset(offset).Error
}

func (s *Storage) Updates(ctx context.Context, filter, update interface{}) error {
	return s.Db.WithContext(ctx).Model(filter).Updates(update).Error
}

func (s *Storage) Delete(ctx context.Context, filter interface{}, update interface{}) error {
	return s.Db.WithContext(ctx).Model(filter).Updates(update).Error
}
