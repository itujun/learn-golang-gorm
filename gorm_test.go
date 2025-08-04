package learn_golang_gorm

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

func OpenConnection() *gorm.DB {
	dialect := mysql.Open("root:@tcp(localhost:3306)/learn_golang_gorm?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dialect, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Set log level to Info
	})
	if err != nil {
		panic(err)		
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db
}

var db = OpenConnection()

func TestOpenConnection(t *testing.T) {
	assert.NotNil(t, db)
}

func TestExecuteSQL(t *testing.T){
	err := db.Exec("insert into sample(id,name) values(?, ?)", "1", "Lev").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample(id,name) values(?, ?)", "2", "Tempest").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample(id,name) values(?, ?)", "3", "Vex").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample(id,name) values(?, ?)", "4", "Cait").Error
	assert.Nil(t, err)
}

type Sample struct {
	Id   string
	Name string
}

func TestRawSQL(t *testing.T) {
	var sample Sample
	err := db.Raw("select id, name from sample where id = ?", "1").Scan(&sample).Error
	assert.Nil(t, err)
	assert.Equal(t, "Lev", sample.Name)

	var samples []Sample
	err = db.Raw("select id, name from sample").Scan(&samples).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(samples))
}

// LEBIH PANJANG JIKA MENGGUNAKAN Rows
// Menggunakan Rows untuk mendapatkan hasil query
// Rows adalah metode yang digunakan untuk mendapatkan hasil query dalam bentuk baris.
func TestSqlRow(t *testing.T){
	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	var samples []Sample
	for rows.Next() {
		var sample Sample
		err := rows.Scan(&sample.Id, &sample.Name)
		assert.Nil(t, err)
		samples = append(samples, sample)
	}
	assert.Equal(t, 4, len(samples))
}

// LEBIH RINGKAS JIKA MENGGUNAKAN ScanRows
// ScanRows adalah metode yang disediakan oleh GORM untuk memindai hasil query ke dalam slice atau struct.
// Ini mengurangi boilerplate code yang diperlukan untuk memindai setiap baris secara manual.
func TestScanRow(t *testing.T){
	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	var samples []Sample
	for rows.Next() {
		err := db.ScanRows(rows, &samples)
		assert.Nil(t, err)
		// Tidak perlu lagi melakukan append, karena ScanRows sudah mengisi slice samples
	}
	assert.Equal(t, 4, len(samples))
}

func TestCreateUser(t *testing.T) {
	user := User{
		ID: "1",
		Name: Name{
			FirstName:  "Lev",
			MiddleName: "Tempest",
			LastName:   "Vex",
		},
		Password: "secret",
	}

	response := db.Create(&user)
	assert.Nil(t, response.Error)
	assert.Equal(t, int64(1), response.RowsAffected)
}

func TestBatchInsertUsers(t *testing.T) {
	var users []User
	for i := 2; i < 10; i++ {
		users = append(users, User{
			ID:       strconv.Itoa(i),
			Name:     Name{
				FirstName: "User" + strconv.Itoa(i),
			},
			Password: "secret",
		})
	}

	response := db.Create(&users)
	assert.Nil(t, response.Error)
	assert.Equal(t, 8, int(response.RowsAffected))
}

func TestManualTransactionSuccess(t *testing.T){
	tx := db.Begin()
	defer tx.Rollback() // Pastikan rollback jika tidak commit

	err := tx.Create(&User{ID: "10", Name: Name{FirstName: "User10"}, Password: "secret"}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{ID: "11", Name: Name{FirstName: "User11"}, Password: "secret"}).Error
	assert.Nil(t, err)

	if err == nil {
		tx.Commit() // Commit jika tidak ada error
	}
}

func TestManualTransactionFailure(t *testing.T){
	tx := db.Begin()
	defer tx.Rollback() // Pastikan rollback jika tidak commit

	err := tx.Create(&User{ID: "12", Name: Name{FirstName: "User12"}, Password: "secret"}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{ID: "11", Name: Name{FirstName: "User11"}, Password: "secret"}).Error
	assert.Nil(t, err)

	if err == nil {
		tx.Commit() // Commit jika tidak ada error
	}
}

func TestTransactionSuccessWithGorm(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&User{ID: "12", Name: Name{FirstName: "User12"}, Password: "secret"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&User{ID: "13", Name: Name{FirstName: "User13"}, Password: "secret"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&User{ID: "14", Name: Name{FirstName: "User14"}, Password: "secret"}).Error; err != nil {
			return err
		}
		return nil
	})
	assert.Nil(t, err)
}

