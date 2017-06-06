package main

import (
	"net/http"
	_"github.com/jinzhu/gorm/dialects/postgres"
	"html/template"
	_ "image/png"
	. "Marketplace/persistence"
	"golang.org/x/crypto/bcrypt"
)

var db = ConnectToDB()

var templates = template.Must(template.ParseFiles("html/index.html",
	"html/login.html",
	"html/register.html",
	"html/about.html"))

func renderTemplateWithPageData(w http.ResponseWriter, tmpl string, data *PageData) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderTemplateWithMessage(w http.ResponseWriter, tmpl string, message string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getMainPageData(category string) *PageData {
	var categories []Category
	db.Find(&categories)

	var news []News
	db.Limit(2).Find(&news)

	var latest []BoardGame
	db.Limit(2).Find(&latest)

	var products []BoardGame
	if len(category) > 0 {
		var c Category
		db.Where(&Category{Name: category}).First(&c)
		if c.ID != 0 {
			db.Where(&BoardGame{Category_id: c.ID}).Find(&products)
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
	data.NewBoardGames = latest
	data.RequestedBoardGames = products
	return data
}

func categoryHandler(w http.ResponseWriter, r *http.Request) {
	cat := r.URL.Query().Get("category")

	renderTemplateWithPageData(w, "index", getMainPageData(cat))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		params := r.Form
		username := params.Get("username")
		password := []byte(params.Get("pass"))

		var checkedUser User
		db.Where(&User{Username: username}).Find(&checkedUser)
		if &checkedUser == nil {
			renderTemplateWithMessage(w, "login", "Username not found!")
		} else {
			err := bcrypt.CompareHashAndPassword([]byte(checkedUser.Password), password)
			if err != nil {
				renderTemplateWithMessage(w, "login", "Bad credentials!")
			} else {
				renderTemplateWithPageData(w, "index", getMainPageData(""))
			}
		}

	}
	if r.Method == http.MethodGet {
		renderTemplateWithPageData(w, "login", nil)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		params := r.Form
		var user = new(User)
		user.Username = params.Get("username")
		hash, err := bcrypt.GenerateFromPassword([]byte(params.Get("password")), bcrypt.DefaultCost)
		Check(err)
		user.Password = string(hash)
		user.Email = params.Get("email")
		db.Debug().Save(user)
		renderTemplateWithMessage(w, "login", "Registration successful!")
	}
	if r.Method == http.MethodGet {
		renderTemplateWithPageData(w, "register", nil)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("product")
	var prod BoardGame
	db.First(&prod, productID)
	renderAboutTemplate(w, "about", &prod)
}

func renderAboutTemplate(w http.ResponseWriter, tmpl string, data *BoardGame) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	defer db.Close()
	//TestFillDB(db)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", categoryHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/about", aboutHandler)
	http.ListenAndServe(":8090", nil)
}
