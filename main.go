package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	_ "one-to-many/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Models for User and Post
type User struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name" gorm:"type:varchar(100);not null"`
	Email string `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Posts []Post `json:"posts" gorm:"foreignKey:UserID"`
}

type Post struct {
	ID     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"user_id"`
	Tags   []Tag  `json:"tags" gorm:"many2many:post_tags;"`
}

// @description Tag model
type Tag struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name" gorm:"type:varchar(100);not null"`
	Posts []Post `json:"posts" gorm:"many2many:post_tags;"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Global database variable
var db *gorm.DB
var err error

// Initialize the DB connection
func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}
	db.AutoMigrate(&User{}, &Post{}) // Auto-migrate the tables
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{
		Code:    code,
		Message: message,
	})
}

// @Summary Get all users
// @Description Get all users along with their posts
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {array} User
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users [get]
// Get all users with their posts
func getUsers(c *gin.Context) {
	var users []User
	if err := db.Preload("Posts.Tags").Find(&users).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching users")
		return
	}
	c.JSON(200, users)
}

// @Summary Get user by ID
// @Description Get a specific user with their posts
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID" // The ID of the user to retrieve
// @Success 200 {object} User // The user object returned in the response
// @Failure 400 {object} ErrorResponse // Bad request if the ID is invalid
// @Failure 404 {object} ErrorResponse // User not found
// @Failure 500 {object} ErrorResponse // Internal server error
// @Router /api/v1/users/{id} [get]
// Get a specific user with their posts
func getUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := db.Preload("Posts.Tags").First(&user, id).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "User not found")
		return
	}
	c.JSON(200, user)
}

// @Summary Create a new user
// @Description Create a new user by providing a name and email
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body User true "New user information"
// @Success 201 {object} User
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users [post]
// Create a new user
func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid input")
		return
	}
	if err := db.Create(&user).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create user")
		return
	}
	c.JSON(201, user)
}

// @Summary Create a post for a specific user
// @Description Add a new post to the specified user by user ID
// @Tags posts
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param post body Post true "Post details"
// @Success 201 {object} Post
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{userID}/posts [post]
// Create a new post for a specific user
func createPost(c *gin.Context) {
	userID := c.Param("userID")

	// Convert userID from string to int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid input")
		return
	}
	// Now assign the integer value to post.UserID
	post.UserID = userIDInt

	// Create the post
	if err := db.Create(&post).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create post")
		return
	}
	c.JSON(201, post)
}

// Handlers for Tags
// @Summary Get all tags
// @Description Get all tags along with associated posts
// @Tags Tags
// @Produce json
// @Success 200 {array} Tag
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/tags [get]
func getTags(c *gin.Context) {
	var tags []Tag
	if err := db.Preload("Posts").Find(&tags).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching tags")
		return
	}
	c.JSON(http.StatusOK, tags)
}

// @Summary Create a new tag
// @Description Create a new tag
// @Tags Tags
// @Accept json
// @Produce json
// @Param tag body Tag true "Tag data"
// @Success 201 {object} Tag
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/tags [post]
func createTag(c *gin.Context) {
	var tag Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid input")
		return
	}
	if err := db.Create(&tag).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create tag")
		return
	}
	c.JSON(http.StatusCreated, tag)
}

// Associate a tag with a post
// @Summary Add a tag to a post
// @Description Associate an existing tag with a specific post
// @Tags Tags
// @Param postID path int true "Post ID"
// @Param tagID path int true "Tag ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/posts/{postID}/tags/{tagID} [post]
func addTagToPost(c *gin.Context) {
	postID := c.Param("postID")
	tagID := c.Param("tagID")

	var post Post
	var tag Tag

	if err := db.First(&post, postID).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Post not found")
		return
	}
	if err := db.First(&tag, tagID).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Tag not found")
		return
	}

	if err := db.Model(&post).Association("Tags").Append(&tag); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to associate tag with post")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag added to post"})
}

// @title Go API with Swagger
// @version 1.0
// @description This is a sample API with Swagger integration.
// @host localhost:8000
// @contact.name API Support
// @contact.url http://localhost:8000/support   // Local URL for your development environment
// @contact.email support@localhost.com
func main() {
	initDB()

	r := gin.Default()
	r.Use(cors.Default())

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	r.GET("/api/v1/users", getUsers)
	r.GET("/api/v1/users/:id", getUser)
	r.POST("/api/v1/users", createUser)

	// Post routes
	r.POST("/api/v1/users/:userID/posts", createPost)

	// Tag routes
	r.GET("/api/v1/tags", getTags)
	r.POST("/api/v1/tags", createTag)
	r.POST("/api/v1/posts/:postID/tags/:tagID", addTagToPost)
	// Start the server
	if err := r.Run(":8000"); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
