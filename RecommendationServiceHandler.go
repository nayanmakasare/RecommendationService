package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/micro/go-micro/util/log"
	RecommendationService "github.com/nayanmakasare/RecommendationService/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type RecommendationServiceHandler struct {
	RedisConnection  *redis.Client
	CloudwalkerRecommendation *CloudwalkerRedrec
	MongoCollection *mongo.Collection
}


var counter = 0

func (h *RecommendationServiceHandler) TileClicked(ctx context.Context, req *RecommendationService.TileClickedRequest, res *RecommendationService.TileClickedResponse) error {
	go h.CloudwalkerRecommendation.Rate(req.TileId, req.UserId, req.TileScore)
	return nil
}

func (h *RecommendationServiceHandler) GetCollabrativeFilteringData(ctx context.Context, req *RecommendationService.GetRecommendationRequest, stream RecommendationService.RecommendationService_GetCollabrativeFilteringDataStream) error{
	log.Info("Colab Trigger ")
	//err := h.CloudwalkerRecommendation.BatchUpdateSimilarUsers(2)
	//if err != nil {
	//	return err
	//}
	err := h.CloudwalkerRecommendation.UpdateSuggestedItems(req.UserId, 50)
	if err != nil {
		return err
	}
	userRedisKey := fmt.Sprintf("user:%s:collabrativeFilter", req.UserId)
	result, err :=  h.CheckInRedis(userRedisKey)
	if err != nil {
		return err
	}
	log.Info("Is Colab There ?", result)
	if result {
		stringResultCmd := h.RedisConnection.SMembers(userRedisKey)
		resultStringArray := stringResultCmd.Val()
		for i := 0 ; i < len(resultStringArray) ; i++ {
			var resultMovie RecommendationService.MovieTile
			err := proto.Unmarshal([]byte(resultStringArray[i]), &resultMovie)
			if err != nil {
				return err
			}
			err = stream.Send(&resultMovie)
			if err != nil {
				return err
			}
		}
		return stream.Close()
	} else {
		resultSet, err := h.CloudwalkerRecommendation.GetUserSuggestions(req.UserId, 5)
		if err != nil {
			return err
		}
		if len(resultSet) > 0 {
			for _,v := range resultSet {
				result := h.RedisConnection.HGet("cloudwalkerTiles", v)
				if result.Err() != nil {
					return result.Err()
				}
				resultInBytes, _ := result.Bytes()
				h.RedisConnection.Del(userRedisKey)
				h.RedisConnection.SAdd(userRedisKey,resultInBytes)
				var movieTiles RecommendationService.MovieTile
				err := proto.Unmarshal(resultInBytes, &movieTiles)
				if err != nil {
					return err
				}
				err = stream.Send(&movieTiles)
				if err != nil {
					return err
				}
			}
			return stream.Close()
		}else {
			return errors.New("No Recomendation")
		}
	}
}

