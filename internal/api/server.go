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

	// словарь вместо массива
	serviceData := map[int]Service{
		1: {Id: 1, Name: "Личная карточка", Description: "Личная карточка — документ, содержащий информацию о военнослужащем. В личной карточке заносятся данные о биографии, физической подготовке, участии в боевых действиях и прохождении службы.", Image: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQX6qztHnPlLg_dBb-LoJIkqWjEWRvKrZDAQw&usqp=CAU"},
		2: {Id: 2, Name: "Наградные", Description: "Наградные документы — документы, удостоверяющие присуждение различных наград военнослужащим за особые заслуги и подвиги. В наградные документы вносятся сведения о военнослужащих, обстоятельствах присуждения и виде награды.", Image: "https://www.ulspu.ru/vov/img/otdelno/pushkareva_g_v/big/lichnaya_kartochka.jpg"},
		3: {Id: 3, Name: "Журнал боевых действий", Description: "Журнал боевых действий — отчётно-информационный документ, входит в состав боевых документов.\nПри составлении (ведении) журнала боевых действий соблюдаются правила, предусмотренные уставами и наставлениями.\nВедётся в штабе объединения, соединения, воинской части, а также на кораблях 1, 2 и 3 ранга в течение всего времени нахождения в составе действующей армии или флота. В журнал боевых действий ежедневно заносятся сведения о подготовке и ходе боевых действий.", Image: "https://10otb.ru/content/dokuments/6gmk_scan/journal_bd_1944-03-1-31/0_oblozhka.jpg"},
		4: {Id: 4, Name: "Список раненых", Description: "Список раненых — документ, содержащий информацию о военнослужащих, получивших ранения во время боевых действий.", Image: "https://www.prlib.ru/sites/default/files/book_preview/d48d307a-f81b-495b-a6b8-6e7192cfa3bb/12976930_doc1.jpg"},
		5: {Id: 5, Name: "Список погибших", Description: "Список погибших — документ, содержащий информацию о военнослужащих, погибших в результате военных действий. Ведется в штабе объединения, соединения, воинской части", Image: "https://ic.pics.livejournal.com/vagante_travel/50304595/19851/19851_original.jpg"},
		7: {Id: 7, Name: "Список военнопленных", Description: "База данных советских военнопленных» содержит персональные данные о советских военнопленных во время Второй мировой войны, находившихся под стражей в лагерях для военнопленных или в трудовых командах на территории бывшего германского рейха. Обнародованная база данных включает в себя базовую информацию: имена и фамилии военнопленных, дату рождения и дату смерти.", Image: "https://russian7.ru/wp-content/uploads/2015/02/1678_434498865_big.jpg"},
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
			"search":   search,
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
			"title":   service.Name,
			"service": service,
		})
	})

	r.Static("/static", "./static/")

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
