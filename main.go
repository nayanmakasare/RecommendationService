package main

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"github.com/nayanmakasare/RecommendationService/neo4j"
	RecommendationService "github.com/nayanmakasare/RecommendationService/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	defaultHost = "mongodb://nayan:tlwn722n@cluster0-shard-00-00-8aov2.mongodb.net:27017,cluster0-shard-00-01-8aov2.mongodb.net:27017,cluster0-shard-00-02-8aov2.mongodb.net:27017/test?ssl=true&replicaSet=Cluster0-shard-0&authSource=admin&retryWrites=true&w=majority"
	//defaultHost = "mongodb://192.168.1.9:27017"
	//defaultHost = "mongodb://192.168.1.143:27017"
)

func main(){

	err := neo4j.InitNeo4j()
	if err != nil {
		log.Fatal(err)
	}

	service := grpc.NewService(
		micro.Name("RecommendationService"),
		micro.Address(":50054"),
		micro.Version("1.0"),
	)
	service.Init()
	uri := defaultHost
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Debug(err)
	}
	client := GetRedisClient()
	// Register Handler
	redrecObj,err := New("localhost:6379")
	if err != nil {
		log.Info("nayan",err)
	}
	h := RecommendationServiceHandler{CloudwalkerRecommendation:redrecObj,
		MongoCollection:mongoClient.Database("test").Collection("cwmovies"),
		RedisConnection:client}

	err = RecommendationService.RegisterRecommendationServiceHandler(service.Server(), &h)
	if err != nil {
		log.Fatal(err)
	}


	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func GetRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Could not connect to redis %v", err)
	}
	return client
}