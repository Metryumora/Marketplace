package main

import (
	"net/http"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/postgres"
	"html/template"
	_ "image/png"
	"fmt"
	"io/ioutil"
	//"github.com/labstack/echo"
	//"github.com/labstack/echo/middleware"
)

type PageData struct {
	Categories        []Category
	NewProducts       []Product
	News              []News
	RequestedProducts []Product
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

type Product struct {
	gorm.Model
	Name        string
	Description string
	Number      string
	Price       string
	Image       string
	Category_id uint `gorm:"index"`
}

type News struct {
	gorm.Model
	Header string
	Text   string
}

type Sale struct {
	gorm.Model
	Product_id int
}

var db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres password=root2017 dbname=marketplace sslmode=disable")

func getMainPageData(category string) *PageData {

	var categories []Category
	db.Find(&categories)
	var news []News
	db.Limit(2).Find(&news)
	var latest []Product
	db.Limit(2).Find(&latest)
	var products []Product
	if len(category) > 0 {
		var c Category
		db.Where(&Category{Name: category}).First(&c)
		if c.ID != 0 {
			db.Where(&Product{Category_id: c.ID}).Find(&products)
		}
	} else {
		db.Find(&products)
	}
	data := new(PageData)
	data.Categories = categories
	data.News = news
	for i := 0; i < 2; i++ {
		latest[i].Description = latest[i].Description[:200]
	}
	data.NewProducts = latest
	data.RequestedProducts = products
	return data
}

//func createUser(c echo.Context) {
//	u := new(User)
//	u.Username = c.Param("name")
//	u.Password = c.Param("password")
//	u.Email = c.Param("email")
//}

func categoryHandler(w http.ResponseWriter, r *http.Request) {
	cat := r.URL.Query().Get("category")

	renderTemplate(w, "index", getMainPageData(cat))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register", nil)
}

func registerUserHandker(w http.ResponseWriter, r *http.Request) {
}

var templates = template.Must(template.ParseFiles("html/index.html", "html/login.html", "html/register.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data *PageData) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.DropTableIfExists(&User{}, &Sale{}, &Product{}, &Category{}, &News{})
	db.AutoMigrate(&User{}, &Category{}, &Product{}, &News{}, &Sale{})

	//// Create debug
	db.Create(&User{Username: "Admin", Email: "metryumora@gmail.com", Password: "pass_2017"})
	db.Create(&Category{Name: "Board games"})
	db.Create(&Category{Name: "Card games"})
	db.Create(&Category{Name: "Puzzles"})
	db.Create(&Category{Name: "Chess & Checkers"})
	db.Create(&Product{Name: "Mombasa", Price: "$41", Category_id: 1})
	db.Create(&Product{Name: "Scythe", Price: "$80", Category_id: 1})
	db.Create(&Product{Name: "Captain Sonar", Price: "$75", Category_id: 2})

	//e := echo.New()
	//
	//// Middleware
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	////CORS
	//e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//	AllowOrigins: []string{"*"},
	//	AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	//}))

	var products []Product
	db.Find(&products)
	for _, product := range products {
		description, err := ioutil.ReadFile("D:/Workspace/Marketplace/assets/products/info/" + fmt.Sprintf("%d", product.ID) + ".txt")
		check(err)
		db.Model(&product).Update("Description", description)
		db.Model(&product).Update("Image", "/assets/products/images/"+fmt.Sprintf("%d", product.ID)+".png")
	}

	db.Create(&News{Header: "Opening!", Text: "Our store is now opened!"})
	db.Create(&News{Header: "Breaking news2!", Text: "It's not really important, but you gonna read this anyway."})
	db.Create(&News{Header: "Breaking news3!", Text: "It's not really important, but you gonna read this anyway."})

	//for i := 0; i < 5; i++ {
	//	rand.Int()
	//	switch rand.Int() {
	//	case 0:
	//		db.Create(&Sale{Product_id: 1})
	//		break
	//	case 1:
	//		db.Create(&Sale{Product_id: 2})
	//		break
	//	case 2:
	//		db.Create(&Sale{Product_id: 3})
	//		break
	//	}
	//
	//}

	//e.POST("/register", )

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", categoryHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.ListenAndServe(":8080", nil)
}
