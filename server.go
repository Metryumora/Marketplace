package main

import (
	"net/http"
	_"github.com/jinzhu/gorm/dialects/postgres"
	"html/template"
	_ "image/png"
	. "Marketplace/persistence"
	"golang.org/x/crypto/bcrypt"
	"github.com/labstack/echo"
	"io"
)

var db = ConnectToDB()

type Template struct {
	templates *template.Template
}

var t = &Template{
	templates: template.Must(template.ParseGlob("html/*.html")),
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
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

func categoryHandler(c echo.Context) error {
	cat := c.QueryParam("category")
	return c.Render(http.StatusOK, "index.html", getMainPageData(cat))
}

func authUser(c echo.Context) error {
	username := c.FormValue("username")
	password := []byte(c.FormValue("pass"))

	var checkedUser User
	db.Where(&User{Username: username}).Find(&checkedUser)
	if &checkedUser == nil {
		return c.Render(http.StatusOK, "login.html", "Username not found!")
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(checkedUser.Password), password)
		if err != nil {
			return c.Render(http.StatusOK, "login.html", "Bad credentials!")
		} else {
			return c.Render(http.StatusOK, "index.html", getMainPageData(""))
		}
	}
}

func loginHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

func registerUser(c echo.Context) error {
	var user = new(User)
	user.Username = c.FormValue("username")
	hash, err := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), bcrypt.DefaultCost)
	Check(err)
	user.Password = string(hash)
	user.Email = c.FormValue("email")
	db.Debug().Save(user)
	c.Render(http.StatusOK, "login.html", "Registration successful!")
	return err
}
func registerHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", nil)
}

func aboutHandler(c echo.Context) error {
	productID := c.QueryParam("product")
	var data AboutProductData
	db.First(&data.Product, productID)
	db.Find(&data.Categories)
	return c.Render(http.StatusOK, "about.html", &data)
}

func main() {
	defer db.Close()
	//TestFillDB(db)
	e := echo.New()
	e.Renderer = t

	e.Static("/assets", "assets")
	e.GET("/", categoryHandler)
	e.GET("/login", loginHandler)
	e.GET("/register", registerHandler)
	e.GET("/about", aboutHandler)
	e.POST("/register", registerUser)
	e.POST("/login", authUser)
	e.Start(":8090")
}
