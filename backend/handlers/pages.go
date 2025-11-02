package handlers

import "github.com/gin-gonic/gin"

func LoginPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "login",
	})
}

func RegisterPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "register",
	})
}
func ProductsPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "products",
	})
}

func IndexPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "index",
	})
}
