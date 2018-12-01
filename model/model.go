package model
import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"log"
	"fmt"
	"crypto/md5"
)

var Config *conf
var DB *gorm.DB


type conf struct {
	HTTP_BIND string `yaml:"HTTP_BIND"`
	SECRET string `yaml:"SECRET"`
	EXPIRES int64 `yaml:"EXPIRES"`
	DATABASE string `yaml:"DATABASE"`
}

type Users struct {
	ID int 	`gorm:"primary_key:AUTO_INCREMENT"`
	Username string		`gorm:"type:varchar(36)"`
	Password string		`gorm:"type:varchar(36);not null;unique"`
	Email string				`gorm:"type:varchar(255);not null;unique"`
	Lastdate	time.Time
}

type Acls struct {
	ID int 	`gorm:"primary_key:AUTO_INCREMENT"`
	Domain string		`gorm:"type:varchar(100);not null;unique"`
	Url string			`gorm:"type:varchar(1024)"`
	Title string			`gorm:"type:varchar(1024)"`
	Keywords string		`gorm:"type:varchar(1024)"`
	Description string		`gorm:"type:varchar(1024)"`
	Favicon	string			`gorm:"type:varchar(1024)"`
	Method string			`gorm:"type:varchar(36)"`
	Comment string 		`gorm:"type:varchar(1024)"`
	Count string			`gorm:"type:varchar(36)"`
	Countid string			`gorm:"type:varchar(36)"`
	Username string		`gorm:"type:varchar(36)"`
}

func Convert(user, pass string)string{
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%s-DD", pass, user))))
}

func dbinit(){
	var isinit = false
	if _, err := os.Stat(Config.DATABASE);err != nil{
		isinit = true
	}
	db, err := gorm.Open("sqlite3", Config.DATABASE)
	if err != nil {
		log.Fatal("数据库初始化失败：", err.Error())
	}
	if isinit{
		db.AutoMigrate(&Users{}, &Acls{})
		user := Users{Username: "admin", Password: Convert("admin", "admin")}
		db.Create(&user)
		if ! db.NewRecord(user){
			log.Println("数据库初始化成功！")
		}
	}
	DB = db
}

func initConfig(){
	read,err_cfg := ioutil.ReadFile("config.yaml")
	if err_cfg != nil{
		log.Fatal("读取配置文件失败:",err_cfg.Error())
	}
	yamlerr := yaml.Unmarshal(read,&Config)
	if yamlerr != nil{
		log.Fatal("配置文件解析失败:",yamlerr.Error())
	}
}

func init(){
	initConfig()
	dbinit()
}