package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang/freetype"
	"github.com/google/uuid"
)

func main() {
	r := gin.Default()

	//必须先加载所有模板
	r.LoadHTMLGlob("views/*")

	r.GET("/quote", func(c *gin.Context) {
		// 获取名言和名人信息
		quote := c.Query("quote")
		author := c.Query("author")

		// 创建一个空白的图片
		img := image.NewRGBA(image.Rect(0, 0, 800, 600))
		draw.Draw(img, img.Bounds(), image.Black, image.ZP, draw.Src)

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
		ctx.SetSrc(image.White)
		ctx.SetDst(img)

		// 绘制名言和名人信息
		pt := freetype.Pt(100, 100)
		_, err = ctx.DrawString(quote, pt)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error: %v", err)
			return
		}
		pt.Y += ctx.PointToFixed(48 * 1.5)
		_, err = ctx.DrawString("- "+author, pt)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error: %v", err)
			return
		}

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