func (h *RecommendationServiceHandler) InitialRecommendationEngine(ctx context.Context, req *RecommendationService.InitRecommendationRequest, res *RecommendationService.InitRecommendationResponse) error {
	log.Info("Triggered ")
	for _,v := range req.Genres {
		genreKey :=  fmt.Sprintf("genre:%s:items", v)
		log.Info(genreKey)
		myStages := mongo.Pipeline{
			//stage 1
			bson.D{{"$match", bson.D{{"metadata.genre", bson.D{{"$in", bson.A{v}}}}}}},
			//stage 2
			bson.D{{"$sort", bson.D{{"created_at", -1}}}},
			//stage 3
			bson.D{{"$project", bson.D{{"_id", 0},{"ref_id", 1}}}},
		}
		cur, err := h.MongoCollection.Aggregate(context.Background(), myStages, options.Aggregate().SetMaxTime(2000*time.Millisecond))
		if err != nil {
			return err
		}
		for cur.Next(ctx) {
			log.Info("Genre",cur.Current.Lookup("ref_id").StringValue())
			h.RedisConnection.SAdd(genreKey, cur.Current.Lookup("ref_id").StringValue())
		}
		err = cur.Close(ctx)
		if err != nil {
			return err
		}
	}

	// language Tiles Segregation
	for _,v := range req.Languages {
		languageKey :=  fmt.Sprintf("languages:%s:items", v)
		log.Info(languageKey)
		myStages := mongo.Pipeline{
			//stage 1
			bson.D{{"$match", bson.D{{"metadata.languages", bson.D{{"$in", bson.A{v}}}}}}},
			//stage 2
			bson.D{{"$sort", bson.D{{"created_at", -1}}}},
			//stage 3
			bson.D{{"$project", bson.D{{"_id", 0},{"ref_id", 1}}}},
		}
		cur, err := h.MongoCollection.Aggregate(context.Background(), myStages, options.Aggregate().SetMaxTime(2000*time.Millisecond))
		if err != nil {
			return err
		}
		for cur.Next(ctx) {
			log.Info("Languages",cur.Current.Lookup("ref_id").StringValue())
			h.RedisConnection.SAdd(languageKey, cur.Current.Lookup("ref_id").StringValue())
		}
		err = cur.Close(ctx)
		if err != nil {
			return err
		}
	}

	for _,v := range req.Categories {
		categoriesKey :=  fmt.Sprintf("categories:%s:items", v)
		log.Info(categoriesKey)
		myStages := mongo.Pipeline{
			//stage 1
			bson.D{{"$match", bson.D{{"metadata.categories", bson.D{{"$in", bson.A{v}}}}}}},
			//stage 2
			bson.D{{"$sort", bson.D{{"created_at", -1}}}},
			//stage 3
			bson.D{{"$project", bson.D{{"_id", 0},{"ref_id", 1}}}},
		}
		cur, err := h.MongoCollection.Aggregate(context.Background(), myStages, options.Aggregate().SetMaxTime(2000*time.Millisecond))
		if err != nil {
			return err
		}
		for cur.Next(ctx) {
			log.Info("Categories",cur.Current.Lookup("ref_id").StringValue())
			h.RedisConnection.SAdd(categoriesKey, cur.Current.Lookup("ref_id").StringValue())
		}
		err = cur.Close(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *RecommendationServiceHandler) CheckInRedis(redisKey string) (bool, error) {
	intCmdResult := h.RedisConnection.Exists(redisKey)
	if intCmdResult.Val() == 1 {
		return  true , nil
	}else {
		return false , nil
	}
}

func (h *RecommendationServiceHandler) GetContentbasedData(ctx context.Context, req *RecommendationService.GetRecommendationRequest, stream RecommendationService.RecommendationService_GetContentbasedDataStream) error {
	isThere, err := h.CheckInRedis(fmt.Sprintf("user:%s:Contentbased", req.UserId))
	if err != nil {
		return err
	}
	log.Info("Is ContentBased There ? ", isThere)
	if isThere {
		// get the updated Tile Clicked By the User
		h.DiffContentBaseToTileClickedByUser(req.UserId)
		result := h.RedisConnection.SMembers(fmt.Sprintf("user:%s:Contentbased", req.UserId))
		if result.Err() != nil {
			return result.Err()
		}
		tilesResult := h.RedisConnection.HMGet("cloudwalkerTiles", result.Val()...)
		for i := 0 ; i < len(tilesResult.Val()) ; i++ {
			if i < 20 {
				var resultMovie RecommendationService.MovieTile
				err := proto.Unmarshal([]byte(fmt.Sprintf("%v", tilesResult.Val()[i])), &resultMovie)
				if err != nil {
					return err
				}
				err = stream.Send(&resultMovie)
				if err != nil {
					return err
				}
			}else {
				break
			}
		}
		return stream.Close()
	}else {
		err = h.PreProcessingContentbasedData(req)
		if err != nil {
			return err
		}
		result := h.RedisConnection.SMembers(fmt.Sprintf("user:%s:Contentbased", req.UserId))
		if result.Err() != nil {
			return result.Err()
		}
		tilesResult := h.RedisConnection.HMGet("cloudwalkerTiles", result.Val()...)
		for i := 0 ; i < len(tilesResult.Val()) ; i++ {
			if i < 20 {
				var resultMovie RecommendationService.MovieTile
				err := proto.Unmarshal([]byte(fmt.Sprintf("%v", tilesResult.Val()[i])), &resultMovie)
				if err != nil {
					return err
				}
				err = stream.Send(&resultMovie)
				if err != nil {
					return err
				}
			}else {
				break
			}
		}
		return stream.Close()
	}
}

func (h *RecommendationServiceHandler) PreProcessingContentbasedData(req *RecommendationService.GetRecommendationRequest) error  {
	// getting genre for user
	result := h.RedisConnection.SMembers(fmt.Sprintf("user:%s:genre", req.UserId))
	if result.Err() != nil {
		return result.Err()
	}
	var tempGenre []string
	for _,v := range result.Val() {
		tempGenre = append(tempGenre,fmt.Sprintf("genre:%s:items", v))
	}
	genreResult := h.RedisConnection.SUnionStore(fmt.Sprintf("user:%s:genre:items", req.UserId),tempGenre...)
	log.Info("genre",genreResult.Val())
	// getting Languages for user
	result = h.RedisConnection.SMembers(fmt.Sprintf("user:%s:languages", req.UserId))
	if result.Err() != nil {
		return result.Err()
	}
	var tempLanguages []string
	for _,v := range result.Val() {
		tempLanguages = append(tempLanguages,fmt.Sprintf("languages:%s:items", v) )
	}
	log.Info(len(tempLanguages))
	h.RedisConnection.SUnionStore(fmt.Sprintf("user:%s:languages:items", req.UserId),tempLanguages...)

	// getting categories for user
	result = h.RedisConnection.SMembers(fmt.Sprintf("user:%s:categories", req.UserId))
	if result.Err() != nil {
		return result.Err()
	}
	var tempCategories []string
	for _,v := range result.Val() {
		tempCategories = append(tempCategories,fmt.Sprintf("categories:%s:items", v) )
	}
	log.Info(len(tempCategories))
	h.RedisConnection.SUnionStore(fmt.Sprintf("user:%s:categories:items", req.UserId),tempCategories...)
	// get tiles which content all its taste of user
	h.RedisConnection.SInterStore(fmt.Sprintf("user:%s:InterStore", req.UserId),
		fmt.Sprintf("user:%s:genre:items", req.UserId),
		fmt.Sprintf("user:%s:languages:items", req.UserId),
		fmt.Sprintf("user:%s:categories:items", req.UserId))

	tileClickedByUser := h.RedisConnection.ZRevRange(fmt.Sprintf("user:%s:items", req.UserId), 0, -1)
	h.RedisConnection.SAdd(fmt.Sprintf("user:%s:temp", req.UserId),tileClickedByUser.Val())
	resultOfDiff := h.RedisConnection.SDiffStore(fmt.Sprintf("user:%s:Contentbased", req.UserId),fmt.Sprintf("user:%s:InterStore", req.UserId),fmt.Sprintf("user:%s:temp", req.UserId))
	if resultOfDiff.Err() != nil {
		log.Fatal(resultOfDiff.Err())
	}
	log.Info(resultOfDiff.Val())
	h.RedisConnection.Del(fmt.Sprintf("user:%s:temp", req.UserId), fmt.Sprintf("user:%s:InterStore", req.UserId))
	return nil
}

func(h *RecommendationServiceHandler) DiffContentBaseToTileClickedByUser(userId string){
	tileClickedByUser := h.RedisConnection.ZRevRange(fmt.Sprintf("user:%s:items", userId), 0, -1)
	h.RedisConnection.SAdd(fmt.Sprintf("user:%s:temp", userId),tileClickedByUser.Val())
	log.Info("Before ", h.RedisConnection.SCard(fmt.Sprintf("user:%s:Contentbased", userId)))
	h.RedisConnection.SDiffStore(fmt.Sprintf("user:%s:Contentbased", userId), fmt.Sprintf("user:%s:Contentbased", userId),fmt.Sprintf("user:%s:temp", userId))
	log.Info("After ", h.RedisConnection.SCard(fmt.Sprintf("user:%s:Contentbased", userId)))
	h.RedisConnection.Del(fmt.Sprintf("user:%s:temp", userId))
}



