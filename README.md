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

4. Buat table users pada database

```bash
CREATE TABLE IF NOT EXISTS users(
    id VARCHAR(100) NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)engine = InnoDB;
```

5. Edit table users pada database

```bash
-- Ubah nama kolom 'name' menjadi 'first_name'
ALTER TABLE users
  CHANGE name first_name VARCHAR(255);

-- Tambahkan kolom baru
ALTER TABLE users
  ADD COLUMN middle_name VARCHAR(100) NULL AFTER first_name,
  ADD COLUMN last_name VARCHAR(100) NULL AFTER middle_name;

```

6. Buat table user_logs pada database

```bash
CREATE TABLE IF NOT EXISTS user_logs(
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)engine = InnoDB;
```

7. Ubah tipe kolom menjadi bigint pada created_at dan updated_at pada table user_logs

```bash
alter table user_logs
	MODIFY created_at BIGINT not null
  MODIFY updated_at BIGINT not null
```