func TestTransactionFailureWithGorm(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&User{ID: "15", Name: Name{FirstName: "User15"}, Password: "secret"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&User{ID: "13", Name: Name{FirstName: "User13"}, Password: "secret"}).Error; err != nil {
			return err
		}
		return nil
	})
	assert.NotNil(t, err)
}

func TestSingleQuery(t *testing.T){
	user := User{}
	err := db.First(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, "1", user.ID)

	user = User{}
	err = db.Last(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, "9", user.ID)
}

func TestSingleSingleObjectInLineCondition(t *testing.T){
	user := User{}
	err := db.Take(&user, "id = ?" , "5").Error
	assert.Nil(t, err)
	assert.Equal(t, "5", user.ID)
	assert.Equal(t, "User5", user.Name.FirstName)
}

func TestQueryAllObjects(t *testing.T){
	var users []User
	err := db.Find(&users, "id in ?" , []string{"1", "2", "3", "4"}).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(users))
}

func TestQueryCondition(t *testing.T){
	var users []User
	err := db.Where("first_name like ?" , "%User%").Where("password = ?", "secret").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 13, len(users))
}

func TestOrOperator(t *testing.T){
	var users []User
	err := db.Where("first_name like ?" , "%User%").Or("password = ?", "secret").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 14, len(users))
}

func TestNotOperator(t *testing.T){
	var users []User
	err := db.Not("first_name like ?" , "%User%").Where("password = ?", "secret").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

func TestSelectFields(t *testing.T){
	var users []User
	err := db.Select("id, first_name").Find(&users).Error
	assert.Nil(t, err)

	for _, user := range users {
		assert.NotEmpty(t, user.ID)
		assert.NotEqual(t, "", user.Name.FirstName)
	}

	assert.Equal(t, 14, len(users))
}

func TestStructCondition(t *testing.T){
	userCondition := User{
		Name: Name{
			FirstName: "User5",
			LastName: "", // LastName akan diabaikan karena dianggap default value
		},
		Password: "secret",
	}
	var users []User
	err := db.Where(userCondition).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

func TestMapCondition(t *testing.T){
	mapCondition := map[string]interface{}{
		"middle_name": "",
	}
	var users []User
	err := db.Where(mapCondition).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 13, len(users))
}

func TestOrderLimitOffset(t *testing.T){
	var users []User
	err := db.Order("id asc, first_name desc").Limit(5).Offset(5).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 5, len(users))
}

type UserResponse struct {
	ID 			string
	FirstName 	string
	LastName	string
}

func TestQueryNonModel(t *testing.T){
	var users []UserResponse
	// dari model users, pilih id, first_name, last_name, lalu simpan ke variable users
	err := db.Model(&User{}).Select("id", "first_name", "last_name").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 14, len(users))
	fmt.Println(users);
}

func TestUpdate(t *testing.T){
	user := User{}
	err := db.Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)

	user.Name.FirstName = "Levi"
	user.Name.MiddleName = ""
	user.Name.LastName = "kun"
	user.Password = "secret123"

	err = db.Save(&user).Error
	assert.Nil(t, err)
}

func TestUpdateSelectedColumns(t *testing.T){
	// updates versi 1
	// dari model (tabel) users, yang id = 1, lakukan update...
	err := db.Model(&User{}).Where("id = ?", "1").Updates(map[string]interface{}{
		"middle_name" : "",
		"last_name" : "San",
	}).Error
	assert.Nil(t, err)

	// update versi 2
	err = db.Model(&User{}).Where("id = ?", "1").Update("password", "secretberubah").Error
	assert.Nil(t, err)

	// update versi 3
	// karena diatas sudah menentukan model yg akan digunakan, maka menambahakan kode db.Model() bisa optional
	err = db.Where("id = ?", "1").Updates(User{
		Name: Name{
			FirstName: "Lev",
			LastName: "Tempest",
		},
	}).Error
	assert.Nil(t, err)
}

func TestAutoIncrement(t *testing.T){
	for i :=  0; i < 10; i++ {
		userLog := UserLog{
			UserId: "1",
			Action: "Test Action",
		}

		err := db.Create(&userLog).Error
		assert.Nil(t, err) 					// Pastikan tidak ada error saat menyimpan UserLog
		assert.NotEqual(t, 0, userLog.ID) 	// Pastikan ID tidak nol
		fmt.Println(userLog.ID) 			// Cetak ID untuk verifikasi;
	}
}

