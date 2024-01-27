package go_orm

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"strconv"
	"testing"
)

func OpenConnection() *gorm.DB {
	dialect := mysql.Open("root:@tcp(127.0.0.1:3306)/golang_orm?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dialect, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
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

func TestTransaction(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {

		err := tx.Create(&User{
			ID: "11",
			Name: Name{
				FirstName: "test 1",
			},
			Password: "test123",
		}).Error
		if err != nil {
			return err
		}

		err = tx.Create(&User{
			ID: "12",
			Name: Name{
				FirstName: "test 2",
			},
			Password: "test123",
		}).Error
		if err != nil {
			return err
		}

		err = tx.Create(&User{
			ID: "13",
			Name: Name{
				FirstName: "test 3",
			},
			Password: "test123",
		}).Error
		if err != nil {
			return err
		}

		return nil
	})

	assert.Nil(t, err)
}

func TestTransactionRollback(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {

		err := tx.Create(&User{
			ID: "14",
			Name: Name{
				FirstName: "test 4",
			},
			Password: "test123",
		}).Error
		if err != nil {
			return err
		}

		err = tx.Create(&User{
			ID: "12",
			Name: Name{
				FirstName: "test 2",
			},
			Password: "test123",
		}).Error
		if err != nil {
			return err
		}

		return nil
	})

	assert.NotNil(t, err)
}

func TestQuerySingleObject(t *testing.T) {
	//ambil data pertama
	users := User{}
	err := db.First(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, "1", users.ID)

	users = User{}
	err = db.Last(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, "9", users.ID)
}

func TestQuerySingleObjectInlineCondition(t *testing.T) {
	users := User{}
	err := db.Take(&users, "id = ?", "5").Error
	assert.Nil(t, err)
	assert.Equal(t, "5", users.ID)
}

func TestQueryAllObject(t *testing.T) {
	var users []User

	err := db.Find(&users, "id in ?", []string{"1", "2", "3", "4"}).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(users))
}

func TestQueryCondition(t *testing.T) {
	var users []User

	err := db.Where("first_name like ?", "%test%").
		Where("password = ?", "test123").Find(&users).Error

	assert.Nil(t, err)
	assert.Equal(t, 3, len(users))
	fmt.Println(users)
}

func TestOrCondition(t *testing.T) {
	var users []User

	err := db.Where("first_name like ?", "%User%").
		Or("password = ?", "rahasia").Find(&users).Error

	assert.Nil(t, err)

}

func TestNotCondition(t *testing.T) {
	var users []User

	err := db.Not("first_name like ?", "%User%").
		Where("password = ?", "rahasia").Find(&users).Error

	assert.Nil(t, err)
}

func TestSelectFields(t *testing.T) {
	var users []User
	err := db.Select("id", "first_name").Find(&users).Error
	assert.Nil(t, err)

	for _, user := range users {
		assert.NotNil(t, user.ID)
		assert.NotEqual(t, "", user.Name.FirstName)
	}
}

func TestStructCondition(t *testing.T) {
	var users []User

	userCondition := User{
		Name: Name{
			FirstName: "User 5",
			LastName:  "", //tidak bisa karena dianggap data default
		},
		Password: "rahaisa",
	}

	err := db.Where(userCondition).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

func TestMapCondition(t *testing.T) {
	var users []User

	userCondition := map[string]interface{}{
		"middle_name": "",
	}

	err := db.Where(userCondition).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 11, len(users))
}

func TestOrderLimitOffset(t *testing.T) {
	var users []User

	err := db.Order("id asc, first_name desc").Limit(5).Offset(5).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 5, len(users))
}

type UserResponse struct {
	ID        string
	FirstName string
	LastName  string
}

func TestQueryNonModel(t *testing.T) {
	var users []UserResponse

	err := db.Model(&User{}).Select("id", "first_name", "last_name").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 12, len(users))

	fmt.Println(users)
}

func TestUpdate(t *testing.T) {
	user := User{}

	err := db.Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)

	user.Name.FirstName = "Test"
	user.Name.MiddleName = "Update"
	user.Name.LastName = ""

	err = db.Save(&user).Error
	assert.Nil(t, err)
}

func TestUpdateSelectedColumn(t *testing.T) {
	err := db.Model(&User{}).Where("id = ?", "1").Updates(map[string]interface{}{
		"middle_name": "alfathar",
		"last_name":   "test",
	}).Error
	assert.Nil(t, err)

	err = db.Model(&User{}).Where("id = ?", "1").Update("password", "barudiubah").Error
	assert.Nil(t, err)

	err = db.Model(&User{}).Where("id = ?", "1").Updates(User{Name: Name{FirstName: "tsetttt"}}).Error
	assert.Nil(t, err)
}

func TestAutoIncrement(t *testing.T) {
	for i := 0; i < 10; i++ {
		userLog := UserLog{
			UserID: "000" + strconv.Itoa(i),
			Action: "test action",
		}
		result := db.Create(&userLog)
		assert.Nil(t, result.Error)
		assert.NotEqual(t, 0, userLog.ID)
		fmt.Println(userLog)
	}
}

