package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/upload", uploadFile)
	e.GET("/files", getFiles)
	e.DELETE("/delete/:name", deleteFile)
	e.Logger.Fatal(e.Start(":8000"))
}

type Document struct {
	Id          string `json:"id" `
	Name        string `json:"name"`
	Size        int
	RawContent  string `json:"raw_content"`
	WorkingHtml string `json:"workingHtml"`
}

func getFiles(c echo.Context) error {
	files, err := FilePathWalkDir("/filesBurea/")
	if err != nil {
		panic(err)
	}
	var allFiles []Document
	for _, fileA := range files {
		content, err := ioutil.ReadFile("/filesBurea/" + fileA.Name())
		if err != nil {
			fmt.Println(err)
		}
		document := Document{
			Id:         strings.Split(fileA.Name(), ".")[0],
			Name:       strings.Split(fileA.Name(), ".")[0],
			Size:       int(fileA.Size()),
			RawContent: string(content),
		}
		allFiles = append(allFiles, document)
		fmt.Println(fileA)
	}
	return c.JSON(http.StatusOK, allFiles)
}

func deleteFile(c echo.Context) error {
	fileName := c.Param("name")

	err := os.Remove("/filesBurea/" + fileName + ".txt")
	if err != nil {
		fmt.Println(err)
	}

	return c.JSON(http.StatusOK, "Archivo eliminado con éxito")
}

func uploadFile(c echo.Context) error {
	var document Document

	err := json.NewDecoder(c.Request().Body).Decode(&document)
	saveDocument(document.Id, document.RawContent, document.Name)
	if err != nil {
		fmt.Println(err)
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>Document saved correctly</p>"))
}

func saveDocument(id string, rawContent string, name string) {
	errAtC := os.Mkdir("/filesBurea/", 0700)
	if errAtC != nil {
		fmt.Println(errAtC)
	}
	if _, err := os.Stat("/filesBurea/"); err != nil {
		fmt.Println(err)
		os.Mkdir("/filesBurea/", 0700)
	}
	if _, err := os.Stat("/filesBurea/" + name + ".txt"); err == nil {
		fmt.Println(err)
	}

	f, err := os.Create("/filesBurea/" + name + ".txt")
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(rawContent)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("done")
}

func FilePathWalkDir(root string) ([]os.FileInfo, error) {
	var files []os.FileInfo
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, info)
		}
		return nil
	})
	return files, err
}
