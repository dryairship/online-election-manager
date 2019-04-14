package db

import (
    "github.com/dryairship/online-election-manager/models"
    "github.com/dryairship/online-election-manager/config"
)

// Function to get all posts from the database.
func (db ElectionDatabase) GetPosts() ([]models.Post, error){
    postsCollection := db.Session.DB(config.MongoDbName).C("posts")
    posts := []models.Post{}
    err := postsCollection.Find(nil).All(&posts)
    return posts, err
}

// Function to insert a new post into the database.
func (db ElectionDatabase) AddNewPost(post *models.Post) error {
    postsCollection := db.Session.DB(config.MongoDbName).C("posts")
    err := postsCollection.Insert(&post)
    return err
}
