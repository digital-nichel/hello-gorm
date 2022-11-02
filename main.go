package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

type author struct {
	ID   uint
	Name string
}

type book struct {
	ID        uint
	AuthorID  uint
	Author    author
	Title     string
	Libraries []library `gorm:"many2many:catalogs;"`
}

type library struct {
	ID    uint
	Name  string
	Books []book `gorm:"many2many:catalogs"`
}

func main() {
	e := echo.New()

	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=testdb sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		e.Logger.Fatal(err)
	}

	// get all authors
	e.GET("/authors", func(c echo.Context) error {
		var authors []author

		err := db.Find(&authors).Error
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, authors)
	})

	// get all books
	e.GET("/books/extended", func(c echo.Context) error {
		var books []book

		err := db.Preload(clause.Associations).Find(&books).Error
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, books)
	})

	// get all books
	e.GET("/books", func(c echo.Context) error {
		var books []book

		err := db.Find(&books).Error
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, books)
	})

	// get books by authorID
	e.GET("/books/:authorID", func(c echo.Context) error {
		var books []book

		err = db.Joins("Author").Where(&author{
			ID: 2,
		}).Preload(clause.Associations).Find(&books).Error
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, books)
	})

	// get library by ID
	e.GET("/library/:id", func(c echo.Context) error {
		libraryId := c.Param("id")

		var library library

		err := db.Preload(clause.Associations).First(&library, libraryId).Error
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, library)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
