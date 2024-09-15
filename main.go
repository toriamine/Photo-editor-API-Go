package main

import (
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"image"
	"image/color"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	r := gin.Default()

	// Загрузка шаблонов
	r.LoadHTMLGlob("/app/html/template/*")

	// Маршруты
	r.POST("/upload", uploadImage)
	r.POST("/filter/:filter", applyFilter)
	r.GET("/download/:filename", downloadImage)
	r.GET("/images", listImages)
	r.DELETE("/images/:filename", deleteImage)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Image Processing API",
		})
	})

	r.Run(":8081")
}

func uploadImage(c *gin.Context) {
	// Обработка файла изображения
	file, _ := c.FormFile("image")
	filename := file.Filename
	c.SaveUploadedFile(file, filepath.Join("uploads", filename))
	c.JSON(http.StatusOK, gin.H{
		"message":  "Изображение загружено успешно",
		"filename": filename,
	})
}

func applyFilter(c *gin.Context) {
	filter := c.Param("filter")
	filename := c.Query("filename")
	img, err := imaging.Open(filepath.Join("uploads", filename))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось открыть изображение: " + err.Error()})
		return
	}

	var filteredImage image.Image
	switch filter {
	case "grayscale":
		filteredImage = imaging.Grayscale(img)
	case "sepia":
		// Применяем эффект сепии с помощью стандартной библиотеки image
		filteredImage = image.NewRGBA(img.Bounds())
		for x := 0; x < img.Bounds().Max.X; x++ {
			for y := 0; y < img.Bounds().Max.Y; y++ {
				c := img.At(x, y)
				r, g, b, _ := c.RGBA()
				newR := uint8((float64(r) * 0.393) + (float64(g) * 0.769) + (float64(b) * 0.189))
				newG := uint8((float64(r) * 0.349) + (float64(g) * 0.686) + (float64(b) * 0.168))
				newB := uint8((float64(r) * 0.272) + (float64(g) * 0.534) + (float64(b) * 0.131))
				filteredImage.(*image.RGBA).Set(x, y, color.RGBA{newR, newG, newB, uint8(r)})
			}
		}
	case "blur":
		filteredImage = imaging.Blur(img, 5)
	case "sharpen":
		filteredImage = imaging.Sharpen(img, 2)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимый фильтр"})
		return
	}

	newFilename := filepath.Join("uploads", "filtered_"+filename)
	err = imaging.Save(filteredImage, newFilename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить отфильтрованное изображение: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Изображение отфильтровано успешно",
		"filename": "filtered_" + filename,
	})
}

func downloadImage(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join("uploads", filename)
	c.FileAttachment(filePath, filename)
}

func listImages(c *gin.Context) {
	files, err := os.ReadDir("uploads")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список изображений: " + err.Error()})
		return
	}

	var imageFiles []string
	for _, file := range files {
		if !file.IsDir() {
			imageFiles = append(imageFiles, file.Name())
		}
	}

	c.JSON(http.StatusOK, imageFiles)
}

func deleteImage(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join("uploads", filename)
	err := os.Remove(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить изображение: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Изображение удалено успешно"})
}
