package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Document struct {
	Document_id int
	Title       string
	Description string
	Image_url   string
	Status      string
}

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.LoadHTMLGlob("templates/**/*")

	documentsData := []Document{
		{Document_id: 0, Title: "Орден Отечественной войны II степени", Description: "Орден Отечественной войны II степени — советское государственное наградное учреждение, учреждённое Указом Президиума Верховного Совета СССР от 20 мая 1942 года «За установление и укрепление военного порядка, за отличие в обороне Советского Союза и в освобождении его территорий от немецко-фашистских захватчиков».", Image_url: "https://upload.wikimedia.org/wikipedia/commons/1/1c/Order_of_the_Patriotic_War_%281st_class%29.png", Status: "active"},
		{Document_id: 1, Title: "ЖБД 34-й стрелковой дивизии", Description: "Журнал боевых действий 34-й стрелковой дивизии", Image_url: "https://sun9-61.userapi.com/impf/1_k_ogYH4PDdygGmvhKaYvlD2_zzlzLkO3pzqg/CSPpSw6FUoc.jpg?size=1280x940&quality=96&sign=f84827e9b86a1635ba031b60d6615d85&type=album", Status: "active"},
		{Document_id: 2, Title: "Оперативный отдел штаба 49 армии", Description: "Отчёт о боевых действиях на Западном направлении. Обеспечение операции материальными ресурсами, инженерным обеспечением, потери в операции", Image_url: "https://podvignaroda.ru/filter/filterimage?path=TV/001/208-0002511-1042/00000167.jpg&id=60332603&id1=a97e93049bac29beadfad6636e96ba21", Status: "active"},
		{Document_id: 3, Title: "15.12.1941 Штаб ВПУ Юго-Западного фронта", Description: "Боевые сводки потерь противника.", Image_url: "https://podvignaroda.ru/filter/filterimage?path=TV/001/251-0000646-0051/00000315.jpg&id=60113134&id1=7db41e7f93af57d1cd6f17add3648de3", Status: "active"},
		{Document_id: 4, Title: "Орден Красного Знамени", Description: "Орден Красного Знамени — советское государственное наградное учреждение, учреждённое Указом Президиума Верховного Совета СССР от 6 ноября 1943 года «За отличие в бою и в других боевых операциях в защите Советского Союза и в освобождении его территорий от немецко-фашистских захватчиков».", Image_url: "https://podvignaroda.ru/img/awards/new/Orden_Krasnogo_Znameni_1st.png", Status: "active"},
	}

	r.GET("/", func(c *gin.Context) {
		searchDocument := strings.ToLower(c.Query("document"))

		filteredServiceData := []Document{}

		for _, service := range documentsData {
			service_name := strings.ToLower(service.Title)
			if searchDocument == "" || strings.Contains(service_name, searchDocument) {
				filteredServiceData = append(filteredServiceData, service)
			}
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":     "Книга памяти",
			"search":    c.Query("document"),
			"documents": filteredServiceData,
		})
	})

	r.GET("/document/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.String(http.StatusBadRequest, "Bad request")
			return
		}

		service := Document{}

		for _, s := range documentsData {
			if s.Document_id == id {
				service = s
				break
			}
		}

		if service != (Document{}) {
			c.HTML(http.StatusOK, "document.html", gin.H{
				"title":    service.Title,
				"document": service,
			})
		} else {
			c.HTML(http.StatusOK, "404.html", gin.H{})
		}
	})

	r.Static("/static", "./static/")

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
