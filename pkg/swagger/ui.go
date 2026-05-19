package swagger

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterScalar(r *gin.Engine, title string, specURL string) {
	r.GET("/docs", func(c *gin.Context) {
		html := `
		<!doctype html>
		<html>
		  <head>
			<title>` + title + `</title>
			<meta charset="utf-8" />
			<meta name="viewport" content="width=device-width, initial-scale=1" />
			<style>body { margin: 0; }</style>
		  </head>
		  <body>
			<script
			  id="api-reference"
			  data-url="` + specURL + `"></script>
			<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
		  </body>
		</html>`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

	r.StaticFile(specURL, "./docs/swagger.json")
}
