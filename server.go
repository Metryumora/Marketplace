package main

import (
	"net/http"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/postgres"
	"html/template"
	_ "image/png"
	//"image"
	//"fmt"
)

type PageData struct {
	Categories      []Category
	PopularProducts []Product
	News            []News
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
	ID          uint `gorm:"primary_key"`
	Name        string
	Description string
	Number      string
	Price       string
	//Image
	Category_id uint `gorm:"index"`
}

type News struct {
	gorm.Model
	Header string
	Text   string
}

type Sales struct {
	gorm.Model
	Product_id int
}

var db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres password=root2017 dbname=marketplace sslmode=disable")

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var categories []Category
	db.Find(&categories)
	data := new(PageData)
	data.Categories = categories
	renderTemplate(w, "index", data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register", nil)
}

var templates = template.Must(template.ParseFiles("html/index.html", "html/login.html", "html/register.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data *PageData) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.DropTableIfExists(&User{}, &Sales{}, &Product{}, &Category{}, &News{})
	db.AutoMigrate(&User{}, &Category{}, &Product{}, &News{}, &Sales{})

	//// Create debug
	db.Create(&User{Username: "Admin", Email: "metryumora@gmail.com", Password: "pass_2017"})
	db.Create(&Category{Name: "Trash"})
	db.Create(&Category{Name: "Junk"})
	db.Create(&Product{Name: "Fancy-looking stuff", Description: "Short description", Price: "$19.99", Category_id: 1})
	db.Create(&Product{Name: "Fancy-looking stuff", Description: "Short description", Price: "$19.99", Category_id: 1})
	db.Create(&Product{Name: "Fancy-looking stuff", Description: "Short description", Price: "$19.99", Category_id: 2})

	//var products []Product
	//db.Where("category_id = ?", "1").Find(&products)
	//fmt.Printf("%s",products)
	//
	//
	//// Read
	//var user User
	//db.First(&user, 1)                       // find user with id 1
	//db.First(&user, "username = ?", "admin") // find user with code l1212

	// Update - update user's price to 2000
	//db.Model(&user).Update("Price", 2000)

	// Delete - delete user
	//db.Delete(&user)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.ListenAndServe(":8080", nil)
}
