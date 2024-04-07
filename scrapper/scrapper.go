package scrapper

import (
	"fmt"
	"jiva-guildes/domain/models/guilde"
	"jiva-guildes/settings"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/gocolly/colly"
)

func connectionAvailable() bool {
	cmd := exec.Command("ping", "-c", "1", "google.com")
	err := cmd.Run()
	return err == nil
}
func get_url_from_style(url string) string {
	regexPattern := `\(([^)]+)\)`
	regex := regexp.MustCompile(regexPattern)
	match := regex.FindStringSubmatch(url)
	if len(match) > 1 {
		return match[1]
	} else {
		return url
	}
}
func Scrap() {
	if !connectionAvailable() {
		println("No connection available for scrapping, try later...")
		return
	}
	var guildes []guilde.Guilde
	c := colly.NewCollector()
	c.OnHTML(".ak-bg-odd", func(e *colly.HTMLElement) {
		var g guilde.Guilde
		img_url := get_url_from_style(e.ChildAttr(".ak-emblem", "style"))
		g_name := e.ChildText("td:nth-child(2)")
		g_page := "https://www.dofus.com/" + e.ChildAttr("td:nth-child(2) > a", "href")
		g.Img_url = img_url
		g.Name = g_name
		g.Page_url = g_page
		guildes = append(guildes, g)

	})
	MAX_PAGE := 228
	page := 1

	for i := page; i < MAX_PAGE+1; i++ {
		link := fmt.Sprintf("https://www.dofus.com/fr/mmorpg/communaute/annuaires/pages-guildes?name=&server_id=292&page=%v", i)
		c.Visit(link)
		fmt.Println(fmt.Sprint(len(guildes)) + " guilde found.")
	}
	file, err := os.Create(settings.AppSettings.CSV_FILE_FROM_SCRAP)
	if err != nil {
		log.Fatal("ERROR! ", err)
	}
	for _, v := range guildes {
		str_line := fmt.Sprintf("%s,%s,%s\n", v.Name, v.Img_url, v.Page_url)
		file.WriteString(str_line)
	}
	file.Close()
}
