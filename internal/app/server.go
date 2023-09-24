package app

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"github.com/gin-gonic/gin"
)

type Service struct {
	Id          int
	Name        string
	Description string
	Image       string
}

func (a *Application) StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.LoadHTMLGlob("templates/**/*")

	r.GET("/", func(c *gin.Context) {

		search := strings.ToLower(c.Query("search"))

		documents, err := a.repository.GetDocuments()
		if err != nil {
			log.Println("Error with running\nServer down")
			return
		}

		var filteredDocuments []*ds.Document

		for _, document := range documents {
			document_name := strings.ToLower(document.Title)
			if search == "" || strings.Contains(document_name, search) {
				filteredDocuments = append(filteredDocuments, document)
			}
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":     "Книга памяти",
			"search":    search,
			"documents": filteredDocuments,
		})
	})

	r.GET("/document/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.String(http.StatusBadRequest, "Bad request")
			return
		}

		document, err := a.repository.GetDocumentByID(id)
		if err != nil {
			log.Println("Error with running\nServer down")
			return
		}

		if err != nil {
			c.String(http.StatusNotFound, "Not found")
			return
		}

		c.HTML(http.StatusOK, "document.html", gin.H{
			"title":    document.Title,
			"document": document,
		})
	})

	r.Static("/static", "./static/")

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
