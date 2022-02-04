package studydb

import (
	"errors"

	"github.com/influenzanet/study-service/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (dbService *StudyDBService) AddConfidentialResponse(instanceID string, studyKey string, response types.SurveyResponse) (string, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()
	if len(response.ParticipantID) < 1 {
		return "", errors.New("participantID must be defined")
	}
	res, err := dbService.collectionRefConfidentialResponses(instanceID, studyKey).InsertOne(ctx, response)
	id := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), err
}

func (dbService *StudyDBService) ReplaceConfidentialResponse(instanceID string, studyKey string, response types.SurveyResponse) error {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{
		"participantID": response.ParticipantID,
		"key":           response.Key,
	}

	upsert := true
	options := options.ReplaceOptions{
		Upsert: &upsert,
	}
	_, err := dbService.collectionRefConfidentialResponses(instanceID, studyKey).ReplaceOne(ctx, filter, response, &options)
	return err
}

func (dbService *StudyDBService) FindConfidentialResponses(instanceID string, studyKey string, participantID string, key string) (responses []types.SurveyResponse, err error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	if participantID == "" {
		return responses, errors.New("participant id must be defined")
	}
	filter := bson.M{"participantID": participantID}
	if key != "" {
		filter["key"] = key
	}

	cur, err := dbService.collectionRefConfidentialResponses(instanceID, studyKey).Find(
		ctx,
		filter,
		nil,
	)

	if err != nil {
		return responses, err
	}
	defer cur.Close(ctx)

	responses = []types.SurveyResponse{}
	for cur.Next(ctx) {
		var result types.SurveyResponse
		err := cur.Decode(&result)
		if err != nil {
			return responses, err
		}

		responses = append(responses, result)
	}
	if err := cur.Err(); err != nil {
		return responses, err
	}

	return responses, nil
}

func (dbService *StudyDBService) DeleteConfidentialResponses(instanceID string, studyKey string, participantID string, key string) (count int64, err error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	if participantID == "" {
		return 0, errors.New("participant id must be defined")
	}
	filter := bson.M{"participantID": participantID}
	if key != "" {
		filter["key"] = key
	}

	res, err := dbService.collectionRefConfidentialResponses(instanceID, studyKey).DeleteMany(ctx, filter)
	return res.DeletedCount, err
}

func (dbService *StudyDBService) UpdateParticipantIDonConfidentialResponses(instanceID string, studyKey string, oldID string, newID string) (count int64, err error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	if oldID == "" || newID == "" {
		return 0, errors.New("participant id must be defined")
	}
	filter := bson.M{"participantID": oldID}
	update := bson.M{"$set": bson.M{"participantID": newID}}

	res, err := dbService.collectionRefConfidentialResponses(instanceID, studyKey).UpdateMany(ctx, filter, update)
	return res.ModifiedCount, err
}
