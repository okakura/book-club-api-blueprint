package bookclub

import (
	// "cgontext"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Votes  int    `json:"votes"`
}

type Club struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Members    []string `json:"members"`
	Proposals  []Book   `json:"proposals"`
	CurrentBook *Book   `json:"current_book,omitempty"`
}

var (
	clubs = map[string]*Club{}
	mu    = sync.Mutex{}
)

func Main() {
	InitDB()
	
	r := gin.Default()

	r.GET("/clubs", listClubs)
	r.POST("/clubs", createClub)
	r.GET("/clubs/:id", getClub)
	r.POST("/clubs/:id/members", addMember)
	r.DELETE("/clubs/:id/members/:name", removeMember)
	r.GET("/clubs/:id/books", listBooks)
	r.POST("/clubs/:id/books", proposeBook)
	r.POST("/clubs/:id/books/:bookID/vote", voteBook)
	r.GET("/clubs/:id/current", getCurrentBook)
	r.GET("/openapi.yaml", func(c *gin.Context) {
		c.File("./docs/openapi.yaml")
	})
	r.Static("/docs", "./docs")


	r.Run(":8080")
}

func listClubs(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	res := make([]*Club, 0, len(clubs))
	for _, club := range clubs {
		res = append(res, club)
	}
	c.JSON(http.StatusOK, res)
}

func createClub(c *gin.Context) {
	type req struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	clubs[r.ID] = &Club{ID: r.ID, Name: r.Name, Members: []string{}, Proposals: []Book{}}
	c.JSON(http.StatusCreated, clubs[r.ID])
}

func getClub(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	id := c.Param("id")
	club, ok := clubs[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		return
	}
	c.JSON(http.StatusOK, club)
}

func addMember(c *gin.Context) {
	type req struct {
		Name string `json:"name"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	club := clubs[c.Param("id")]
	club.Members = append(club.Members, r.Name)
	c.JSON(http.StatusOK, club)
}

func removeMember(c *gin.Context) {
	id := c.Param("id")
	name := c.Param("name")
	mu.Lock()
	defer mu.Unlock()
	club := clubs[id]
	filtered := []string{}
	for _, m := range club.Members {
		if m != name {
			filtered = append(filtered, m)
		}
	}
	club.Members = filtered
	c.JSON(http.StatusOK, club)
}

func listBooks(c *gin.Context) {
	club := clubs[c.Param("id")]
	c.JSON(http.StatusOK, club.Proposals)
}

func proposeBook(c *gin.Context) {
	type req struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Author string `json:"author"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	club := clubs[c.Param("id")]
	club.Proposals = append(club.Proposals, Book{ID: r.ID, Title: r.Title, Author: r.Author, Votes: 0})
	c.JSON(http.StatusOK, club.Proposals)
}

func voteBook(c *gin.Context) {
	club := clubs[c.Param("id")]
	bookID := c.Param("bookID")
	mu.Lock()
	defer mu.Unlock()
	for i, b := range club.Proposals {
		if b.ID == bookID {
			club.Proposals[i].Votes++
			break
		}
	}
	c.JSON(http.StatusOK, club.Proposals)
}

func getCurrentBook(c *gin.Context) {
	club := clubs[c.Param("id")]
	if len(club.Proposals) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No books proposed yet"})
		return
	}
	var top *Book
	for i := range club.Proposals {
		if top == nil || club.Proposals[i].Votes > top.Votes {
			top = &club.Proposals[i]
		}
	}
	club.CurrentBook = top
	c.JSON(http.StatusOK, top)
}