func TestSaveOrUpdate(t *testing.T){
	userLog := UserLog{
		UserId: "1",
		Action: "Test Action",
	}

	err := db.Save(&userLog).Error	// insert
	assert.Nil(t, err)

	userLog.UserId = "2"
	err = db.Save(&userLog).Error	// update
	assert.Nil(t, err)
}

func TestSaveOrUpdateNonAutoIncrement(t *testing.T){
	user := User{
		ID	: "99",
		Name: Name{
			FirstName: "User99",
		},
	}

	err := db.Save(&user).Error	// insert
	assert.Nil(t, err)

	user.Name.FirstName = "User 99 Updated"
	err = db.Save(&user).Error	// update
	assert.Nil(t, err)
}

func TestOnConflict(t *testing.T){
	user := User{
		ID	: "88",
		Name: Name{
			FirstName: "User88",
		},
	}

	err := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&user).Error // insert
	assert.Nil(t, err)
}

func TestDeke(t *testing.T){
	var user User
	// cara 1 : select dulu baru update
	err := db.Take(&user, "id = ?", "88").Error
	assert.Nil(t, err)
	err = db.Delete(&user).Error
	assert.Nil(t, err)

	// cara 2 : langsung delete
	err = db.Delete(&User{}, "id = ?", "99").Error
	assert.Nil(t, err)

	// cara 3 : langsung delete
	err = db.Where("id = ?", "14").Delete(&User{}).Error
	assert.Nil(t, err)
}

func TestSoftDelete(t *testing.T){
	todo := Todo{
		UserId:  		"1",
		Title:   		"Title 1",
		Description:   	"Description 1",
	}
	err := db.Create(&todo).Error	// insert
	assert.Nil(t, err)

	err = db.Delete(&todo).Error	// delete (update deleted_at)
	assert.Nil(t, err)
	// assert.NotNil(t, todo.Deleted_at) // 

	var todos []Todo
	err = db.Find(&todos).Error		// select (where deleted_at is null)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(todos))
}

func TestUnscoped(t *testing.T){
	var todo Todo
	err := db.Unscoped().First(&todo, "id = ?", 1).Error
	assert.Nil(t, err)

	err = db.Unscoped().Delete(&todo).Error	// delete permanen
	assert.Nil(t, err)

	var todos []Todo
	err = db.Unscoped().Find(&todos).Error
	assert.Nil(t, err)
}

func TestLock(t *testing.T){
	err := db.Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.Clauses(clause.Locking{
			Strength: "UPDATE", 			// Lock for update
		}).First(&user, "id = ?", "1").Error
		if err != nil {
			return err
		}

		user.Name.FirstName = "Levi"
		user.Name.LastName = "Tempesto"
		err = tx.Save(&user).Error
		return err
		})
	assert.Nil(t, err)
}

func TestCreateWallet(t *testing.T){
	wallet := Wallet{
		ID: "1",
		UserId: "1",
		Balance: 1000000,
	}

	err := db.Create(&wallet).Error
	assert.Nil(t, err)
}

func TestRetrieveRelation(t *testing.T){
	var user User
	err := db.Model(&User{}).Preload("Wallet").Take(&user, "id = ?", "1").Error // preload melakukan 2x query
	assert.Nil(t, err)

	assert.Equal(t, 1000000, user.Wallet.Balance)
}

func TestRetrieveRelationJoin(t *testing.T){
	var user User
	err := db.Model(&User{}).Joins("Wallet").Take(&user, "users.id = ?", "1").Error // joins melakukan 1x query (cocok untuk relasi One to One)
	assert.Nil(t, err)

	assert.Equal(t, 1000000, user.Wallet.Balance)
}

func TestAutoCreateUpdate(t *testing.T){
	user := User{
		ID: "20",
		Password: "secret",
		Name: Name{
			FirstName: "User 20",
		},
		Wallet: Wallet{
			ID: "20",
			UserId: "20",
			Balance: 1000000,
		},
	}

	err := db.Create(&user).Error
	assert.Nil(t, err)
}

func TestSkipAutoCreateUpdate(t *testing.T){
	user := User{
		ID: "21",
		Password: "secret",
		Name: Name{
			FirstName: "User 21",
		},
		Wallet: Wallet{
			ID: "21",
			UserId: "21",
			Balance: 1000000,
		},
	}

	err := db.Omit(clause.Associations).Create(&user).Error
	assert.Nil(t, err)
}

