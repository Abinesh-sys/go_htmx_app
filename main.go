package main

import (
	"html/template"
	"io"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Contact struct {
	Name  string
	Email string
}

func newContact(name, email string) *Contact {
	return &Contact{
		Name:  name,
		Email: email,
	}
}

type contacts = []Contact

type Role struct {
	ID       int
	RoleName string
	Icon     string
}

type Data struct {
	Contacts []Contact
	Roles    []Role
}

func newData() *Data {
	return &Data{
		Contacts: []Contact{
			*newContact("Abinesh", "Abinesh@gmail.com"),
			*newContact("Bobby", "Bobby@gmail.com"),
		},
		Roles: []Role{
			{ID: 1, RoleName: "Admin", Icon: "/icons/admin.png"},
			{ID: 2, RoleName: "User", Icon: "/icons/user.png"},
			{ID: 3, RoleName: "Manager", Icon: "/icons/manager.png"},
		},
	}
}

func main() {
	e := echo.New()

	// Use middleware for logging
	e.Use(middleware.Logger())

	// Set the template renderer
	e.Renderer = newTemplate()

	// Sample data
	data := newData()

	// Handle the root route to render the page
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", data)
	})

	// Handle the POST request to add new contacts
	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")
		newContact := *newContact(name, email)
		data.Contacts = append(data.Contacts, newContact)
		return c.Render(http.StatusOK, "index", data)
	})

	// Serve static files for icons and CSS
	e.Static("/icons", "icons")
	e.Static("/css", "views/css")

	// Start the server
	e.Logger.Fatal(e.Start(":42069"))
}