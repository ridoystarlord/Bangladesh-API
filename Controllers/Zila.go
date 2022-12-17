package Controllers

import (
	"bangladesh-api/DBManager"
	"bangladesh-api/Models"
	"bangladesh-api/Responses"
	"bangladesh-api/Utils"
	"context"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateNewZila(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Zila
	divisionCollection := DBManager.SystemCollections.Division

	var self Models.Zila
	c.BodyParser(&self)

	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}

	self.Status=true
	self.Upazilas=[]primitive.ObjectID{}
	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}
	divisionDoc, _ := GetDivisionById(self.ParentId)
	newId := res.InsertedID

	updateData := bson.M{
		"$push": bson.M{
			"zilas": newId,
		},
	}

	_, err = divisionCollection.UpdateOne(context.Background(), bson.M{"_id": divisionDoc.ID}, updateData)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}

	Responses.Created(c, "Zila", res)
	return nil
}

func GetAllZila(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Zila

	// Fill the received search obj data
	var self Models.ZilaSearch
	c.BodyParser(&self)
	opts := options.Find().SetProjection(bson.D{{"_id", 0},{"upazilas", 0},{"status",0},{"parentid",0}})
	b, results := Utils.FindByFilter(collection, self.GetZilaSearchBSONObj(),opts)
	if !b {
		return Responses.NotFound(c, "Zila")
	}
	Responses.Get(c, "Zila", results)
	return nil
}

func GetAllZilaNames(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Zila

	// Fill the received search obj data
	var self Models.ZilaSearch
	c.BodyParser(&self)
	opts := options.Find().SetProjection(bson.D{})
	b, results := Utils.FindByFilter(collection, self.GetZilaSearchBSONObj(),opts)
	if !b {
		return Responses.NotFound(c, "Zila")
	}

	byteArr, _ := json.Marshal(results)
	var ResultDocs []Models.Zila
	json.Unmarshal(byteArr, &ResultDocs)
	populatedResult := make([]string, len(ResultDocs))

	for i, v := range ResultDocs {
		populatedResult[i], _ = GetZilaName(v.ID, &v)
	}

	Responses.Get(c, "Zila", populatedResult)
	return nil
}

func GetZilaName(objID primitive.ObjectID, ptr *Models.Zila) (string, error) {
	var zilaDoc Models.Zila
	if ptr == nil {
		zilaDoc, _ = GetZilaById(objID)
	} else {
		zilaDoc = *ptr
	}

	return zilaDoc.Name, nil
}

func GetZilaById(id primitive.ObjectID) (Models.Zila, error) {
	collection := DBManager.SystemCollections.Zila

	filter := bson.M{"_id": id}
	var self Models.Zila
	opts := options.Find().SetProjection(bson.D{{"_id",0}})
	_, results := Utils.FindByFilter(collection, filter,opts)
	if len(results) == 0 {
		return self, errors.New("Zila not found")
	}

	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func GetAllZilaNamesByDivision(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Zila
	divisionName:= c.Params("divisionName")
	division,_ := GetDivisionByName(divisionName)
	filter := bson.M{"parentid":division.ID}
	opts := options.Find().SetProjection(bson.D{})
	b, results := Utils.FindByFilter(collection, filter, opts)
	if !b {
		return Responses.NotFound(c, "Zila")
	}

	byteArr, _ := json.Marshal(results)
	var ResultDocs []Models.Zila
	json.Unmarshal(byteArr, &ResultDocs)
	populatedResult := make([]string, len(ResultDocs))

	for i, v := range ResultDocs {
		populatedResult[i], _ = GetZilaName(v.ID, &v)
	}

	Responses.Get(c, divisionName+" Divison Zila", populatedResult)
	return nil
}