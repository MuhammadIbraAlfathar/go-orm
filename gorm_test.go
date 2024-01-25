package go_orm

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
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

func TestRawSql(t *testing.T) {
	var sample Sample
	err := db.Raw("select id, name from sample where id = ?", "1").Scan(&sample).Error
	assert.Nil(t, err)
	assert.Equal(t, "Jhon", sample.Name)

	var samples []Sample
	err = db.Raw("select id, name from sample").Scan(&samples).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(samples))
}

func TestSqlRows(t *testing.T) {
	var samples []Sample
	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	for rows.Next() {
		var id string
		var name string

		err := rows.Scan(&id, &name)
		assert.Nil(t, err)
		samples = append(samples, Sample{
			Id:   id,
			Name: name,
		})
	}

	assert.Equal(t, 4, len(samples))
}

func TestCreateUser(t *testing.T) {
	user := User{
		ID: "1",
		Name: Name{
			FirstName:  "Muhammad",
			MiddleName: "Ibra",
			LastName:   "Alfathar",
		},
		Password: "rahaisanegara",
	}

	//insert data
	response := db.Create(&user)
	assert.Nil(t, response.Error)
	assert.Equal(t, int64(1), response.RowsAffected)
}

func TestBatchInsert(t *testing.T) {
	var users []User

	for i := 2; i < 10; i++ {
		users = append(users, User{
			ID: strconv.Itoa(i),
			Name: Name{
				FirstName: "User " + strconv.Itoa(i),
			},
			Password: "rahaisa",
		})
	}

	db.Create(users)

}
