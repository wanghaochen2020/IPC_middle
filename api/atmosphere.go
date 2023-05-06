package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"middle/model"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Atmosphere struct {
	Time   int         `json:"time"`
	Result Atmosphere2 `json:"data"`
}

type Atmosphere2 struct {
	Wind        Atmosphere3
	Temperature Atmosphere3
	Humidity    Atmosphere3
	Radiation   Atmosphere3
}

type Atmosphere3 struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Kekong struct {
	Time int     `json:"time"`
	D4   float64 `json:"d4"`
	D5   float64 `json:"d5"`
	D6   float64 `json:"d6"`
}

// 获取气象实时数据的请求
func GetData(c *gin.Context) {
	index := c.Query("index")
	base := c.Query("base")
	zutuan := c.Query("zutuan")

	base2, _ := strconv.Atoi(base)

	a := getData2(index, base2, zutuan)

	c.JSON(http.StatusOK, gin.H{
		"data": a,
	})
}

func getData2(index string, base int, zutuan string) []float64 {
	var array []float64
	array = make([]float64, 24)
	limit0, limit1 := FindIntervalDay(base)

	type Update struct {
		Id        int
		UpdatedAt string
		TableName string
	}

	switch index {
	case "roomRate":
		var devices []Kekong
		data, _ := model.Db.Collection("keKong").Find(context.TODO(), bson.M{"time": bson.M{"$gte": limit0, "$lt": limit1}})
		err := data.All(context.TODO(), &devices)
		if err != nil {
			log.Println(err)
		}
		if zutuan == "D4组团" {
			for i := 0; i < 24; i++ {
				array[i] = devices[i].D4
			}
		} else if zutuan == "D5组团" {
			for i := 0; i < 24; i++ {
				array[i] = devices[i].D5
			}
		} else if zutuan == "D6组团" {
			for i := 0; i < 24; i++ {
				array[i] = devices[i].D6
			}
		} else {
			for i := 0; i < 24; i++ {
				array[i] = devices[i].D4
			}
		}
	default:
		var devices []Atmosphere
		data, _ := model.Db.Collection("atmosphere").Find(context.TODO(), bson.M{"time": bson.M{"$gte": limit0, "$lt": limit1}})
		err := data.All(context.TODO(), &devices)
		if err != nil {
			log.Println(err)
		}
		if index == "temperature" {
			for i := 0; i < 24; i++ {
				array[i], _ = strconv.ParseFloat(devices[i].Result.Temperature.Value, 64)
			}
		} else if index == "humidity" {
			for i := 0; i < 24; i++ {
				array[i], _ = strconv.ParseFloat(devices[i].Result.Humidity.Value, 64)
			}
		} else if index == "radiation" {
			for i := 0; i < 24; i++ {
				array[i], _ = strconv.ParseFloat(devices[i].Result.Radiation.Value, 64)
			}
		} else if index == "wind" {
			for i := 0; i < 24; i++ {
				array[i], _ = strconv.ParseFloat(devices[i].Result.Wind.Value, 64)
			}
		}
	}

	return array
}

func FindIntervalDay(value int) (int, int) {
	Time := time.Unix(int64(value), 0)
	Time1 := time.Date(Time.Year(), Time.Month(), Time.Day(), 0, 0, 0, 0, Time.Location())
	Time2 := time.Date(Time.Year(), Time.Month(), Time.Day(), 24, 0, 0, 0, Time.Location())
	return int(Time1.Unix()), int(Time2.Unix())
}

// 将字符串中的空格转为%20
func trans(input string) string {
	res := strings.Replace(input, " ", "%20", -1)
	return res
}
