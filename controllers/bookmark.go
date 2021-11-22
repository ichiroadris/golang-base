package controllers

import (
	"github.com/gin-gonic/gin"
	"golang-gogin/forms"
	"golang-gogin/helpers"
	"golang-gogin/models"
	"golang-gogin/services"
	"log"
	"time"
)

var BookmarkModel = new(models.BookmarkModel)

type BookmarkController struct{}

func responseWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"message": message})
}

func (b *BookmarkController) FetchBookmarks(c *gin.Context) {
	user := c.MustGet("User").(models.User)

	if user.Email == "" {
		responseWithError(c, 404, "Please login")
		return
	}

	var linkModel models.BookmarkModel

	results, err := linkModel.FetchBookmarks(user.ID)

	if err != nil {
		responseWithError(c, 500, "Problem fetching your articles")
		return
	}

	if results != nil {
		c.JSON(200, gin.H{"bookmarks": results})
	} else {
		c.JSON(200, gin.H{"bookmarks": []string{}})
	}
}

// DeleteBookmark controller handles deleting a single bookmark
func (b *BookmarkController) DeleteBookmark(c *gin.Context) {
	// Get the user we set in the Authenticate middleware
	user := c.MustGet("User").(models.User)

	// Check if we have user and respond if not
	if user.Email == "" {
		responseWithError(c, 404, "Please login")
		return
	}

	// Get query parameter value that holds the bookmakr id
	bookmarkID, found := c.GetQuery("bookmark_id")

	// Check if the query parameter `bookmark_id` is provided, respond if not
	if !found {
		responseWithError(c, 400, "Please provide bookmark id")
		return
	}

	// Define variable to hold the model methods
	var linkModel models.BookmarkModel

	// Delete record
	err := linkModel.DeleteBookmark(bookmarkID)

	// Check if we got an error while deleting the file
	if err != nil {
		responseWithError(c, 500, "Problem deleting bookmark")
		return
	}

	// Respond with a 204 No Content on successful delete
	c.JSON(204, gin.H{"message": "Deleted bookmark successfully"})
}

// CreateBookmak controller handles creating a bookmark of a specifi user
func (b *BookmarkController) CreateBookmak(c *gin.Context) {
	// Get the user we set in the Authenticate middleware
	user := c.MustGet("User").(models.User)

	// Check if we have user and respond if not
	if user.Email == "" {
		responseWithError(c, 404, "Please login")
		return
	}

	// Define variable to hold the payload structure
	var data forms.BookmarkPayload

	// Check if required fields are provided
	if c.BindJSON(&data) != nil {
		log.Fatal(c.BindJSON(&data))
		responseWithError(c, 406, "Please provide link, and name")
		return
	}

	// Define variable to hold the model methods
	var linkModel models.BookmarkModel

	// Check if the url is valid and respond if its not
	if !helpers.IsValidURL(data.Link) {
		responseWithError(c, 400, "Link is invalid")
	}

	// Define a variable to hold out scrapper methods
	var scrapper services.Scrapper

	// Make a website request based on the link provided on the request body
	var meta services.Meta = scrapper.CallWebsite(data.Link, c)

	// Define the data to be save to the database
	var bookmarkPayload models.Link = models.Link{
		Name:            data.Name,
		MetaImage:       meta.Image,
		MetaDescription: meta.Description,
		MetaSite:        meta.Site,
		MetaURL:         meta.URL,
		Link:            data.Link,
		Owner:           user.ID,
		CreateAt:        time.Now(),
	}

	// Save the bookmakr
	err := linkModel.CreateBookmark(bookmarkPayload)

	// Check if we got an error while saving and respond if we do
	if err != nil {
		responseWithError(c, 500, "Problem saving your bookmark")
		log.Fatal(err)
		return
	}

	// Return a 201 Created on successful creation with the bookmark saved
	c.JSON(201, gin.H{"message": "Bookmark saved", "bookmark": bookmarkPayload})
}
