package main

import (
	"log"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micmonay/keybd_event"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	// Select keys to be pressed
	kb.SetKeys(keybd_event.VK_MEDIA_PLAY_PAUSE)

	// Press the selected keys
	play := func() error {
		kb.SetKeys(keybd_event.VK_MEDIA_PLAY_PAUSE)
		err = kb.Launching()
		if err != nil {
			return err
		}
		return nil
	}
	space := func() error {
		kb.SetKeys(keybd_event.VK_SPACE)
		err = kb.Launching()
		if err != nil {
			return err
		}
		return nil
	}
	alttab := func() error {
		kb.HasALT(true)
		kb.SetKeys(keybd_event.VK_TAB)
		err = kb.Launching()
		kb.HasALT(false)
		if err != nil {
			return err
		}
		return nil
	}
	left := func() error {
		kb.SetKeys(keybd_event.VK_LEFT)
		err = kb.Launching()
		if err != nil {
			return err
		}
		return nil
	}
	right := func() error {
		kb.SetKeys(keybd_event.VK_RIGHT)
		err = kb.Launching()
		if err != nil {
			return err
		}
		return nil
	}
	up := func() error {
		kb.SetKeys(keybd_event.VK_UP)
		err = kb.Launching()
		if err != nil {
			return err
		}
		return nil
	}
	down := func() error {
		kb.SetKeys(keybd_event.VK_DOWN)
		err = kb.Launching()
		if err != nil {
			return err
		}
		return nil
	}

	qrc, err := qrcode.New("http://" + GetOutboundIP().String() + ":9145")
	if err != nil {
		log.Printf("could not generate QRCode: %v", err)
		return
	}

	w, err := standard.New("qr.png")
	if err != nil {
		log.Printf("standard.New failed: %v", err)
		return
	}

	// save file
	if err = qrc.Save(w); err != nil {
		log.Printf("could not save image: %v", err)
	}

	router := gin.Default()
	router.LoadHTMLFiles("index.html")
	api := router.Group("/")
	{
		api.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{})
		})
		api.GET("/plpa", func(c *gin.Context) {
			err := play()
			if err != nil {
				c.JSON(200, gin.H{"status": "error", "message": err})
				return
			}
			c.JSON(200, gin.H{"status": "success", "message": "play/pause"})
		})
		api.GET("/space", func(c *gin.Context) {
			err := space()
			if err != nil {
				c.JSON(200, gin.H{"status": "error", "message": err})
				return
			}
			c.JSON(200, gin.H{"status": "success", "message": "space"})
		})
		api.GET("/alttab", func(c *gin.Context) {
			err := alttab()
			if err != nil {
				c.JSON(200, gin.H{"status": "error", "message": err})
				return
			}
			c.JSON(200, gin.H{"status": "success", "message": "alttab"})
		})
		api.GET("/left", func(c *gin.Context) {
			err := left()
			if err != nil {
				c.JSON(200, gin.H{"status": "error", "message": err})
				return
			}
			c.JSON(200, gin.H{"status": "success", "message": "left"})
		})
		api.GET("/right", func(c *gin.Context) {
			err := right()
			if err != nil {
				c.JSON(200, gin.H{"status": "error", "message": err})
				return
			}
			c.JSON(200, gin.H{"status": "success", "message": "right"})
		})
		api.GET("/up", func(c *gin.Context) {
			err := up()
			if err != nil {
				c.JSON(200, gin.H{"status": "error", "message": err})
				return
			}
			c.JSON(200, gin.H{"status": "success", "message": "up"})
		})
		api.GET("/down", func(c *gin.Context) {
			err := down()
			if err != nil {
				c.JSON(200, gin.H{"status": "error", "message": err})
				return
			}
			c.JSON(200, gin.H{"status": "success", "message": "down"})
		})
	}
	router.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{}) })
	// Start listening and serving requests
	router.Run(":9145")

}