func TestConflict(t *testing.T) {
	user := User{
		ID: "88",
		Name: Name{
			FirstName: "User 88 Updated",
		},
	}

	err := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&user).Error
	assert.Nil(t, err)
}

func TestDelete(t *testing.T) {
	var user []User
	err := db.Take(&user, "id = ?", "88").Error
	assert.Nil(t, err)

	err = db.Delete(&User{}, "id = ?", "88").Error
	assert.Nil(t, err)

	err = db.Where("id = ?", "10").Delete(&User{}).Error
	assert.Nil(t, err)
}

func TestSoftDelete(t *testing.T) {
	todo := Todo{
		UserId:      "3",
		Title:       "Todo 3",
		Description: "Description todo 3",
	}

	err := db.Create(&todo).Error
	assert.Nil(t, err)

	err = db.Where("id = ?", 3).Delete(&todo).Error
	assert.Nil(t, err)

}

func TestUnscoped(t *testing.T) {
	var todo Todo

	err := db.Unscoped().Take(&todo, "id = ?", 3).Error
	assert.Nil(t, err)
	fmt.Println(todo)

	err = db.Unscoped().Where("id = ?", 3).Delete(&todo).Error
	assert.Nil(t, err)

	var todos []Todo
	err = db.Unscoped().Find(&todos).Error
	assert.Nil(t, err)
	//assert.Equal(t, 1, len(todos))
}

func TestLock(t *testing.T) {
	var user User
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&user, "id = ?", "2").Error
		if err != nil {
			return err
		}

		user.Name.FirstName = "test locking"
		user.Name.MiddleName = "locking"
		user.Name.LastName = "last locking"

		err = tx.Save(&user).Error

		return err
	})
	assert.Nil(t, err)
}

func TestCreateWaller(t *testing.T) {
	wallet := Wallet{
		ID:      "2",
		UserId:  "2",
		Balance: 10000000,
	}

	err := db.Create(&wallet).Error
	assert.Nil(t, err)
}

func TestRetrieveRelation(t *testing.T) {
	var user User

	err := db.Model(&User{}).Preload("Wallet").Take(&user, "id = ?", "2").Error
	assert.Nil(t, err)
	fmt.Println(user)
}

func TestRetrieveRelationJoin(t *testing.T) {
	var user User

	err := db.Model(&User{}).Joins("Wallet").Take(&user, "users.id = ?", "29").Error
	assert.Nil(t, err)
	fmt.Println(user)
}

func TestUpsertOneToOne(t *testing.T) {
	user := User{
		ID: "29",
		Name: Name{
			FirstName: "User 29",
		},
		Password: "testing123",
		Wallet: Wallet{
			ID:      "29",
			UserId:  "29",
			Balance: 1000000,
		},
	}

	err := db.Create(&user).Error
	assert.Nil(t, err)
}

func TestUserAndAddress(t *testing.T) {
	user := User{
		ID:   "53",
		Name: Name{FirstName: "User 53"},
		Wallet: Wallet{
			ID:      "53",
			UserId:  "53",
			Balance: 2000000,
		},
		Address: []Address{
			{
				UserId:  "53",
				Address: "Jl. Nangka",
			},
		},
	}

	err := db.Create(&user).Error
	assert.Nil(t, err)
}

func TestPreloadJoin(t *testing.T) {
	var users User
	err := db.Model(&User{}).Preload("Address").Joins("Wallet").Find(&users).Error
	assert.Nil(t, err)
	fmt.Println(users)
}

func TestPreloadOneToMany(t *testing.T) {
	var users []User
	err := db.Model(&User{}).Preload("Address").Joins("Wallet").Take(&users, "users.id = ?", "53").Error
	assert.Nil(t, err)

	fmt.Println(users)
}

func TestBelongsTo(t *testing.T) {
	fmt.Println("Preload")
	var address []Address
	err := db.Model(&Address{}).Preload("User").Find(&address).Error
	assert.Nil(t, err)
	fmt.Println(address)

	address = []Address{}
	fmt.Println("Joins")
	err = db.Model(&Address{}).Joins("User").Find(&address).Error
	assert.Nil(t, err)
	fmt.Println(address)
}

func TestBelongsToWallet(t *testing.T) {
	fmt.Println("Preload")
	var wallet []Wallet
	err := db.Model(&Wallet{}).Preload("User").Find(&wallet).Error
	assert.Nil(t, err)
	fmt.Println(wallet)

	wallet = []Wallet{}
	fmt.Println("Joins")
	err = db.Model(&Wallet{}).Joins("User").Find(&wallet).Error
	assert.Nil(t, err)
	fmt.Println(wallet)
}

