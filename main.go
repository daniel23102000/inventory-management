package main

import (
    "github.com/gin-gonic/gin"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "net/http"
)

// Struct untuk Produk
type Product struct {
    ID          uint    `json:"id" gorm:"primaryKey"`
    Nama        string  `json:"nama"`
    Deskripsi   string  `json:"deskripsi"`
    Harga       float64 `json:"harga"`
    Kategori    string  `json:"kategori"`
}

var db *gorm.DB

func main() {
    // Koneksi ke database
    dsn := "root:password@tcp(127.0.0.1:3306)/inventory_db?charset=utf8mb4&parseTime=True&loc=Local"
    var err error
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database")
    }

    // Migrasi database
    db.AutoMigrate(&Product{})

    // Setup router Gin
    r := gin.Default()

    // Routes
    r.GET("/products", GetProducts)       // Lihat semua produk
    r.POST("/products", AddProduct)      // Tambah produk
    r.POST("/upload/:id", UploadImage)   // Unggah gambar
    r.GET("/download/:id", DownloadImage)// Unduh gambar

    // Jalankan server
    r.Run(":8080")
}

// Handler untuk produk
func GetProducts(c *gin.Context) {
    var products []Product
    db.Find(&products)
    c.JSON(http.StatusOK, products)
}

func AddProduct(c *gin.Context) {
    var product Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    db.Create(&product)
    c.JSON(http.StatusOK, gin.H{"message": "Product added"})
}

// Handler untuk file
func UploadImage(c *gin.Context) {
    file, _ := c.FormFile("file")
    path := "uploads/" + file.Filename
    if err := c.SaveUploadedFile(file, path); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "File uploaded", "path": path})
}

func DownloadImage(c *gin.Context) {
    id := c.Param("id")
    path := "uploads/" + id + ".jpg"
    c.File(path)
}