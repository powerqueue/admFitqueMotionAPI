package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/powerqueue/admFitqueMotionAPI/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var (
	audience string
	domain   string
)

func main() {

	// connect to mongo and migrate schema
	dbConfigs := &models.DBConfigs1{
		Host:     "ln004prd",
		Port:     27017,
		User:     "",
		Password: "",
	}
	models.ConnectAndMigrate(dbConfigs, "fitqueue-db")

	// apiPrefix := "/fitqueue-login-api/v1"

	r := gin.Default()

	r.Use(CORSMiddleware())

	// r.GET("/todo", GetTodoListHandler)
	// r.POST("/todo", AddTodoHandler)
	// r.DELETE("/todo/:id", DeleteTodoHandler)
	// r.PUT("/todo", CompleteTodoHandler)

	v1 := r.Group("/fitqueue-motion-api/v1")
	{
		v1.POST("/retrieve-motion", RetrieveMotionHandler)
		v1.POST("/create-motion", CreateMotion)
	}

	// authorized := r.Group("/")
	// authorized.Use(authRequired())
	// authorized.GET("/cases/:caseType/:page", GetCaseListHandler)
	// authorized.GET("/case/:caseCode", GetCaseDetailsHandler)
	// authorized.POST("/case", AddCaseHandler)
	// authorized.DELETE("/case/:caseCode", TermCaseHandler)
	// authorized.PUT("/case", UpdateCaseHandler)

	err := r.Run(":3101")
	if err != nil {
		panic(err)
	}
}

//CORSMiddleware Cross-Origin Resource Sharing helper Class
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, GET, OPTIONS, POST, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func terminateWithError(statusCode int, message string, c *gin.Context) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}

func convertHTTPBodyToMotionDefinition(httpBody io.ReadCloser) (models.MotionDefinition, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return models.MotionDefinition{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	return convertJSONBodyToMotionDefinition(body)
}

func convertJSONBodyToMotionDefinition(jsonBody []byte) (models.MotionDefinition, int, error) {
	var motion models.MotionDefinition
	err := json.Unmarshal(jsonBody, &motion)
	if err != nil {
		return models.MotionDefinition{}, http.StatusBadRequest, err
	}
	return motion, http.StatusOK, nil
}

//RetrieveMotionHandler - handler definition
func RetrieveMotionHandler(c *gin.Context) {
	fmt.Println("Inside RetrieveMotionHandler Route Handler")
	motion, statusCode, err := convertHTTPBodyToMotionDefinition(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("fitque-db").Collection("motion")

	// create a value into which the result can be decoded
	filter := bson.D{{"LocationID", motion.LocationID}, {"SensorID", motion.SensorID}, {"MotionStartDt", motion.MotionStartDt}}

	err = collection.FindOne(context.TODO(), filter).Decode(&motion)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", motion)

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

	c.JSON(statusCode, motion)
}

//CreateMotion - handler method definition
func CreateMotion(c *gin.Context) {
	fmt.Println("Inside CreateMotion Route Handler")
	motion, statusCode, err := convertHTTPBodyToMotionDefinition(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}

	fmt.Println("Inside CreateMotion Model Function")
	motion.ID = primitive.NewObjectID()
	motion.LocationID = strings.ToUpper(motion.LocationID)
	motion.SensorID = strings.ToUpper(motion.SensorID)

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://ln004prd:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error create %s", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error create %s", err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("fitqueue-db").Collection("motion")

	//Insert SINGLE here
	// insertResult, err := collection.InsertOne(context.TODO(), ash)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	//Insert MULTIPLE here
	motions := []interface{}{motion}

	insertManyResult, err := collection.InsertMany(context.TODO(), motions)
	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error create %s", err)
	}

	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	//disconnect client
	err = client.Disconnect(context.TODO())

	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error create %s", err)
	}
	fmt.Println("Connection to MongoDB closed.")

	bytes, err := json.Marshal(motion)
	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error create %s", err)

	}

	fmt.Println(string(bytes))
	// return login, err

	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "Category: %v\n", loginDef)

	c.JSON(statusCode, motion)
}
