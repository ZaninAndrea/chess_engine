package main

import (
	"time"

	"github.com/ZaninAndrea/chess_engine/chessboard"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"ok": true,
		})
	})

	r.GET("/bestmove", func(c *gin.Context) {
		fen := c.DefaultQuery("fen", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

		game := chessboard.NewGameFromFEN(fen)
		engine := chessboard.NewBruteForceEngine(&game)

		game.Move(engine.BestMove(60))
		pos := game.Position()
		c.JSON(200, gin.H{
			"fen":    pos.FEN(),
			"result": game.Result().String(),
		})
	})

	r.Run()
}
