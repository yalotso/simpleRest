package app

import (
	"github.com/jinzhu/gorm"
	"time"
)

type RequestScope interface {
	Logger
	UserId() int
	SetUserId(id int)
	Tx() *gorm.DB
	SetTx(tx *gorm.DB)
	Rollback() bool
	SetRollback(bool)
	Now() time.Time
}

type requestScope struct {
	Logger
	now      time.Time
	userId   int
	rollback bool
	tx       *gorm.DB
}

func (rs *requestScope) UserId() int {
	return rs.userId
}

func (rs *requestScope) SetUserId(id int) {
	rs.userId = id
}

func (rs *requestScope) Tx() *gorm.DB {
	return rs.tx
}

func (rs *requestScope) SetTx(tx *gorm.DB) {
	rs.tx = tx
}

func (rs *requestScope) Rollback() bool {
	return rs.rollback
}

func (rs *requestScope) SetRollback(v bool) {
	rs.rollback = v
}

func (rs *requestScope) Now() time.Time {
	return rs.now
}

func newRequestScope(l Logger, now time.Time) RequestScope {
	return &requestScope{
		Logger: l,
		now:    now,
	}
}
