/*
Read/Write cookies in Echo.

By Thanh Ba Nguyen @ btnguyen2k/gosample-echo.g8
*/
package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func handlerHome(c echo.Context) error {
	html := `<h1>Echo sample: Cookie</h1>`
	html += `<style>table, th, td {border: 1px solid black;}</style>`
	html += `<table><thead><tr><th>Cookie</th><th>Value</th><th>Action</th></tr></thead>`
	html += `<tbody>`
	for _, cookie := range c.Request().Cookies() {
		html += "<tr><td>" + cookie.Name + "</td><td>" + cookie.Value + `</td><td><a href="/?a=delete&name=` + cookie.Name + `">Delete</a></td></tr>`
	}
	html += `</tbody></table>`

	html += `<form method="get" action="/">`
	html += `<p>Add/Set cookie:</p>`
	html += `<input type="hidden" name="a" value="add">`
	html += `Name <input name="name" type="text" style="width: 128px"> / `
	html += `Value <input name="value" type="text" style="width: 128px"> / `
	html += `<input type="submit" value="Add">`
	html += `</form>`
	return c.HTML(http.StatusOK, html)
}

func doAddCookie(c echo.Context) error {
	name := strings.TrimSpace(c.QueryParam("name"))
	if name != "" {
		value := strings.TrimSpace(c.QueryParam("value"))
		cookie := &http.Cookie{ // create new cookie with default expiry
			Name:  name,
			Value: value,
			Path:  "/",
		}
		c.SetCookie(cookie)
	}
	return c.Redirect(http.StatusFound, "/")
}

func doDeleteCookie(c echo.Context) error {
	name := strings.TrimSpace(c.QueryParam("name"))
	if name != "" {
		cookie := &http.Cookie{ // set max-age < 0 to remove cookie
			Name:   name,
			Path:   "/",
			MaxAge: -1,
		}
		c.SetCookie(cookie)
	}
	return c.Redirect(http.StatusFound, "/")
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		action := c.QueryParam("a")
		switch action {
		case "add":
			return doAddCookie(c)
		case "delete":
			return doDeleteCookie(c)
		default:
			return handlerHome(c)
		}
	})

	const port = 8080                                 // listen port
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port))) // server will listen on 0.0.0.0:port
}
