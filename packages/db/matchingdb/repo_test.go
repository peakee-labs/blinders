package matchingdb

import (
	"context"
	"slices"
	"testing"
	"time"

	dbutils "blinders/packages/db/utils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	MongoTestURL = "mongodb://localhost:27017"
	MongoTestDB  = "blinder-test"
)

func TestInsertNewRawMatchInfo(t *testing.T) {
	r := GetTestRepo(t)
	defer CleanRepo(t, r)
	rawUser := MatchInfo{
		UserID:    primitive.NewObjectID(),
		Name:      "name",
		Gender:    "male",
		Major:     "student",
		Native:    "vietnamese",
		Country:   "vn",
		Learnings: make([]string, 0),
		Interests: make([]string, 0),
		Age:       0,
	}
	insertedUser, err := r.InsertRaw(&rawUser)
	assert.Nil(t, err)
	assert.Equal(t, rawUser.UserID, insertedUser.UserID)
	assert.Equal(t, rawUser.Gender, insertedUser.Gender)
	assert.Equal(t, rawUser.Major, insertedUser.Major)
	assert.Equal(t, rawUser.Native, insertedUser.Native)
	assert.Equal(t, rawUser.Country, insertedUser.Country)
	assert.Equal(t, rawUser.Learnings, insertedUser.Learnings)
	assert.Equal(t, rawUser.Interests, insertedUser.Interests)
	assert.Equal(t, rawUser.Age, insertedUser.Age)

	gotWithUserID, err := r.GetByUserID(rawUser.UserID)
	assert.Nil(t, err)
	assert.Equal(t, *insertedUser, *gotWithUserID)

	deleted, err := r.DropMatchInfoByUserID(rawUser.UserID)
	assert.Nil(t, err)
	assert.Equal(t, *insertedUser, *deleted)
}

func TestGetMatchInfoByUserID(t *testing.T) {
	r := GetTestRepo(t)
	defer CleanRepo(t, r)
	rawUser := MatchInfo{
		UserID:    primitive.NewObjectID(),
		Name:      "name",
		Gender:    "male",
		Major:     "student",
		Native:    "vietnamese",
		Country:   "vn",
		Learnings: make([]string, 0),
		Interests: make([]string, 0),
		Age:       0,
	}
	insertedUser, err := r.InsertRaw(&rawUser)
	assert.Nil(t, err)
	assert.NotNil(t, insertedUser)

	gotWithUserID, err := r.GetByUserID(rawUser.UserID)
	assert.Nil(t, err)
	assert.Equal(t, *insertedUser, *gotWithUserID)

	deleted, err := r.DropMatchInfoByUserID(rawUser.UserID)
	assert.Nil(t, err)
	assert.Equal(t, *insertedUser, *deleted)

	gotFailed, err := r.GetByUserID(rawUser.UserID)
	assert.NotNil(t, err)
	assert.Nil(t, gotFailed)
}

func TestGetUsersByLanguage(t *testing.T) {
	r := GetTestRepo(t)
	defer CleanRepo(t, r)
	rawUser := MatchInfo{
		UserID:    primitive.NewObjectID(),
		Name:      "name",
		Gender:    "male",
		Major:     "student",
		Native:    "vietnamese",
		Country:   "vn",
		Learnings: make([]string, 0),
		Interests: make([]string, 0),
		Age:       0,
	}
	numReturn := uint32(10)

	deletedUser, err := r.DropMatchInfoByUserID(rawUser.UserID)
	if err != nil {
		assert.Nil(t, deletedUser)
	} else {
		assert.NotNil(t, deletedUser)
	}

	failedGot, err := r.GetUsersByLanguage(rawUser.UserID, 10)
	assert.NotNil(t, err)
	assert.Len(t, failedGot, 0)

	insertedUser, err := r.InsertRaw(&rawUser)
	assert.Nil(t, err)
	assert.NotNil(t, insertedUser)

	got, err := r.GetUsersByLanguage(rawUser.UserID, numReturn)
	assert.Nil(t, err)

	assert.GreaterOrEqual(t, numReturn, uint32(len(got)))

candidateLoop:
	for _, id := range got {
		oid, err := primitive.ObjectIDFromHex(id)
		assert.Nil(t, err)
		assert.False(t, oid.IsZero())

		candidate, err := r.GetByUserID(oid)
		assert.Nil(t, err)
		assert.NotNil(t, candidate)
		// at here, candidate must be learning same language with curr user or natively speak the language that current
		// user is learning as well as learning language that current user is natively speak.
		for _, language := range candidate.Learnings {
			if slices.Contains[[]string, string](insertedUser.Learnings, language) {
				// user and candidate learning same language
				continue candidateLoop
			}
		}
		assert.Contains(t, insertedUser.Learnings, candidate.Native)
		assert.Contains(t, candidate.Learnings, insertedUser.Native)
	}
	dropUser, err := r.DropMatchInfoByUserID(rawUser.UserID)
	assert.Nil(t, err)
	assert.Equal(t, *insertedUser, *dropUser)
}

func TestDropUserByUserID(t *testing.T) {
	r := GetTestRepo(t)
	defer CleanRepo(t, r)
	rawUser := MatchInfo{
		UserID:    primitive.NewObjectID(),
		Name:      "name",
		Gender:    "male",
		Major:     "student",
		Native:    "vietnamese",
		Country:   "vn",
		Learnings: make([]string, 0),
		Interests: make([]string, 0),
		Age:       0,
	}
	insertedUser, err := r.InsertRaw(&rawUser)
	assert.Nil(t, err)
	assert.NotNil(t, insertedUser)

	deleted, err := r.DropMatchInfoByUserID(insertedUser.UserID)
	assert.Nil(t, err)
	assert.Equal(t, *insertedUser, *deleted)

	failed, err := r.DropMatchInfoByUserID(insertedUser.UserID)
	assert.NotNil(t, err)
	assert.Nil(t, failed)
}

func CleanRepo(t *testing.T, repo *MatchingRepo) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	assert.Nil(t, repo.Drop(ctx))
}

func GetTestRepo(t *testing.T) *MatchingRepo {
	client, err := dbutils.InitMongoClient(MongoTestURL)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	assert.Nil(t, client.Ping(ctx, nil))

	return NewMatchingRepo(client.Database(MongoTestDB))
}
