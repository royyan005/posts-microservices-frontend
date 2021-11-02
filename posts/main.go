package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Post struct {
    Id              uint        `json:"id"`
    Title           string      `json:"title"`
    Description     string      `json:"description"`
    Comments        []Comment   `json:"comments" gorm:"-" default:"[]"`
}

type Comment struct {
    Id      uint    `json:"id"`
    PostId  uint    `json:"post_id"`
    Text    string  `json:"text"`
}

func main() {
    db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/posts_ms"), &gorm.Config{})

    if err != nil {
        panic(err)
    }

    db.AutoMigrate(Post{})
    
    app := fiber.New()

    app.Use(cors.New())

    app.Get("/api/posts", func(c *fiber.Ctx) error {
        var posts []Post

        db.Find(&posts)

        for i,post := range posts {
            response, err := http.Get(fmt.Sprintf("http://localhost:8001/api/posts/%d/comments", post.Id))

            if err != nil {
                return err
            }

            var comments []Comment

            json.NewDecoder(response.Body).Decode(&comments)

            posts[i].Comments = comments
        }

        return c.JSON(posts)
    })

    app.Post("/api/posts", func(c *fiber.Ctx) error {
        var post Post

        if err := c.BodyParser(&post); err != nil {
            return err
        }

        db.Create(&post)

        return c.JSON(post)
    })

    app.Listen(":8000")
}