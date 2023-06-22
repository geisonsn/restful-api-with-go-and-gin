// A simple CRUD (Create, Read, Update and Delete) example using Gin Gonic

package main

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type AlbumSlice []album

func (a AlbumSlice) Len() int {
	return len(a)
}

func (a AlbumSlice) Less(i, j int) bool {
	return a[i].ID < a[j].ID
}

func (a AlbumSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

var albums = AlbumSlice{}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	newAlbum, err := fillNewAlbumFromRequest(c)
	if err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	sort.Sort(albums)
	c.IndentedJSON(http.StatusOK, albums)
}

// gertAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// updateAlbums updates an album whose ID value matches the id parameter sent by the client.
func updateAlbums(c *gin.Context) {
	id := c.Param("id")

	var updatedAlbums []album
	for _, a := range albums { // TODO: improve performance by using a search algorithm
		if a.ID != id {
			updatedAlbums = append(updatedAlbums, a)
		} else {
			updatedAlbum, err := fillNewAlbumFromRequest(c)
			if err != nil {
				return
			}
			updatedAlbum.ID = id
			updatedAlbums = append(updatedAlbums, updatedAlbum)
		}
	}

	albums = updatedAlbums
}

// deleteAlbums deletes the album whose ID value matches the id parameter sent by the client.
func deleteAlbums(c *gin.Context) {
	id := c.Param("id")

	var newAlbumSlice []album
	for _, a := range albums {
		if a.ID != id {
			newAlbumSlice = append(newAlbumSlice, a)
		}
	}
	albums = newAlbumSlice
	c.IndentedJSON(http.StatusOK, albums)
}

// fillNewAlbumFromRequest create a new 'album' and fill its content from the request body
func fillNewAlbumFromRequest(c *gin.Context) (album, error) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to new Album.
	err := c.BindJSON(&newAlbum)
	if err != nil {
		return album{}, nil
	}

	return newAlbum, err
}

func init() {
	albums = AlbumSlice{
		{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.PUT("/albums/:id", updateAlbums)
	router.DELETE("/albums/:id", deleteAlbums)

	router.Run("localhost:8080")
}
