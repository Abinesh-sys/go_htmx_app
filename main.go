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
	ID     int
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
			{ID: 1, Name: "Abinesh", Email: "Abinesh@gmail.com", RoleID: 1},
			{ID: 2, Name: "Bobby", Email: "Bobby@gmail.com", RoleID: 2},
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

	// Display Contact Details
	e.GET("/contact/:id", func(c echo.Context) error {
		// Get the contact ID from the URL parameter
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid contact ID"})
		}

		// Find the contact and role by ID
		var contact Contact
		var role Role
		for _, c := range data.Contacts {
			if c.ID == id {
				contact = c
				break
			}
		}
		for _, r := range data.Roles {
			if r.ID == contact.RoleID {
				role = r
				break
			}
		}

		// Prepare the data for the detail page
		detailData := struct {
			Contact Contact
			Role    Role
		}{
			Contact: contact,
			Role:    role,
		}

		return c.Render(http.StatusOK, "contact", detailData)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")
		roleIDStr := c.FormValue("role")

		roleID, err := strconv.Atoi(roleIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role ID"})
		}

		newContact := Contact{ID: len(data.Contacts) + 1, Name: name, Email: email, RoleID: roleID}
		data.Contacts = append(data.Contacts, newContact)

		return c.Render(http.StatusOK, "index", data)
	})

	// Serve static files for icons and CSS
	e.Static("/css", "views/css")
	e.Static("/icons", "views/icons")

	// Start server
	e.Logger.Fatal(e.Start(":42069"))
}

