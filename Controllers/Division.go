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

func CreateNewDivision(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Division

	var self Models.Division
	c.BodyParser(&self)

	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}

	self.Status=true
	self.Zilas=[]primitive.ObjectID{}
	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}

	Responses.Created(c, "Division", res)
	return nil
}

func GetAllDivision(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Division

	// Fill the received search obj data
	var self Models.DivisionSearch
	c.BodyParser(&self)
	opts := options.Find().SetProjection(bson.D{{"_id",0},{"zilas", 0},{"status",0}})
	b, results := Utils.FindByFilter(collection, self.GetDivisionSearchBSONObj(),opts)
	if !b {
		return Responses.NotFound(c, "Division")
	}

	Responses.Get(c, "Division", results)
	return nil
}

func GetAllDivisionNames(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Division

	// Fill the received search obj data
	var self Models.DivisionSearch
	c.BodyParser(&self)
	opts := options.Find().SetProjection(bson.D{})
	b, results := Utils.FindByFilter(collection, self.GetDivisionSearchBSONObj(),opts)
	if !b {
		return Responses.NotFound(c, "Division")
	}

	byteArr, _ := json.Marshal(results)
	var ResultDocs []Models.Division
	json.Unmarshal(byteArr, &ResultDocs)
	populatedResult := make([]string, len(ResultDocs))

	for i, v := range ResultDocs {
		populatedResult[i], _ = GetDivisionName(v.ID, &v)
	}

	Responses.Get(c, "Division", populatedResult)
	return nil
}

func GetDivisionName(objID primitive.ObjectID, ptr *Models.Division) (string, error) {
	var divisionDoc Models.Division
	if ptr == nil {
		divisionDoc, _ = GetDivisionById(objID)
	} else {
		divisionDoc = *ptr
	}

	return divisionDoc.Name, nil
}


func GetAllDivisionByPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Division

	// Fill the received search obj data
	var self Models.DivisionSearch
	c.BodyParser(&self)
	opts := options.Find().SetProjection(bson.D{{"_id", 0}})
	b, results := Utils.FindByFilter(collection, self.GetDivisionSearchBSONObj(),opts)
	if !b {
		return Responses.NotFound(c, "Division")
	}
	byteArr, _ := json.Marshal(results)
	var ResultDocs []Models.Division
	json.Unmarshal(byteArr, &ResultDocs)
	populatedResult := make([]Models.DivisionPopulated, len(ResultDocs))

	for i, v := range ResultDocs {
		populatedResult[i], _ = GetDivisionByIdPopulated(v.ID, &v)
	}
	
	Responses.Get(c, "Division", populatedResult)
	return nil
}

func GetDivisionByIdPopulated(objID primitive.ObjectID, ptr *Models.Division) (Models.DivisionPopulated, error) {
	var divisionDoc Models.Division
	if ptr == nil {
		divisionDoc, _ = GetDivisionById(objID)
	} else {
		divisionDoc = *ptr
	}

	populatedResult := Models.DivisionPopulated{}
	populatedResult.CloneFrom(divisionDoc)

	// populate for SOF array
	populatedResult.Zilas = make([]Models.Zila, len(divisionDoc.Zilas))
	for i, element := range divisionDoc.Zilas {
		populatedResult.Zilas[i], _ = GetZilaById(element)
	}

	return populatedResult, nil
}



func GetDivisionById(id primitive.ObjectID) (Models.Division, error) {
	collection := DBManager.SystemCollections.Division

	filter := bson.M{"_id": id}
	var self Models.Division
	opts := options.Find().SetProjection(bson.D{})
	_, results := Utils.FindByFilter(collection, filter,opts)
	if len(results) == 0 {
		return self, errors.New("Division not found")
	}

	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}
func GetDivisionByName(name string) (Models.Division, error) {
	collection := DBManager.SystemCollections.Division
	filter := bson.M{"name": name}
	var self Models.Division
	opts := options.Find().SetProjection(bson.D{})
	_, results := Utils.FindByFilter(collection, filter,opts)
	if len(results) == 0 {
		return self, errors.New("Division not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}