package go_orm

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func OpenConnection() *gorm.DB {
	dialect := mysql.Open("root:@tcp(127.0.0.1:3306)/golang_orm?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

var db = OpenConnection()

func TestConnectionDatabase(t *testing.T) {
	assert.NotNil(t, db)
}

func TestExecuteSQL(t *testing.T) {
	err := db.Exec("insert into sample (id, name) values (?,?)", "1", "Jhon").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample (id, name) values (?,?)", "2", "Doe").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample (id, name) values (?,?)", "3", "Shella").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample (id, name) values (?,?)", "4", "Sarah").Error
	assert.Nil(t, err)
}

type Sample struct {
	Id   string
	Name string
}
