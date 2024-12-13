## Installation

1. go version: 1.16.7
2. require\
   These library will automatically install when run below command\
   github.com/go-resty/resty/v2 v2.7.0\
   github.com/go-sql-driver/mysql v1.6.0\
   github.com/gofiber/fiber/v2 v2.21.0\
   github.com/golang-jwt/jwt v3.2.2+incompatible\
   github.com/google/uuid v1.3.0 // indirect\
   github.com/spf13/cobra v1.2.1\
   github.com/spf13/viper v1.9.0\
   go.uber.org/zap v1.17.0\
   gorm.io/driver/mysql v1.1.3\
   gorm.io/gorm v1.22.2\
3. app.env file

## How to run

1. In the source project directory run this command to get dependency library

```
go mod download
```

or try this command

```
go get
```

2. Go mod tidy cleans up unused dependencies or adds missing dependencies

```
go mod tidy
```

3. Make sure in your GOPATH has those install library. Something like this\
   I has $GOPATH locate in this path

```

vmd@vmd:echo $GOPATH
/home/vmd/go
ls $GOPATH/pkg/mod/github.com
 ajg             coreos        go-check        go-stack          karrick         mergermarket   rogpeppe        stretchr
 andybalholm     cosiner       go-delve        grpc-ecosystem    kataras         microcosm-cc   russross        subosito
 aquasecurity    cpuguy83      gofiber         hashicorp         kkyr            mitchellh      ryanuber        teramoby
 armon           cweill        gofrs           haya14busa        klauspost       modern-go      satori          tmc
 aymerick        davecgh       gogo            imkira           '!knetic'        moul           schollz         twinj
 becheran        denisenkom    golang          inconshreveable   konsorten       natefinch      sergi           ugorji
 benbjohnson     dgrijalva     golang-jwt      iris-contrib      kr              nsf            shareed2k       uudashr
 beorn7          dgryski       golang-sql      jackal-xmpp       labstack        onokonem      '!shopify'       valyala
 bgentry         eknkc         go-mysql-org    jackc             lann            ortuman
```

4. Run Makefile command

```
make run
```

5. If you dont have "make" install you could run the following cmd. Step by step\
   Go will create executable file by run this cmd.

```
go build
```

6.  Then run the file created above with the `start` options. Your "app.env" file locate in root directory then run this cmd

```
./mine start
```

7.  If you want to specify "mine.env" in another path then run this, with --env flag that specify path-to-survey.env

```
./mine start --env=./mine.env
```

8. Your env file must have these key-value or in variable environment

| Key            |         Description |
| :------------- | ------------------: |
| DB_HOST        |      Host of the DB |
| DB_USERNAME    |         DB username |
| DB_PASSWORD    |         DB password |
| DB_DATABASE    | Name of DB instance |
| DB_PORT        |  Port connect to DB |
| USE_PRODUCTION |              0 or 1 |
