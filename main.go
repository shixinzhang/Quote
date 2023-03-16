package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	gorm "gorm.io/gorm"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"net/http"
	"os"
)

type Poster struct {
	ID     uint   `gorm:"primaryKey"`
	Image  []byte `gorm:"not null"`
	Text   string `gorm:"not null"`
	Author string `gorm:"not null"`
}

var db *gorm.DB

func main() {
	learnGin()

	//quote()
}

func learnGin() {
	//r := gin.Default()
	//
	//db, err := gorm.Open(mysql.Open("todo.db"), &gorm.Config{})
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func quote() {

	dsn := "user:password@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Poster{})

	r := gin.Default()

	//必须先加载所有模板
	r.LoadHTMLGlob("views/*")

	r.POST("/posters", generatePoster)

	r.GET("/quote", func(c *gin.Context) {
		// 获取名言和名人信息
		quote := c.Query("quote")
		author := c.Query("author")

		// 创建一个空白的图片
		img := image.NewRGBA(image.Rect(0, 0, 800, 600))
		draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{R: 255, G: 255, B: 255, A: 255}), image.ZP, draw.Src)

		// 打开字体文件
		fontFile, err := os.Open("font.ttf")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error: %v", err)
			return
		}
		defer fontFile.Close()

		// 读取字体文件并创建字体
		fontBytes, err := io.ReadAll(fontFile)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error: %v", err)
			return
		}
		font, err := freetype.ParseFont(fontBytes)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error: %v", err)
			return
		}

		// 创建绘图上下文
		ctx := freetype.NewContext()
		ctx.SetDPI(72)
		ctx.SetFont(font)
		ctx.SetFontSize(48)
		ctx.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 255, A: 255})) // 设置字体颜色为黄色
		ctx.SetDst(img)

		// 绘制名言和名人信息
		pt := freetype.Pt(100, 100)
		_, err = ctx.DrawString(quote, pt)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error: %v", err)
			return
		}
		println("DrawString: " + quote)

		pt.Y += ctx.PointToFixed(48 * 1.5)
		_, err = ctx.DrawString("- "+author, pt)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error: %v", err)
			return
		}
		println("DrawString: " + author)

		// 保存图片并返回URL
		fileName := fmt.Sprintf("%s.jpg", uuid.New().String())
		filePath := fmt.Sprintf("posters/%s", fileName)
		file, err := os.Create(filePath)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error: %v", err)
			return
		}
		defer file.Close()
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 80})
		if err != nil {
			c.String(http.StatusInternalServerError, "Error: %v", err)
			return
		}
		//url := fmt.Sprintf("https://example.com/%s", fileName)
		url := filePath

		println("quote: " + quote + ", author:" + author + ", url: " + url)

		//c.JSON(http.StatusOK, gin.H{
		//	"quote":  quote,
		//	"author": author,
		//	"url":    url,
		//})

		//模板名称
		c.HTML(http.StatusOK, "poster.html", gin.H{
			"quote":  quote,
			"author": author,
			"url":    url,
		})
	})
	r.Run(":8089")
}

func generatePoster(c *gin.Context) {
	//imageData, err := ioutil.ReadAll(c.Request.Body)
	//if err != nil {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "读取图片数据时发生错误：" + err.Error()})
	//	return
	//}
	//
	//image, err := imaging.Decode(bytes.NewReader(imageData))
	//if err != nil {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "解码图片时发生错误：" + err.Error()})
	//	return
	//}
	//
	//text := c.PostForm("text")
	//author := c.PostForm("author")
	//if text == "" {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "缺少文字参数"})
	//	return
	//}
}
