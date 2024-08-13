package repositories

import (
	"context"
	domain "example/go-clean-architecture/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type taskRepository struct {
	database   mongo.Database
	collection string
}

func NewTaskRepository(db mongo.Database, collection string) domain.TaskRepository {
	return &taskRepository{
		database:   db,
		collection: collection,
	}
}

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

func (tr *taskRepository) FindTaskById(ctx context.Context, taskId string) (domain.Task, error){
	collection := tr.database.Collection(tr.collection)
	objID, err := primitive.ObjectIDFromHex(taskId)
	if err != nil{
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

func (tr *taskRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error){
	collection := tr.database.Collection(tr.collection)
	task.ID = primitive.NewObjectID()
    // Insert the task into the collection
    _, err := collection.InsertOne(ctx, task)
    if err != nil {
        return domain.Task{}, err
	}
    return task, nil
}

func (tr *taskRepository) UpdateTaskById(ctx context.Context, updatedTask domain.Task, id string) (domain.Task, error){
	collection := tr.database.Collection(tr.collection)
	update := bson.M{
        "$set": bson.M{
            "title":       updatedTask.Title,
            "description": updatedTask.Description,
            "due_date":    updatedTask.DueDate,
            "status":      updatedTask.Status,
        },
    }

    filter := bson.M{"_id": id}
    opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
    var updated domain.Task
    err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return domain.Task{}, err
        }
        return domain.Task{}, err
    }
    return updated, nil
}

func (tr *taskRepository) DeleteTask(ctx context.Context, taskId string) error{
	collection := tr.database.Collection(tr.collection)
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": taskId})
	if err != nil{
		return err
	}
	return nil
}





