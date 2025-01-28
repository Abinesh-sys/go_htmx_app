package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Templates struct to handle rendering
type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	// Define custom template functions
	funcMap := template.FuncMap{
		"sub": func(a, b int) int { return a - b },
	}

	// Parse templates with custom functions
	return &Templates{
		templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("views/*.html")),
	}
}

type Contact struct {
	Name   string
	Email  string
	RoleID int
}

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
			{Name: "Abinesh", Email: "Abinesh@gmail.com", RoleID: 1},
			{Name: "Bobby", Email: "Bobby@gmail.com", RoleID: 2},
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
	e.Use(middleware.Logger())

	// Set the renderer
	e.Renderer = newTemplate()

	// Sample data
	data := newData()

	// Handlers
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", data)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")
		roleIDStr := c.FormValue("role")

		roleID, err := strconv.Atoi(roleIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role ID"})
		}

		newContact := Contact{Name: name, Email: email, RoleID: roleID}
		data.Contacts = append(data.Contacts, newContact)

		return c.Render(http.StatusOK, "index", data)
	})

	// Serve static files for icons and CSS
	e.Static("/css", "views/css")
	e.Static("/icons", "views/icons")

	// Start server
	e.Logger.Fatal(e.Start(":42069"))
}
