package learn_golang_gorm

import (
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