package Repositories

import (
	"context"
	domain "example/go-clean-architecture/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// taskRepository represents a repository for managing tasks.
type taskRepository struct {
	database   mongo.Database
	collection string
}

// NewTaskRepository creates a new instance of the TaskRepository interface.
// It takes a mongo.Database and a collection name as parameters.
// It returns a pointer to a taskRepository struct that implements the TaskRepository interface.
func NewTaskRepository(db mongo.Database, collection string) domain.TaskRepository {
	return &taskRepository{
		database:   db,
		collection: collection,
	}
}

// FindAlltasks retrieves all tasks from the task repository.
// It returns a slice of domain.Task and an error if any.
func (tr *taskRepository) FindAlltasks(ctx context.Context) ([]domain.Task, error) {
	collection := tr.database.Collection(tr.collection)
	var tasks []domain.Task
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task domain.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// FindTaskById retrieves a task from the database by its ID.
// It takes a context.Context and a taskId string as parameters.
// It returns a domain.Task and an error.
// The domain.Task represents the retrieved task from the database.
// The error is returned if there was an issue retrieving the task.
func (tr *taskRepository) FindTaskById(ctx context.Context, taskId string) (domain.Task, error) {
	collection := tr.database.Collection(tr.collection)
	objID, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return domain.Task{}, err
	}
	filter := bson.M{"_id": objID}

	var task domain.Task
	err = collection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Task{}, err // No document found
		}
		return domain.Task{}, err
	}
	return task, nil
}

// CreateTask creates a new task in the task repository.
// It takes a context and a task object as parameters.
// It returns the created task and an error, if any.
func (tr *taskRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	collection := tr.database.Collection(tr.collection)
	task.ID = primitive.NewObjectID()
	// Insert the task into the collection
	_, err := collection.InsertOne(ctx, task)
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

// UpdateTaskById updates a task in the database with the given ID.
// It takes the updatedTask object containing the new values for the task fields,
// the id string representing the ID of the task to be updated.
// It returns the updated task object and an error if any occurred.
func (tr *taskRepository) UpdateTaskById(ctx context.Context, updatedTask domain.Task, id string) (domain.Task, error) {
	collection := tr.database.Collection(tr.collection)
	update := bson.M{
		"$set": bson.M{
			"title":       updatedTask.Title,
			"description": updatedTask.Description,
			"due_date":    updatedTask.DueDate,
			"status":      updatedTask.Status,
		},
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Task{}, err
	}
	filter := bson.M{"_id": objID}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated domain.Task
	err = collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Task{}, err
		}
		return domain.Task{}, err
	}
	return updated, nil
}

// DeleteTask deletes a task from the database.
// It takes a context.Context and a taskId string as parameters.
// It returns an error if there was a problem deleting the task.
func (tr *taskRepository) DeleteTask(ctx context.Context, taskId string) error {
	collection := tr.database.Collection(tr.collection)
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": taskId})
	if err != nil {
		return err
	}
	return nil
}
