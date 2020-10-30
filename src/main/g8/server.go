/*
Work with sessions in Echo.

By Thanh Ba Nguyen @ btnguyen2k/gosample-echo.g8

Summary: similar to cookies, working with sessions in Echo is simple. There is, however, one important point
worth to mention: server is responsible for managing session storage. Echo supports a middleware architecture
that can facilitates different session management implementations
(https://github.com/gorilla/sessions#store-implementations). Each implementation has its own pros and cons.
The default one provides cookie and filesystem based session store, which we will cover in this sample.

	- Filesystem based session store: good performance, can store big (in size) sessions, but hard to scale beyond
      single server.
	- Cookie bases session store: excellent option for scaling and performance, has a limit in size (around 3Kb).
      One may have security concern as session is stored in cookie - which is managed by user's browser.
      To address the size limit issue, session data can be compressed. However, it may further raise concern regarding
      security.
*/
package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"main/cocostore"
)

const sessionStoreName = "my_sess_store"

func handlerHome(c echo.Context) error {
	html := `<h1>Echo sample: Session</h1>`
	html += `<style>table, th, td {border: 1px solid black;}</style>`
	html += `<table><thead><tr><th>Session</th><th>Value</th><th>Action</th></tr></thead>`
	html += `<tbody>`
	sess, _ := session.Get(sessionStoreName, c)
	for sKey, sVal := range sess.Values {
		html += fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td><a href="/?a=delete&name=%s">Delete</a></td></tr>`, sKey, sVal, sKey)
	}
	html += `</tbody></table>`

	html += `<form method="get" action="/">`
	html += `<p>Add/Set session:</p>`
	html += `<input type="hidden" name="a" value="add">`
	html += `Name <input name="name" type="text" style="width: 128px"> / `
	html += `Value <input name="value" type="text" style="width: 128px"> / `
	html += `<input type="submit" value="Add">`
	html += `</form>`
	return c.HTML(http.StatusOK, html)
}

func doAddSession(c echo.Context) error {
	name := strings.TrimSpace(c.QueryParam("name"))
	if name != "" {
		value := strings.TrimSpace(c.QueryParam("value"))
		sess, _ := session.Get(sessionStoreName, c)
		sess.Values[name] = value
		sess.Save(c.Request(), c.Response())
	}
	return c.Redirect(http.StatusFound, "/")
}

func doDeleteSession(c echo.Context) error {
	name := strings.TrimSpace(c.QueryParam("name"))
	if name != "" {
		sess, _ := session.Get(sessionStoreName, c)
		delete(sess.Values, name)
		sess.Save(c.Request(), c.Response())
	}
	return c.Redirect(http.StatusFound, "/")
}

func main() {
	e := echo.New()

	// session store secret must be secured!
	sessionStoreSecret := make([]byte, 16)
	rand.Read(sessionStoreSecret)

	{
		// // use cookie based session store:
		// e.Use(session.Middleware(sessions.NewCookieStore(sessionStoreSecret)))

		// use cookie based session store with max compression:
		e.Use(session.Middleware(cocostore.NewCompressedCookieStore(sessionStoreSecret).SetCompressionLevel(9)))

		// // use filesystem based session store:
		// const sessionStorePath = "./temp/cookies"
		// if err := os.MkdirAll(sessionStorePath, 0700); err != nil {
		// 	panic(err)
		// }
		// e.Use(session.Middleware(sessions.NewFilesystemStore(sessionStorePath, sessionStoreSecret)))

		// other session store implementations can be found at https://github.com/gorilla/sessions#store-implementations
	}

	e.GET("/", func(c echo.Context) error {
		action := c.QueryParam("a")
		switch action {
		case "add":
			return doAddSession(c)
		case "delete":
			return doDeleteSession(c)
		default:
			return handlerHome(c)
		}
	})

	const port = 8080                                 // listen port
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port))) // server will listen on 0.0.0.0:port
}
