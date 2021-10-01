package handlers

import (
	"context"
	"os"

	"github.com/PoppyPenguin-Metadata/app/config"
	"github.com/PoppyPenguin-Metadata/constants"
	"github.com/PoppyPenguin-Metadata/db"
	"github.com/PoppyPenguin-Metadata/structs"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetRarityById endpoints accepts id of a single Poppy and information for a single Poppy.
//
// If no Poppy Penguin is found returns empty response
func GetRarityById(id int) structs.RarityServiceResponse {
	godotenv.Load()
	poppyPenguinDBName := os.Getenv("POPPYPENGUIN_DB")
	rarityCollectionName := os.Getenv("RARITY_COLLECTION")
	collection, err := db.GetMongoDbCollection(poppyPenguinDBName, rarityCollectionName)
	if err != nil {
		return structs.RarityServiceResponse{}
	}

	findOptions := options.FindOneOptions{}
	removePrivateFieldsSingle(&findOptions)

	var filter bson.M = bson.M{}
	filter = bson.M{constants.MorphFieldNames.TokenId: id}

	var result = structs.RarityServiceResponse{}
	curr := collection.FindOne(context.Background(), filter, &findOptions)

	curr.Decode(&result)

	return result
}

// removePrivateFieldsSingle removes internal fields that are of no interest to the users of the API.
//
// Configuration of these fields can be found in helpers.apiConfig.go
func removePrivateFieldsSingle(findOptions *options.FindOneOptions) {
	noProjectionFields := bson.M{}
	for _, field := range config.MORPHS_NO_PROJECTION_FIELDS {
		noProjectionFields[field] = 0
	}
	findOptions.SetProjection(noProjectionFields)
}
