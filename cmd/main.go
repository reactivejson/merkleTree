package main

import (
	"github.com/gin-gonic/gin"
	"github.com/reactivejson/merkleTree/api"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */
// The core entry point into the app. will setup the config, and run the App
func main() {
	router := gin.Default()
	router.Use(api.ErrorHandler)
	router.POST("/create", api.CreateTree)
	router.PUT("/update", api.UpdateLeaf)
	router.POST("/verify", api.VerifyProof)
	router.POST("/visual/proof", api.VisualizeProof)
	router.Run(":8080")
}
