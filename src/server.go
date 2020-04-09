package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"time"
)

type Employee struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required"`
	Age      int    `form:"age" binding:"required"`
	Position string `form:"position" binding:"required"`
	Gender   string `form:"gender" binding:"required"`
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/gocrud"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary()) // if connection fails throws error
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("connected to the database")
	}

	coll := client.Database("gocrud").Collection("employees")

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
			"title": "Dashboard",
		})
	})

	router.GET("/all", func(c *gin.Context) {
		var allemp []bson.M
		cursor, err := coll.Find(
			context.Background(),
			bson.D{},
		)
		if err != nil {
			log.Fatal(err)
		} else {
			if err = cursor.All(context.Background(), &allemp); err != nil {
				log.Fatal(allemp)
			}
			c.JSON(http.StatusOK, gin.H{"employees": allemp})
		}
	})

	router.POST("/add", func(c *gin.Context) {
		var emp Employee
		if err := c.ShouldBind(&emp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		addEmp, err := coll.InsertOne(context.Background(), bson.D{
			{"name", emp.Name},
			{"email", emp.Email},
			{"age", emp.Age},
			{"position", emp.Position},
			{"gender", emp.Gender},
		})
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(addEmp.InsertedID)
			c.String(http.StatusOK, "success")
		}

	})

	router.PUT("/update/:id", func(c *gin.Context) {
		var emp Employee
		if err := c.ShouldBind(&emp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id, _ := primitive.ObjectIDFromHex(c.Param("id"))
		updEmp, err := coll.UpdateOne(context.Background(), bson.M{
			"_id": id,
		},
			bson.D{
				{"$set", bson.D{
					{"name", emp.Name},
					{"email", emp.Email},
					{"age", emp.Age},
					{"position", emp.Position},
					{"gender", emp.Gender},
				}},
			},
		)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(updEmp.ModifiedCount)
		}
	})

	router.DELETE("/delete/:id", func(c *gin.Context) {
		id, _ := primitive.ObjectIDFromHex(c.Param("id"))
		fmt.Println(id)
		deleted, err := coll.DeleteOne(context.Background(), bson.D{
			{"_id", id},
		})
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(deleted.DeletedCount)
		}

	})

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}
