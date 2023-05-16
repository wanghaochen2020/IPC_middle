package api

import (
	"context"
	"encoding/json"
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
	Time int         `json:"time"`
	Data Atmosphere2 `json:"data"`
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

type Forecast struct {
	Time int
	Data []Forecast2
}

type Forecast2 struct {
	Time        string
	Temperature string
	Humidity    string
	Wind        string
}

// 获取气象实时数据的请求
// index:气象四个,roomrate,predict,7d
func GetData(c *gin.Context) {
	index := c.Query("index")
	base := c.Query("base")
	zutuan := c.Query("zutuan")

	if index == "predict" {
		a := getPredict()
		c.String(http.StatusOK, a)
	} else if index == "7d" {
		a := Get7d()
		c.JSON(http.StatusOK, gin.H{
			"data": a,
		})
	} else {
		base2, _ := strconv.Atoi(base)
		a := getData2(index, base2, zutuan)
		c.JSON(http.StatusOK, gin.H{
			"data": a,
		})
	}
}

func Get7d() []float64 {
	now := int(time.Now().Unix())
	var devices []Forecast
	data, _ := model.Db.Collection("weatherForecast").Find(context.TODO(), bson.M{"time": bson.M{"$gte": now - 3600, "$lt": now}}) //获取最新的天气预报
	err := data.All(context.TODO(), &devices)
	if err != nil {
		log.Println(err)
	}

	var result []float64
	result = make([]float64, 7)
	var sum float64
	for i := 0; i < 7; i++ {
		sum = 0
		for j := 0; j < 24; j++ {
			a, _ := strconv.ParseFloat(devices[0].Data[i*24+j].Temperature, 64)
			sum += a
		}
		result[i] = sum / 24
	}

	return result
}

func getPredict() string {
	now := int(time.Now().Unix())
	var devices []Forecast
	data, _ := model.Db.Collection("weatherForecast").Find(context.TODO(), bson.M{"time": bson.M{"$gte": now - 3600, "$lt": now}}) //获取最新的天气预报
	err := data.All(context.TODO(), &devices)
	if err != nil {
		log.Println(err)
	}
	a, _ := json.Marshal(devices[len(devices)-1])
	return string(a)
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
			for i := 0; i < len(devices); i++ {
				array[i], _ = strconv.ParseFloat(devices[i].Data.Temperature.Value, 64)
			}
		} else if index == "humidity" {
			for i := 0; i < len(devices); i++ {
				array[i], _ = strconv.ParseFloat(devices[i].Data.Humidity.Value, 64)
			}
		} else if index == "radiation" {
			for i := 0; i < len(devices); i++ {
				array[i], _ = strconv.ParseFloat(devices[i].Data.Radiation.Value, 64)
			}
		} else if index == "wind" {
			for i := 0; i < len(devices); i++ {
				array[i], _ = strconv.ParseFloat(devices[i].Data.Wind.Value, 64)
			}
		}

		if len(devices) < 24 {
			for i := len(devices); i < 24; i++ {
				array[i] = 0
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