func TestUserAndAddresses(t *testing.T){
	user := User{
		ID: "2",
		Password: "secret",
		Name: Name{
			FirstName: "User 2",
		},
		Wallet: Wallet{
			ID: "2",
			UserId: "2",
			Balance: 1000000,
		},
		Addresses: []Address{
			{
				Address: "Jalan 1",
			},
			{
				Address: "Jalan 2",
			},
		},
	}

	// err := db.Create(&user).Error
	err := db.Save(&user).Error
	assert.Nil(t, err)
}

func TestPreloadJoinOneToMany(t *testing.T){
	var users []User
	err := db.Model(&User{}).Preload("Addresses").Joins("Wallet").Find(&users).Error // preload untuk relasi One to Many, join untuk relasi One to One
	assert.Nil(t, err)
}

func TestTakePreloadJoinOneToMany(t *testing.T){
	var user User
	err := db.Model(&User{}).Preload("Addresses").Joins("Wallet").Take(&user, "users.id = ?", "50").Error // preload untuk relasi One to Many, join untuk relasi One to One
	assert.Nil(t, err)
}

func TestBelongsTo(t *testing.T){
	fmt.Println("Preload");
	var addresses []Address
	err := db.Model(&Address{}).Preload("User").Find(&addresses).Error // preload melakukan 2x query
	assert.Nil(t, err)
	assert.Equal(t, 4, len(addresses))

	fmt.Println("Joins");
	addresses = []Address{}
	err = db.Model(&Address{}).Joins("User").Find(&addresses).Error	// joins melakukan 1x query
	assert.Nil(t, err)
	assert.Equal(t, 4, len(addresses))
}

func TestBelongsToWallet(t *testing.T){ 	// BelongsToOneToOne
	fmt.Println("Preload");
	var wallets []Wallet
	err := db.Model(&Wallet{}).Preload("User").Find(&wallets).Error // preload melakukan 2x query
	assert.Nil(t, err)

	fmt.Println("Joins");
	wallets = []Wallet{}
	err = db.Model(&Wallet{}).Joins("User").Find(&wallets).Error	// joins melakukan 1x query
	assert.Nil(t, err)
}

func TestCreateManyToMany(t *testing.T){	// cara manual (lebih baik menggunakan associations /seperti baris kode 635)
	product := Product{
		ID: "P001",
		Name: "Product 1",
		Price: 1000000,
	}
	err := db.Create(&product).Error
	assert.Nil(t, err)

	err = db.Table("user_like_product").Create(map[string]interface{}{
		"user_id" : "1",
		"product_id" : "P001",
	}).Error
	assert.Nil(t, err)

	err = db.Table("user_like_product").Create(map[string]interface{}{
		"user_id" : "2",
		"product_id" : "P001",
	}).Error
	assert.Nil(t, err)
}

func TestPreloadManyToMany(t *testing.T){
	var product Product
	err := db.Preload("LikedByUsers").Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(product.LikedByUsers))
}

func TestPreloadManyToManyUser(t *testing.T){
	var user User
	err := db.Preload("LikeProducts").Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(user.LikeProducts))
}

func TestAssociationFind(t *testing.T){
	var product Product
	err := db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	var users []User
	err = db.Model(&product).Where("users.first_name LIKE ?", "User%").Association("LikedByUsers").Find(&users)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

func TestAssociationAdd(t *testing.T){
	var user User
	err := db.Take(&user, "id = ?", "3").Error
	assert.Nil(t, err)

	var product Product
	err = db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Append(&user)
	assert.Nil(t, err)
}

func TestAssociationReplace(t *testing.T){	// cocok untuk relasi Belongs To (One to One)
	err := db.Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.Take(&user, "id = ?", "1").Error
		assert.Nil(t, err)

		wallet := Wallet{
			ID: "01",
			UserId: user.ID,
			Balance: 1000000,
		}

		err = tx.Model(&user).Association("Wallet").Replace(&wallet)
		return err
	})
	assert.Nil(t, err)
}

func TestAssociationDelete(t *testing.T){
	var user User
	err := db.Take(&user, "id = ?", "3").Error
	assert.Nil(t, err)

	var product Product
	err = db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Delete(&user)
	assert.Nil(t, err)
}

func TestAssociationClear(t *testing.T){
	var product Product
	err := db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Clear()
	assert.Nil(t, err)
}

func TestPreloadingWithCondition(t *testing.T){
	var user User
	err := db.Preload("Wallet", "balance > ?", 1000000).Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)

	fmt.Println(user);
}

func TestPreloadingNested(t *testing.T){
	var wallet Wallet
	err := db.Preload("User.Addresses").Take(&wallet, "id = ?", "1").Error
	assert.Nil(t, err)

	fmt.Println(wallet);
	fmt.Println(wallet.User);
	fmt.Println(wallet.User.Addresses);
}

