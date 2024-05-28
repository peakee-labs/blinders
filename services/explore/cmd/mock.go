package main

import (
	"context"
	"math/rand"

	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"blinders/packages/auth"
	"blinders/packages/db/matchingdb"
	"blinders/packages/db/usersdb"
	dbutils "blinders/packages/db/utils"
	"blinders/packages/explore"
	"blinders/packages/transport"
	"blinders/packages/utils"
	exploreapi "blinders/services/explore/api"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	err     error
	manager *exploreapi.Manager
)

func init() {
	envFile := ".env.dev"
	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("failed to load env", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	redisClient := utils.NewRedisClientFromEnv(ctx)

	var usersDB *mongo.Database
	var matchingDB *mongo.Database
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		usersDB, err = dbutils.InitMongoDatabaseFromEnv("USERS")
		if err != nil {
			log.Fatal("failed to init users db:", err)
		}
	}()
	go func() {
		defer wg.Done()
		matchingDB, err = dbutils.InitMongoDatabaseFromEnv("MATCHING")
		if err != nil {
			log.Fatal("failed to init matching db:", err)
		}
	}()
	wg.Wait()

	matchingRepo := matchingdb.NewMatchingRepo(matchingDB)
	usersRepo := usersdb.NewUsersRepo(usersDB)

	embedderEndpoint := fmt.Sprintf("http://localhost:%s/embedd", os.Getenv("EMBEDDER_SERVICE_PORT"))
	fmt.Println("embedder endpoint: ", embedderEndpoint)

	tp := transport.NewLocalTransportWithConsumers(
		transport.ConsumerMap{
			transport.Embed: embedderEndpoint,
		},
	)

	core := explore.NewExplorer(matchingRepo, usersRepo, redisClient)
	service := exploreapi.NewService(core, redisClient, tp)

	adminJSON, _ := utils.GetFile("firebase.admin.json")
	auth, err := auth.NewFirebaseManager(adminJSON)
	if err != nil {
		panic(err)
	}

	manager = exploreapi.NewManager(nil, auth, usersRepo, service)
}

// random list of major
var majors = []string{
	"Computer Science",
	"Mathematics",
	"Physics",
	"Chemistry",
	"Biology",
	"Engineering",
}

// random list of language in RFC 5646 standard
var languages = []string{
	"en",
	"vi",
}

// random list of country code in ISO 3166 2 characters standard
var countries = []string{
	"US",
	"CA",
	"GB",
	"DE",
	"FR",
	"IT",
	"JP",
	"VN",
}

// list of interests
var interests = []string{
	"Music",
	"Art",
	"Sport",
	"Reading",
	"Travel",
	"Photography",
}

// list of genders
var genders = []string{"male", "female"}

func generateMatchingProfile(user *usersdb.User) matchingdb.MatchInfo {
	// generate matching profile from user record, with random filed values
	matchingProfile := matchingdb.MatchInfo{
		UserID: user.ID,
		// random age from 20-30
		Age:     20 + (rand.Intn(10)),
		Name:    user.Name,
		Major:   majors[rand.Intn(len(majors))],
		Gender:  genders[rand.Intn(len(genders))],
		Country: countries[rand.Intn(len(countries))],
		Interests: []string{
			interests[rand.Intn(len(interests))],
			interests[rand.Intn(len(interests))],
		},
		Learnings: []string{
			languages[rand.Intn(len(languages))],
		},
		Native: languages[len(languages)-1-rand.Intn(len(languages))],
	}
	return matchingProfile
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	cur, err := manager.UsersRepo.Find(ctx, bson.M{})
	if err != nil {
		log.Panicln(err)
	}
	users := []usersdb.User{}
	if err := cur.All(ctx, &users); err != nil {
		log.Panicln(err)
	}

	for _, user := range users {
		matchingProfile := generateMatchingProfile(&user)
		fmt.Println(matchingProfile)
		embedding, err := manager.Service.HandleGetEmbedding(&matchingProfile)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("fetch embedding", len(embedding), "vectors")
		info, err := manager.Service.Core.GetMatchingProfile(matchingProfile.UserID)
		if err != nil {
			log.Println("cannot get matching profile", err)
			info, err = manager.Service.Core.AddUserMatchInformation(&matchingProfile)
			if err != nil {
				log.Println("cannot add user match information", err)
			}
		} else {
			matchingProfile.SetID(info.ID)
			matchingProfile.SetInitTime(info.CreatedAt.Time())
			matchingProfile.SetUpdatedAtByNow()

			info, err = manager.Service.Core.UpdaterUserMatchInformation(&matchingProfile)
			if err != nil {
				log.Println("cannot update user match information", err)
				return
			}
		}

		if err := manager.Service.Core.AddEmbedding(info.UserID, embedding); err != nil {
			log.Println("cannot add user embed", err)
			err = manager.Service.Core.UpdateEmbedding(info.UserID, embedding)
			if err != nil {
				log.Println("cannot update user embed", err)
				return
			}
			return
		}
		fmt.Println("mock for user", user.ID.Hex(), "done!")
	}
}
