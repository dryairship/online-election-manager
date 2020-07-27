package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/dryairship/online-election-manager/models"
)

// Function to get all posts from the database.
func (db ElectionDatabase) GetPosts() ([]models.Post, error) {
	var posts []models.Post

	cursor, err := db.PostsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return posts, err
	}

	err = cursor.All(context.Background(), &posts)
	return posts, err
}

// Function to insert a new post into the database.
func (db ElectionDatabase) AddNewPost(post *models.Post) error {
	_, err := db.PostsCollection.InsertOne(context.Background(), post)
	return err
}

func (db ElectionDatabase) DeletePost(postId string) error {
	_, err := db.PostsCollection.DeleteOne(context.Background(), bson.M{"postid": postId})
	return err
}

func (db ElectionDatabase) MarkPostResolved(postId string) error {
	_, err := db.PostsCollection.UpdateOne(context.Background(),
		bson.M{"postid": postId}, bson.M{"$set": bson.M{"resolved": true}})
	return err
}

func (db ElectionDatabase) GetPostsForCEO() ([]models.Post, error) {
	var posts []models.Post

	cursor, err := db.PostsCollection.Find(context.Background(), bson.M{"resolved": false})
	if err != nil {
		return posts, err
	}

	err = cursor.All(context.Background(), &posts)
	return posts, err
}
