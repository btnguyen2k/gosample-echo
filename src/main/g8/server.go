/*
A simple web application in Echo.

By Thanh Ba Nguyen @ btnguyen2k/gosample-echo.g8
*/
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func handlerHello(c echo.Context) error {
	output := `Hello, world!

Available URIs:

DELETE /delete
GET    /get
HEAD   /head
PATCH  /patch
POST   /post
PUT    /put
`
	return c.String(http.StatusOK, output)
}

func handlerDeleteGetHead(c echo.Context) error {
	html := "<p><a href=\"/\">Home</a></p>"
	html += "<table><tr><td>Method:</td><td>%s</td></tr>"
	html += "<tr><td nowrap='nowrap'>Request URI:</td><td>%s</td></tr>"
	html += "<tr><td nowrap='nowrap'>Request path:</td><td>%s</td></tr>"
	html += "<tr><td nowrap='nowrap'>Request host:</td><td>%s</td></tr>"
	html += "<tr><td nowrap='nowrap'>Query string:</td><td>%s</td></tr>"
	html += "<tr><td nowrap='nowrap'>Query params:</td><td>%s</td></tr>"
	html += "<tr><td nowrap='nowrap'>Headers:</td><td>%s</td></tr>"
	html += "<tr><td nowrap='nowrap'>Cookies:</td><td>%s</td></tr>"
	html += "</table>"
	cookies := ""
	for _, cookie := range c.Request().Cookies() {
		cookies += fmt.Sprintf("%s<br/>", cookie)
	}
	params := ""
	for k, _ := range c.QueryParams() {
		params += fmt.Sprintf("%s=%s<br/>", k, c.QueryParam(k))
	}
	html = fmt.Sprintf(html,
		c.Request().Method,                // method
		c.Request().RequestURI,            // URI
		c.Path(),                          // path
		c.Scheme()+"://"+c.Request().Host, // host
		c.QueryString(),                   // query string
		params,                            // query params
		c.Request().Header,                // headers
		cookies,                           // cookies
	)
	return c.HTML(http.StatusOK, html)
}

func handlerPatchPostPut(c echo.Context) error {
	defer c.Request().Body.Close()
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Error reading request body")
	}
	html := "<p><a href=\"/\">Home</a></p>"
	html += "<table><tr><td>Method:</td><td>%s</td></tr>"
	html += "<tr><td nowrap='nowrap'>Request URI:</td><td>%s</td></tr>"
	html += "<tr><td nowrap='nowrap'>Request path:</td><td>%s</td></tr>"
	html += "<tr><td nowrap='nowrap'>Request host:</td><td>%s</td></tr>"
	html += "<tr><td nowrap='nowrap'>Query string:</td><td>%s</td></tr>"
	html += "<tr><td nowrap='nowrap'>Query params:</td><td>%s</td></tr>"
	html += "<tr><td nowrap='nowrap'>Headers:</td><td>%s</td></tr>"
	html += "<tr><td nowrap='nowrap'>Cookies:</td><td>%s</td></tr>"
	html += "<tr><td nowrap='nowrap'>Body:</td><td>%s</td></tr>"
	html += "</table>"
	cookies := ""
	for _, cookie := range c.Request().Cookies() {
		cookies += fmt.Sprintf("%s<br/>", cookie)
	}
	params := ""
	for k, _ := range c.QueryParams() {
		params += fmt.Sprintf("%s=%s<br/>", k, c.QueryParam(k))
	}
	html = fmt.Sprintf(html,
		c.Request().Method,                // method
		c.Request().RequestURI,            // URI
		c.Path(),                          // path
		c.Scheme()+"://"+c.Request().Host, // host
		c.QueryString(),                   // query string
		params,                            // query params
		c.Request().Header,                // headers
		cookies,                           // cookies
		body,                              // body
	)
	return c.HTML(http.StatusOK, html)
}

func main() {
	e := echo.New()

	e.GET("/", handlerHello)
	e.DELETE("/delete", handlerDeleteGetHead)
	e.GET("/get", handlerDeleteGetHead)
	e.HEAD("/head", handlerDeleteGetHead)

	e.PATCH("/patch", handlerPatchPostPut)
	e.POST("/post", handlerPatchPostPut)
	e.PUT("/put", handlerPatchPostPut)

	const port = 8080                                 // listen port
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port))) // server will listen on 0.0.0.0:port
}
