package Utils

import (
	"context"
	"encoding/base64"
	"math"
	"os"
	"strconv"
	"strings"

	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindByFilter(collection *mongo.Collection, filter bson.M,opts *options.FindOptions) (bool, []bson.M) {
	var results []bson.M
	
	cur, err := collection.Find(context.Background(), filter,opts)
	if err != nil {
		return false, results
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)

	return true, results
}

func CollectionGetById(col interface{}, objID primitive.ObjectID, self interface{}) error {
	var filter bson.M = bson.M{}
	filter = bson.M{"_id": objID}
	collection := col.(*mongo.Collection)
	var results []bson.M
	opts := options.Find().SetProjection(bson.D{})
	b, results := FindByFilter(collection, filter,opts)
	if !b || len(results) == 0 {
		return errors.New("Object not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return nil
}

func ArrayStringContains(arr []string, elem string) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func UploadImage(c *fiber.Ctx) (string, error) {
	file, err := c.FormFile("image")
	if err != nil {
		return "", err
	}

	// Save file to root directory
	var filePath = fmt.Sprintf("images/img_%d_%d.png", rand.Intn(1024), MakeTimestamp())
	saving_err := c.SaveFile(file, "./public/"+filePath)
	if saving_err != nil {
		return "", saving_err
	} else {
		c.Status(200).Send([]byte("Saved Successfully"))
		return filePath, nil
	}
}

func UtilsUploadImageBase64(stringBase64 string, imageDocType string, imageName string) (string, error) {
	i := strings.Index(stringBase64, ",")
	if i != -1 {
		file, _ := base64.StdEncoding.DecodeString(stringBase64[i+1:])
		var filePath = fmt.Sprintf("Public/Images/"+imageName+"_img_%d_%d.%s", rand.Intn(1024), MakeTimestamp(), imageDocType)

		f, err := os.Create("./" + filePath)
		if err != nil {
			return "", err
		}
		defer f.Close()

		if _, err := f.Write(file); err != nil {
			return "", err
		}
		f.Sync()
		return filePath, nil
	}
	return "", nil
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func FindByFilterProjected(collection *mongo.Collection, filter bson.M, fields bson.M) ([]bson.M, error) {
	var results []bson.M
	opts := options.FindOptions{Projection: fields}
	cur, err := collection.Find(context.Background(), filter, &opts)
	if err != nil {
		return results, err
	}
	defer cur.Close(context.Background())

	cur.All(context.Background(), &results)

	return results, err
}

func HashPassword(password string) string {
	return fmt.Sprintf("%X", sha256.Sum256([]byte(password)))
}

func Contains(arr []primitive.ObjectID, elem primitive.ObjectID) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func ContainsItem(arr []primitive.ObjectID, elem primitive.ObjectID) (bool, int) {
	for i, v := range arr {
		if v == elem {
			return true, i
		}
	}
	return false, -1
}

func AnnotateMoney(total float64) string {
	if total > 999999999 {
		result := total / 1000000000
		return fmt.Sprintf("%.2f"+" B", result)
	} else if total > 999999 && total < 999999999 {
		result := total / 1000000
		return fmt.Sprintf("%.2f"+" M", result)
	} else if total > 999 && total < 999999 {
		result := total / 1000
		return fmt.Sprintf("%.2f"+" K", result)
	}

	return fmt.Sprintf("%.2f", total)
}

func AnnotateTimePeriod(invoiceDate primitive.DateTime) string {
	oldestInvoiceDate := invoiceDate.Time()
	currentDay := time.Now()
	diff := math.Ceil(math.Abs(currentDay.Sub(oldestInvoiceDate).Hours() / 24))

	if diff >= 365 {
		year := diff / 365
		str := fmt.Sprintf("%.2f", year)
		date := fmt.Sprintf(str + " Y")
		return date
	} else if diff < 365 && diff >= 30 {
		month := diff / 30
		str := fmt.Sprintf("%.2f", month)
		date := fmt.Sprintf(str + " M")
		return date
	}

	date := fmt.Sprintf(strconv.Itoa(int(diff)) + " D")
	return date
}
