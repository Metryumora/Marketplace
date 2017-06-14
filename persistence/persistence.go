package persistence

import (
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

type PageData struct {
	Categories          []Category
	NewBoardGames       []BoardGame
	News                []News
	RequestedBoardGames []BoardGame
	Token               json.Token
}

type AboutProductData struct {
	Categories []Category
	Product    BoardGame
}

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
}

type Category struct {
	ID   uint `gorm:"primary_key"`
	Name string
}

type BoardGame struct {
	gorm.Model
	Name        string
	Description string
	Number      string
	Price       string
	Image       string
	Category_id uint `gorm:"index"`
	Ages        string
	Players     string
}

type News struct {
	gorm.Model
	Header string
	Text   string
}

type Sale struct {
	gorm.Model
	BoardGame_id int
}

func ConnectToDB() *gorm.DB {
	db, err := gorm.Open("postgres", "postgres",
		"host=localhost"+
			" port=5432"+
			" user=postgres"+
			" password=root2017"+
			" dbname=marketplace"+
			" sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func TestFillDB(db *gorm.DB) {

	// Migrate the schema
	db.DropTableIfExists(&User{}, &Sale{}, &BoardGame{}, &Category{}, &News{})
	db.AutoMigrate(&User{}, &Category{}, &BoardGame{}, &News{}, &Sale{})

	// Filling tables
	hash, err := bcrypt.GenerateFromPassword([]byte("Password1"), bcrypt.DefaultCost)
	Check(err)
	db.Create(&User{Username: "Metr_yumora", Email: "metryumors@gmail.com", Password: string(hash)})

	db.Create(&Category{Name: "Board games"})
	db.Create(&Category{Name: "Card games"})
	db.Create(&Category{Name: "Puzzles"})
	db.Create(&Category{Name: "Chess & Checkers"})
	db.Create(&BoardGame{Name: "Mombasa", Price: "$41", Category_id: 1, Players: "2-4", Ages: "12+"})
	db.Create(&BoardGame{Name: "Scythe", Price: "$80", Category_id: 1, Players: "2-5", Ages: "-"})
	db.Create(&BoardGame{Name: "Captain Sonar", Price: "$75", Category_id: 2, Players: "2-8", Ages: "14+"})

	var products []BoardGame
	db.Find(&products)
	for _, product := range products {
		description, err := ioutil.ReadFile("./assets/products/info/" + fmt.Sprintf("%d", product.ID) + ".txt")
		db.Model(&product).Update("Description", description)
		db.Model(&product).Update("Image", "./assets/products/images/"+fmt.Sprintf("%d", product.ID)+".png")
		Check(err)
	}

	db.Create(&News{Header: "Opening!", Text: "Our store is now opened, be the first one to buy!"})
	db.Create(&News{Header: "Breaking news2!", Text: "It's not really important, but you gonna read this anyway."})
	db.Create(&News{Header: "Breaking news3!", Text: "It's not really important, but you gonna read this anyway."})
}
