package explore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"blinders/packages/db/matchingdb"
	"blinders/packages/db/usersdb"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Explorer interface {
	// SuggestWithContext returns list of users that maybe match with given user
	SuggestWithContext(userID string) ([]matchingdb.MatchInfo, error)
	// AddUserMatchInformation adds user match information to the database.
	AddUserMatchInformation(info matchingdb.MatchInfo) (matchingdb.MatchInfo, error)
	// AddEmbedding adds user embed vector to the vector database.
	AddEmbedding(userID primitive.ObjectID, embed EmbeddingVector) error
	// SuggestRandom returns list of random 5 users that maybe match with given user
	SuggestRandom() ([]matchingdb.MatchInfo, error)
}

type MongoExplorer struct {
	MatchingRepo *matchingdb.MatchingRepo
	UsersRepo    *usersdb.UsersRepo
	RedisClient  *redis.Client
}

func NewExplorer(
	matchingRepo *matchingdb.MatchingRepo,
	usersRepo *usersdb.UsersRepo,
	redisClient *redis.Client,
) *MongoExplorer {
	return &MongoExplorer{
		MatchingRepo: matchingRepo,
		UsersRepo:    usersRepo,
		RedisClient:  redisClient,
	}
}

/*
Suggest  recommends 5 users who are not friends of the current user.

TODO: The goal is to recommend users with whom the current user may communicate effectively.
These users should either be fluent in the language the current user is learning or actively learning the same language.
To achieve this, we will filter the Users database to extract users who are native speakers of the language the current user is learning,
or users who are currently learning the same language as the current user.

We will then use KNN-search in the filtered space to identify 5 users that may match with the current user.
*/
func (m *MongoExplorer) SuggestWithContext(userID string) ([]matchingdb.MatchInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf(
			"explore: cannot parse object id with given hex string(%s), err: %v\n",
			userID,
			err,
		)
		return nil, err
	}

	user, err := m.UsersRepo.GetUserByID(oid)
	if err != nil {
		log.Println("explore: cannot get user by id, err:", err)
		return nil, err
	}

	// JSONGet return value wrapped in an array.
	// at here, if there aren't entries in redis, the jsonStr will be empty, we could check it here then return
	jsonStr, err := m.RedisClient.JSONGet(ctx, CreateMatchKeyWithUserID(userID), "$.embed").Result()
	if err != nil || jsonStr == "" {
		log.Println("explore: cannot get explore entry in redis, err:", err)
		return []matchingdb.MatchInfo{}, fmt.Errorf(
			"explore profile not found, might need to check onboarding status",
		)
	}

	var embedArr []EmbeddingVector
	if err := json.Unmarshal([]byte(jsonStr), &embedArr); err != nil {
		log.Println("explore: cannot unmarshall embed vector, err:", err)
		return []matchingdb.MatchInfo{}, fmt.Errorf("something went wrong")
	}
	embed := embedArr[0]

	// exclude friends of current user
	excludeFilter := userID
	for _, friendID := range user.FriendIDs {
		excludeFilter += " | " + friendID.Hex()
	}
	excludeFilter = fmt.Sprintf("-@id:(%s)", excludeFilter)

	candidates, err := m.MatchingRepo.GetUsersByLanguage(user.ID, 1000)
	if err != nil {
		log.Println("explore: cannot explore candidates, err:", err)
		return nil, err
	}

	includeFilter := ""
	if len(candidates) != 0 {
		includeFilter = candidates[0]
		for idx := 1; idx < len(candidates); idx++ {
			includeFilter += " | " + candidates[idx]
		}
		includeFilter = fmt.Sprintf("@id:(%s)", includeFilter)
	}

	prefilter := fmt.Sprintf("(%s %s)", excludeFilter, includeFilter)

	cmd := m.RedisClient.Do(ctx,
		"FT.SEARCH",
		"idx:match_vss",
		fmt.Sprintf("%s=>[KNN 5 @embed $query_vector as vector_score]", prefilter),
		"SORTBY", "vector_score",
		"PARAMS", "2",
		"query_vector", &embed,
		"DIALECT", "2",
		"RETURN", "1", "id",
	)
	if err := cmd.Err(); err != nil {
		log.Println("explore: cannot perform knn search in vector database, err:", err)
		return nil, err
	}

	var res []matchingdb.MatchInfo
	for _, doc := range cmd.Val().(map[any]any)["results"].([]any) {
		userID := doc.(map[any]any)["extra_attributes"].(map[any]any)["id"].(string)
		oid, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, err
		}
		user, err := m.MatchingRepo.GetByUserID(oid)
		if err != nil {
			return nil, err
		}
		res = append(res, *user)
	}

	// TODO: After the suggestion process, mark these users as suggested to prevent them from being recommended in future suggestions.
	// Idea: Recommended users will be assigned extra points, which will be added to their vector space during the vector search, making their vectors more distant from the current vector.
	// Redis does not support sorting by expression.
	return res, nil
}

/*
AddUserMatchInformation inserts information into the match database.

Currently, embedding will be handled by another service. The caller of this method must trigger a new event
to notify that a new user has been created. This allows the embedding service to update the embedding vector
in the vector database.
*/
func (m *MongoExplorer) AddUserMatchInformation(
	info matchingdb.MatchInfo,
) (matchingdb.MatchInfo, error) {
	_, err := m.UsersRepo.GetUserByID(info.UserID)
	if err != nil {
		return matchingdb.MatchInfo{}, err
	}

	// duplicated match information will be handled by the repository since we have already indexed the collection with firebaseUID.
	matchInfo, err := m.MatchingRepo.InsertRaw(&info)
	if err != nil {
		return matchingdb.MatchInfo{}, err
	}
	return *matchInfo, nil
}

func (m *MongoExplorer) AddEmbedding(userID primitive.ObjectID, embed EmbeddingVector) error {
	_, err := m.MatchingRepo.GetByUserID(userID)
	if err != nil {
		return err
	}
	fmt.Println("checkpoint")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = m.RedisClient.JSONSet(ctx,
		CreateMatchKeyWithUserID(userID.Hex()),
		"$",
		map[string]any{"embed": embed, "id": userID},
	).Err()
	return err
}

func (m *MongoExplorer) SuggestRandom() ([]matchingdb.MatchInfo, error) {
	return nil, nil
}