func TestPreloadingAll(t *testing.T){
	var user User
	err := db.Preload(clause.Associations).Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)
}

func TestJoinQuery(t *testing.T){
	var users []User
	err := db.Joins("join wallets on wallets.user_id = users.id").Find(&users).Error	// inner join
	assert.Nil(t, err)
	assert.Equal(t, 4, len(users))

	users = []User{}
	err = db.Joins("Wallet").Find(&users).Error	// left join
	assert.Nil(t, err)
	assert.Equal(t, 16, len(users))
}

func TestJoinWithCondition(t *testing.T){
	var users []User
	err := db.Joins("join wallets on wallets.user_id = users.id AND wallets.balance > ?", 500000).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(users))

	users = []User{}
	err = db.Joins("Wallet").Where("Wallet.balance > ?", 500000).Find(&users).Error	// alias menggunakan nama field	
	assert.Nil(t, err)
	assert.Equal(t, 4, len(users))
}

func TestCount(t *testing.T){
	var count int64
	err := db.Model(&User{}).Joins("Wallet").Where("Wallet.balance > ?", 500000).Count(&count).Error
	assert.Nil(t, err)
	assert.Equal(t, int64(4), count)
}

type AggregationResult struct {
	TotalBalance int64
	MinBalance int64
	MaxBalance int64
	AvgBalance float64
}

func TestAggregation(t *testing.T){
	var result AggregationResult
	err := db.Model(&Wallet{}).Select("sum(balance) as total_balance", "max(balance) as max_balance", "min(balance) as min_balance", "avg(balance) as avg_balance").Take(&result).Error
	assert.Nil(t, err)

	assert.Equal(t, int64(4000000), result.TotalBalance)
	assert.Equal(t, int64(1000000), result.MaxBalance)
	assert.Equal(t, int64(1000000), result.MinBalance)
	assert.Equal(t, float64(1000000), result.AvgBalance)
}

func TestAggregationGroupByAndHaving(t *testing.T){
	var results []AggregationResult
	err := db.Model(&Wallet{}).Select("sum(balance) as total_balance", "max(balance) as max_balance", "min(balance) as min_balance", "avg(balance) as avg_balance").Joins("User").Group("User.id").Having("sum(balance) > ?", 500000).Find(&results).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(results))
}

func TestWithContext(t *testing.T){
	ctx := context.Background()
	
	var users []User
	err := db.WithContext(ctx).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 16, len(users))
}

func BrokeWalletBalance(db *gorm.DB) *gorm.DB {
	return db.Where("balance = ?", 0)
}

func SultanWalletBalance(db *gorm.DB) *gorm.DB {
	return db.Where("balance = ?", 1000000)
}

func TestScopes(t *testing.T){
	var wallets []Wallet
	err := db.Scopes(BrokeWalletBalance).Find(&wallets).Error
	assert.Nil(t, err)
	fmt.Println("BrokeWalletBalance", len(wallets));

	err = db.Scopes(SultanWalletBalance).Find(&wallets).Error
	assert.Nil(t, err)
	fmt.Println("SultanWalletBalance", len(wallets));
}

func TestMigrator(t *testing.T){	// lebih disarankan menggunakan migrator manual
	err := db.Migrator().AutoMigrate(&GuestBook{})
	assert.Nil(t, err)
}

// HOOK untuk Create
// 1. // begin transaction
// 2. BeforeSave()
// 3. BeforeCreate()
// 4. // save before assosations
// 5. // insert into database
// 6. // save after assosations
// 7. AfterCreate()
// 8. AfterSave()
// 9. // commit or rollback transaction

// HOOK untuk Update
// 1. // begin transaction
// 2. BeforeSave()
// 3. BeforeUpdate()
// 4. // save before assosations
// 5. // update database
// 6. // save after assosations
// 7. AfterUpdate()
// 8. AfterSave()
// 9. // commit or rollback transaction

// HOOK untuk Delete
// 1. // begin transaction
// 2. BeforeDelete()
// 3. // delete from database
// 4. AfterDelete()
// 5. // commit or rollback transaction

// HOOK untuk Find
// 1. // load data from database
// 2. // preloading (eager loading)
// 3. AfterFind()

func TestHook(t *testing.T){
	user := User{
		Password: "secret",
		Name: Name{
			FirstName: "User 100",
		},
	}
	err := db.Create(&user).Error
	assert.Nil(t, err)
	assert.NotEqual(t, "", user.ID)
	fmt.Println("user.ID: ", user.ID);
}