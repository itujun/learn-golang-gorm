package learn_golang_gorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenConnection() *gorm.DB {
	dialect := mysql.Open("root:@tcp(localhost:3306)/learn_golang_gorm?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dialect)
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