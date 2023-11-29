package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

var useMLModel bool = true // Configuration option to enable/disable ML model response

func main() {
	r := gin.Default()
	r.GET("/recommendations", func(c *gin.Context) {
		var recommendations []string
		if useMLModel {
			recommendations = getRecommendationsFromMLModel()
		} else {
			recommendations = getPredefinedRecommendations()
		}
		c.JSON(200, gin.H{
			"recommendations": recommendations,
		})
	})
	r.Run(":8080")
}

func getRecommendationsFromMLModel() []string {
	// Code to fetch personalized beer recommendations from the ML model serving REST API
	url := "https://credit-endpoint-fde6e229.eastus2.inference.ml.azure.com/score"
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + "hbvUnMRdYFtoRqPuoV1DFYRtsG0vFmTy"
	var boday = []byte(`{
		  "input_data": {
			"columns": [0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22],
			"index": [0, 1],
			"data": [
					[20000,2,2,1,24,2,2,-1,-1,-2,-2,3913,3102,689,0,0,0,0,689,0,0,0,0],
					[10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 10, 9, 8]
				]
		  },
		  "params": {}
	}`)
	// Create a new request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(boday))
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("azureml-model-deployment", "credit-defaults-model-1")

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	log.Println("resp from ML model:" + string([]byte(body)))
	recommendations := []string{"Beer X", "Beer Y", "Beer Z"}
	return recommendations
}

func getPredefinedRecommendations() []string {
	// Code to fetch pre-defined recommendations from the Inventory DB
	// Replace this with the actual implementation to fetch pre-defined recommendations
	recommendations := []string{"Beer A", "Beer B", "Beer C"}
	return recommendations
}
