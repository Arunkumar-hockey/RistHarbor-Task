package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"test/database"
	"test/model"

	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var blocksCollection *mongo.Collection = database.OpenCollection(database.Client, "blocks")


func SaveBlocks(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)


		// Establish a connection to an Ethereum client
		url := "https://mainnet.infura.io/v3/a08d8375f6ab4bd597f3c93387de6db1"
		client, err := ethclient.Dial(url)
		if err != nil {
			log.Fatal(err)
		}
	
		// Get the latest block number
		header, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}

		defer cancel()
	
		// Retrieve the latest block
		block, err := client.BlockByNumber(context.Background(), header.Number)
		if err != nil {
			log.Fatal(err)
		}
	
		// Iterate over the transactions in the block
		for _, tx := range block.Transactions() {
			fmt.Println("Transaction Hash:", tx.Hash().Hex())
			//fmt.Println("From:", tx.From().Hex())
			fmt.Println("To:", tx.To().Hex())
			fmt.Println("Value:", tx.Value().String())
			fmt.Println()
			

			a := []interface{}{
				&model.Blocks{
					TransactionHash:tx.Hash().Hex(),
				},
			}
			
			
			

			//blocks = append(blocks, a)

		

			_, err := blocksCollection.InsertMany(ctx, a)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("====Inserted Successfully")
			c.JSON(http.StatusOK, a)
			
		}

		c.JSON(http.StatusOK, "Inserted Successfully")
		
	}



func QueryBlocks(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)


	findOptions := options.Find()
    // Sort by recent
    findOptions.SetSort(bson.D{{"$natural", 1}})

	filter := bson.M{}
	sort :=options.Find().SetSort(map[string]int{"$natural":1}).SetLimit(5)


	result, err := blocksCollection.Find(context.Background(),filter,sort)

	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing table items"})
	}
	var recentBlocks []bson.M
	if err = result.All(ctx, &recentBlocks); err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, recentBlocks)
}

// filter := bson.M{"area_code": mode.AreaCode}

// var allGames []bson.M

// result, err := gameModeCollection.Find(context.TODO(), filter, options.Find().SetSort(map[string]int{"online_players": -1}))