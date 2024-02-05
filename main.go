package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

type CommonHdr struct {
	First bool
	Last  bool
	Count int
}

type Topo struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	TopoId string `json:"topoid"`
}

type Device struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	DeviceId string `json:"deviceid"`
}

type Topos struct {
	CommonHdr
	Data []Topo `json:"data"`
}

type Devices struct {
	CommonHdr
	Data []Device `json:"data"`
}

func getTopos(c *gin.Context) {
	log.Printf("Get topos")
	topos := Topos{
		CommonHdr: CommonHdr{First: true, Last: true, Count: 3},
	}
	for i := 0; i < 3; i++ {
		uu, _ := uuid.NewRandom()
		t := Topo{
			Id:     fmt.Sprintf("%d", i),
			Name:   fmt.Sprintf("Topo %d", i),
			TopoId: uu.String(),
		}
		topos.Data = append(topos.Data, t)
	}

	c.JSON(http.StatusOK, topos)
}
func getDevice(c *gin.Context) {
	dev := c.Param("id")
	log.Printf("dev is %s", dev)
	device := Device{
		Name:     fmt.Sprintf("A device %s", dev),
		DeviceId: fmt.Sprintf("A devID %s", dev),
	}
	time.Sleep(500 * time.Millisecond)
	c.JSON(http.StatusOK, device)
}

func getDevices(c *gin.Context) {
	log.Printf("Get devices")
	devices := Devices{
		CommonHdr: CommonHdr{First: true, Last: true, Count: 3},
	}

	for i := 0; i < 3; i++ {
		uu, _ := uuid.NewRandom()
		t := Device{
			Id:       fmt.Sprintf("%d", i),
			Name:     fmt.Sprintf("Device %d", i),
			DeviceId: uu.String(),
		}
		devices.Data = append(devices.Data, t)
	}

	c.JSON(http.StatusOK, devices)

}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.GET("/topos", getTopos)
	r.GET("/devices", getDevices)
	r.GET("/devices/:id", getDevice)
	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
