package services

import (
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

type Scrapper struct{}

type Meta struct {
	Image, Description, URL, Title, Site string
}

func (scrapper *Scrapper) CallWebsite(websiteURL string, c *gin.Context) Meta {
	var meta Meta = Meta{
		Image:       "",
		Description: "",
		URL:         "",
		Title:       "",
		Site:        "",
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	request, err := http.NewRequest("GET", websiteURL, nil)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err})
	}

	request.Header.Set("pragma", "no-cache")

	request.Header.Set("cache-control", "no-cache")

	request.Header.Set("dnt", "1")

	request.Header.Set("upgrade-insecure-requests", "1")

	request.Header.Set("referer", websiteURL)

	resp, err := client.Do(request)

	if resp.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(resp.Body)

		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"message": err})
		}

		doc.Find("meta").Each(func(i int, s *goquery.Selection) {
			metaProperty, _ := s.Attr("property")
			metaContent, _ := s.Attr("content")

			// If we happen to get any of the two then assign the Site attribute for Meta
			if metaProperty == "og:site_name" || metaProperty == "twitter:site" {
				meta.Site = metaContent
			}

			// If we happen to get any of the two then assign the URL attribute for Meta
			if metaProperty == "og:url" {
				meta.URL = metaContent
			}
			// If we happen to get any of the two then assign the Image attribute for Meta
			if metaProperty == "og:image" || metaProperty == "twitter:image" {
				meta.Image = metaContent
			}

			// If we happen to get any of the two then assign the Title attribute for Meta
			if metaProperty == "og:title" || metaProperty == "twitter:title" {
				meta.Title = metaContent
			}

			// If we happen to get any of the two then assign the Description attribute for Meta
			if metaProperty == "og:description" || metaProperty == "twitter:description" {
				meta.Description = metaContent
			}

		})
	}

	return meta
}
