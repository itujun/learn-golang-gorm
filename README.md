# BELAJAR GOLANG GORM

## DEPENDENCY TAMBAHAN

- (Testify)[github.com/stretchr/testify]

Installation

```bash
go get github.com/stretchr/testify
```

- (Gorm)[https://gorm.io/docs/]

Installation

```bash
go get -u gorm.io/gorm

# lalu pilih driver db yg akan digunakan
# go get -u gorm.io/driver/sqlite
# go get -u gorm.io/driver/mysql
# go get -u gorm.io/driver/postgres
# go get -u gorm.io/driver/sqlserver
# go get -u gorm.io/driver/clickhouse

# pada projek ini, saya menggunakan mysql
go get -u gorm.io/driver/mysql
```

### SETUP PROJECT

1. Buat Database: learn_golang_gorm
2. Hubungkan project ke Database
   (cek dokumentasinya disini)[https://gorm.io/docs/connecting_to_the_database.html]

```bash
import (
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

func main() {
  // refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
  dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
```

3. Buat table sample pada Database

```bash
CREATE TABLE IF NOT EXISTS sample(
    id VARCHAR(100) NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
)engine = INNODB
```
