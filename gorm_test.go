package learn_golang_gorm

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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