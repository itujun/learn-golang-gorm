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
