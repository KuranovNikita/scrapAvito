package parse_kommersant

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Item struct {
	DateTime string
	Title    string
	Url      string
}

type ItemsArray struct {
	Items []Item
}

func ParseKommersant(url string) ItemsArray {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body.", err)
	}

	var itemsArray = ItemsArray{}

	document.Find("article.uho.rubric_lenta__item.js-article").Each(func(index int, element *goquery.Selection) {
		// класс для извлечения текста времени и даты
		dateTimeText := element.Find("p.uho__tag.rubric_lenta__item_tag.hide_desktop").Text()
		dateTime := strings.TrimSpace(dateTimeText)
		// класс для извлечения текста названия
		title := element.Find("h2.uho__name.rubric_lenta__item_name span.vam").Text()
		href, exists := element.Find(".uho__name.rubric_lenta__item_name a").Attr("href")
		if exists {
			href = strings.TrimSpace(href)
		}
		firstrune := []rune(href)[0]
		url = ""
		if firstrune == '/' {
			url = "https://www.kommersant.ru" + href
		} else {
			url = href
		}

		newItem := Item{DateTime: dateTime, Title: title, Url: url}

		itemsArray.Items = append(itemsArray.Items, newItem)
	})

	return itemsArray
}