func TestCreateManyToMany(t *testing.T) {
	//product := Product{
	//	ID:    "P001",
	//	Name:  "Contoh Product",
	//	Price: 100000,
	//}

	//err := db.Create(&product).Error
	//assert.Nil(t, err)

	err := db.Table("user_like_product").Create(map[string]interface{}{
		"user_id":    "11",
		"product_id": "P001",
	}).Error

	assert.Nil(t, err)

	err = db.Table("user_like_product").Create(map[string]interface{}{
		"user_id":    "12",
		"product_id": "P001",
	}).Error

	assert.Nil(t, err)

}

func TestPreloadManyToMany(t *testing.T) {
	var product Product
	err := db.Preload("LikedByUsers").Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(product.LikedByUsers))
	fmt.Println(product.LikedByUsers)
}

func TestPreloadManyToManyProduct(t *testing.T) {
	var user User
	err := db.Preload("LikeProducts").Take(&user, "id = ?", "11").Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(user.LikeProducts))
	fmt.Println(user.LikeProducts)
}

func TestAssociationFind(t *testing.T) {
	var product Product
	err := db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	var user []User
	err = db.Model(&product).Where("first_name LIKE ?", "test%").Association("LikedByUsers").Find(&user)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(user))
	fmt.Println(user)

}

func TestAssociationAppend(t *testing.T) {
	var user User
	err := db.Take(&user, "id = ?", "29").Error
	assert.Nil(t, err)

	var product Product
	err = db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Append(&user)
	assert.Nil(t, err)

}

func TestAssociationReplace(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {

		var user User

		err := tx.Take(&user, "id = ?", "11").Error
		assert.Nil(t, err)

		wallet := Wallet{
			ID:      "01",
			UserId:  user.ID,
			Balance: 20000000,
		}

		err = tx.Model(&user).Association("Wallet").Replace(&wallet)
		assert.Nil(t, err)

		return err
	})

	assert.Nil(t, err)

}

func TestAssociationDelete(t *testing.T) {
	var user User
	err := db.Take(&user, "id = ?", "11").Error
	assert.Nil(t, err)

	var product Product
	err = db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Delete(&user)
	assert.Nil(t, err)

}

// menghapus semua user yang like product P001
func TestAssociationClear(t *testing.T) {

	var product Product
	err := db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Clear()
	assert.Nil(t, err)

}

func TestPreloadWithCondition(t *testing.T) {
	var user User
	err := db.Preload("Wallet", "balance > 1000").Take(&user, "id = ?", "20").Error
	assert.Nil(t, err)
	fmt.Println(user)
}

func TestPreloadWithNested(t *testing.T) {
	var wallets Wallet

	err := db.Preload("User.Address").Take(&wallets, "id = ?", "53").Error
	assert.Nil(t, err)

	fmt.Println(wallets)
	fmt.Println(wallets.User)
	fmt.Println(wallets.User.Address)
}

func TestPreloadAll(t *testing.T) {
	var user User
	err := db.Preload(clause.Associations).Take(&user, "id = ?", "50").Error
	assert.Nil(t, err)

	fmt.Println(user)
}

func TestJoinQuery(t *testing.T) {
	// inner join (jadi datanya harus ada di keduanya)
	var user []User
	err := db.Joins("join wallets on wallets.user_id = users.id").Find(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, 6, len(user))

	// left join
	user = []User{}
	err = db.Joins("Wallet").Find(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, 15, len(user))
}

func TestJoinQueryWithCondition(t *testing.T) {
	// inner join (jadi datanya harus ada di keduanya)
	var user []User
	err := db.Joins("join wallets on wallets.user_id = users.id AND wallets.balance > ?", 10000).Find(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, 6, len(user))
	fmt.Println(user)

	// left join
	user = []User{}
	err = db.Joins("Wallet").Where("Wallet.balance > ?", 10000).Find(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, 6, len(user))
	fmt.Println(user)
}

type AggregationResult struct {
	MaxBalance int64
	MinBalance int64
	SumBalance int64
}

func TestAggregation(t *testing.T) {
	var result AggregationResult
	err := db.Model(&Wallet{}).Select("max(balance) as max_balance", "min(balance) as min_balance").Take(&result).Error

	assert.Nil(t, err)
	assert.Equal(t, int64(20000000), result.MaxBalance)
	assert.Equal(t, int64(1000000), result.MinBalance)

	fmt.Println(result)
}

func TestAggregationGroupByHaving(t *testing.T) {
	var result []AggregationResult
	err := db.Model(&Wallet{}).Select("sum(balance) as sum_balance", "max(balance) as max_balance", "min(balance) as min_balance").
		Joins("User").Group("User.id").Having("sum(balance) > 1000000").
		Find(&result).Error

	assert.Nil(t, err)
	assert.Equal(t, 4, len(result))

	fmt.Println(result)
}

func TestContext(t *testing.T) {
	ctx := context.Background()

	var user []User
	err := db.WithContext(ctx).Find(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, 15, len(user))
}
