package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Service struct {
	Id          int
	Name        string
	Description string
	Image       string
}

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.LoadHTMLGlob("templates/**/*")

	// serviceData := []Service{
	// 	{Id: 1, Name: "Личная карточка", Description: "Добавить личную карточку", Image: "https://medal.spbarchives.ru/images/awarding-history/regulator.jpg"},
	// 	{Id: 2, Name: "Наградные", Description: "Добавить наградные", Image: "https://ordenov.net/upload/shop_1/5/9/7/item_5970/shop_property_file_5970_219.jpg"},
	// 	{Id: 3, Name: "Журнал боевых действий", Description: "Добавить журнал боевых действий", Image: "https://10otb.ru/content/dokuments/6gmk_scan/journal_bd_1944-03-1-31/0_oblozhka.jpg"},
	// }

	// словарь вместо массива
	serviceData := map[int]Service{
		1: {Id: 1, Name: "Личная карточка", Description: "Добавить личную карточку", Image: "https://medal.spbarchives.ru/images/awarding-history/regulator.jpg"},
		2: {Id: 2, Name: "Наградные", Description: "Добавить наградные", Image: "https://ordenov.net/upload/shop_1/5/9/7/item_5970/shop_property_file_5970_219.jpg"},
		3: {Id: 3, Name: "Журнал боевых действий", Description: "Журнал боевых действий — отчётно-информационный документ, входит в состав боевых документов.\nПри составлении (ведении) журнала боевых действий соблюдаются правила, предусмотренные уставами и наставлениями.\nВедётся в штабе объединения, соединения, воинской части, а также на кораблях 1, 2 и 3 ранга в течение всего времени нахождения в составе действующей армии или флота. В журнал боевых действий ежедневно заносятся сведения о подготовке и ходе боевых действий.", Image: "https://10otb.ru/content/dokuments/6gmk_scan/journal_bd_1944-03-1-31/0_oblozhka.jpg"},
		4: {Id: 4, Name: "Список раненых", Description: "Добавить список раненых", Image: "https://www.ww2.dk/ground/infanterie/infregt/1-IR-1.jpg"},
		5: {Id: 5, Name: "Список погибших", Description: "Добавить список погибших", Image: "https://www.ww2.dk/ground/infanterie/infregt/1-IR-1.jpg"},
		6: {Id: 6, Name: "Список пропавших без вести", Description: "Добавить список пропавших без вести", Image: "https://www.ww2.dk/ground/infanterie/infregt/1-IR-1.jpg"},
		7: {Id: 7, Name: "Список военнопленных", Description: "Добавить список военнопленных", Image: "https://www.ww2.dk/ground/infanterie/infregt/1-IR-1.jpg"},
	}

	r.GET("/", func(c *gin.Context) {
		search := strings.ToLower(c.Query("search"))

		filteredServiceData := []Service{}

		for _, service := range serviceData {
			service_name := strings.ToLower(service.Name)
			if search == "" || strings.Contains(service_name, search) {
				filteredServiceData = append(filteredServiceData, service)
			}
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":    "Книга памяти",
			"services": filteredServiceData,
		})
	})

	r.GET("/service/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.String(http.StatusBadRequest, "Bad request")
			return
		}

		service, ok := serviceData[id]

		if !ok {
			c.String(http.StatusNotFound, "Not found")
			return
		}

		c.HTML(http.StatusOK, "service.html", gin.H{
			"title":   "Список сервисов",
			"service": service,
		})
	})

	r.Static("/static", "./static/")

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
