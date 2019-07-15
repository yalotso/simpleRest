package app

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// Transactional returns a handler that encloses the nested handlers with a DB transaction.
// If a nested handler returns an error or a panic happens, it will rollback the transaction.
// Otherwise it will commit the transaction after the nested handlers finish execution.
// By calling app.Context.SetRollback(true), you may also explicitly request to rollback the transaction.
func Transactional(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := db.Begin()
		err := tx.Error
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		rs := GetRequestScope(c)
		rs.SetTx(tx)

		c.Next()
		if len(c.Errors) > 0 || rs.Rollback() {
			tx.Rollback()
			return
		}
		tx.Commit()
		return
	}
}
