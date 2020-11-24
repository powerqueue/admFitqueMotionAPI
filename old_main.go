package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"path"
// 	"path/filepath"
// 	"strconv"
// 	"sync"

// 	"github.com/gin-gonic/gin"
// 	"github.com/rs/xid"

// 	auth0 "github.com/auth0-community/auth0-go"
// 	jose "gopkg.in/square/go-jose.v2"

// 	"github.com/jinzhu/gorm"
// 	_ "github.com/jinzhu/gorm/dialects/mssql"
// )

// var (
// 	audience string
// 	domain   string
// )

// func main() {
// 	setAuth0Variables()
// 	initialMigration()

// 	r := gin.Default()

// 	r.Use(CORSMiddleware())

// 	r.NoRoute(func(c *gin.Context) {
// 		dir, file := path.Split(c.Request.RequestURI)
// 		ext := filepath.Ext(file)
// 		if file == "" || ext == "" {
// 			c.File("./ui/dist/ui/index.html")
// 		} else {
// 			c.File("./ui/dist/ui/" + path.Join(dir, file))
// 		}
// 	})

// 	// r.GET("/todo", GetTodoListHandler)
// 	// r.POST("/todo", AddTodoHandler)
// 	// r.DELETE("/todo/:id", DeleteTodoHandler)
// 	// r.PUT("/todo", CompleteTodoHandler)

// 	r.POST("/login", LoginUserHandler)

// 	authorized := r.Group("/")
// 	authorized.Use(authRequired())
// 	authorized.GET("/cases/:caseType/:page", GetCaseListHandler)
// 	authorized.GET("/case/:caseCode", GetCaseDetailsHandler)
// 	authorized.POST("/case", AddCaseHandler)
// 	authorized.DELETE("/case/:caseCode", TermCaseHandler)
// 	authorized.PUT("/case", UpdateCaseHandler)

// 	err := r.Run(":3001")
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func initialMigration() {
// 	db, err := gorm.Open("mssql", "sqlserver://PwrQ_DB_adm:f-azLbi4@192.168.1.166:1433?database=PowerQueueDB-2018-10-27-11-29")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		panic("failed to connect database")
// 	}
// 	defer db.Close()

// 	// Migrate the schema
// 	db.AutoMigrate(&Case{})
// }

// func setAuth0Variables() {
// 	// audience = os.Getenv("AUTH0_API_IDENTIFIER")
// 	// domain = os.Getenv("AUTH0_DOMAIN")
// 	audience = "https://go-angular-api"
// 	domain = "dev-hyiyx4g0.auth0.com"
// }

// //CORSMiddleware Cross-Origin Resource Sharing helper Class
// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, GET, OPTIONS, POST, PUT")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	}
// }

// // ValidateRequest will verify that a token received from an http request
// // is valid and signyed by Auth0
// func authRequired() gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		var auth0Domain = "https://" + domain + "/"
// 		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: auth0Domain + ".well-known/jwks.json"}, nil)
// 		configuration := auth0.NewConfiguration(client, []string{audience}, auth0Domain, jose.RS256)
// 		validator := auth0.NewValidator(configuration, nil)

// 		_, err := validator.ValidateRequest(c.Request)

// 		if err != nil {
// 			log.Println(err)
// 			terminateWithError(http.StatusUnauthorized, "token is not valid", c)
// 			return
// 		}
// 		c.Next()
// 	}
// }

// func terminateWithError(statusCode int, message string, c *gin.Context) {
// 	c.JSON(statusCode, gin.H{"error": message})
// 	c.Abort()
// }

// var (
// 	list []Case
// 	mtx  sync.RWMutex
// 	once sync.Once
// )

// // func init() {
// // 	once.Do(initialiseList)
// // }

// // func initialiseList() {
// // 	list = []Todo{}
// // }

// //Case Framework Object
// type Case struct {
// 	CaseCode              string `json:"CaseCode"`
// 	CaseType              string `json:"CaseType"`
// 	CaseName              string `json:"CaseName"`
// 	CaseShortDesc         string `json:"CaseShortDesc"`
// 	CaseNotes             string `json:"CaseNotes"`
// 	CaseTypeSLA           string `json:"CaseTypeSLA"`
// 	CoordinatorAssignedTo string `json:"CoordinatorAssignedTo"`
// 	CaseSeverity          string `json:"CaseSeverity"`
// 	CasePriority          string `json:"CasePriority"`
// 	CaseStatus            string `json:"CaseStatus"`
// 	CaseCompleted         bool   `json:"CaseCompleted"`
// 	CaseCreatedBy         string `json:"CaseCreatedBy"`
// 	CaseUpdateBy          string `json:"CaseUpdateBy"`
// 	CaseTermedBy          string `json:"CaseTermedBy"`
// 	gorm.Model
// }

// func convertHTTPBodyToCase(httpBody io.ReadCloser) (Case, int, error) {
// 	body, err := ioutil.ReadAll(httpBody)
// 	if err != nil {
// 		return Case{}, http.StatusInternalServerError, err
// 	}
// 	defer httpBody.Close()
// 	return convertJSONBodyToCase(body)
// }

// func convertJSONBodyToCase(jsonBody []byte) (Case, int, error) {
// 	var cs Case
// 	err := json.Unmarshal(jsonBody, &cs)
// 	if err != nil {
// 		return Case{}, http.StatusBadRequest, err
// 	}
// 	return cs, http.StatusOK, nil
// }

// // GetCaseListHandler returns all current todo items
// func GetCaseListHandler(c *gin.Context) {

// 	pgStr := c.Param("page")

// 	// width := "42"
// 	u64, err := strconv.ParseInt(pgStr, 10, 32)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	page := int(u64)

// 	caseType := c.Param("caseType")

// 	c.JSON(http.StatusOK, allCases(caseType, page))
// }

// //GetCaseDetailsHandler returns single specific Provider Details
// func GetCaseDetailsHandler(c *gin.Context) {
// 	caseCode := c.Param("caseCode")
// 	c.JSON(http.StatusOK, getCase(caseCode))
// }

// // AddCaseHandler adds a new todo to the todo list
// func AddCaseHandler(c *gin.Context) {
// 	cs, statusCode, err := convertHTTPBodyToCase(c.Request.Body)
// 	if err != nil {
// 		c.JSON(statusCode, err)
// 		return
// 	}
// 	//(caseType string, name string, shortDesc string, notes string, dueDate string, assignedTo string, severity string, priority string, completed bool)
// 	c.JSON(statusCode, gin.H{"CaseCode": newCase(cs.CaseType, cs.CaseName, cs.CaseShortDesc, cs.CaseNotes, cs.CaseTypeSLA, cs.CoordinatorAssignedTo, cs.CaseSeverity, cs.CasePriority, cs.CaseCompleted, cs.CaseCreatedBy, cs.CaseStatus)})
// }

// // TermCaseHandler will delete a specified todo based on user http input
// func TermCaseHandler(c *gin.Context) {
// 	caseID := c.Param("caseCode")
// 	if err := termCase(caseID); err != nil {
// 		c.JSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	c.JSON(http.StatusOK, "")
// }

// // UpdateCaseHandler will complete a specified todo based on user http input
// func UpdateCaseHandler(c *gin.Context) {
// 	caseItem, statusCode, err := convertHTTPBodyToCase(c.Request.Body)
// 	if err != nil {
// 		c.JSON(statusCode, err)
// 		return
// 	}
// 	if updateCase(caseItem) != nil {
// 		c.JSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	c.JSON(http.StatusOK, "")
// }

// func allCases(caseType string, page int) []Case {
// 	db, err := gorm.Open("mssql", "sqlserver://PwrQ_DB_adm:f-azLbi4@192.168.1.166:1433?database=PowerQueueDB-2018-10-27-11-29")
// 	if err != nil {
// 		panic("failed to connect database")
// 	}
// 	defer db.Close()

// 	var cases []Case
// 	// db.Where("CaseType = ?", caseType).Find(&cases)

// 	switch codeType := caseType; codeType {
// 	case "all":
// 		db.Table("cases").Select("case_code, case_type, case_name, case_type_sla, coordinator_assigned_to, created_at").Group("case_code, case_type, case_name, case_type_sla, coordinator_assigned_to, created_at").Order("created_at asc").Offset((page - 1) * 10).Limit(10).Scan(&cases)
// 	default:
// 		db.Table("cases").Select("case_code, case_type, case_name, case_type_sla, coordinator_assigned_to, created_at").Group("case_code, case_type, case_name, case_type_sla, coordinator_assigned_to, created_at").Order("created_at asc").Where("case_type = ?", caseType).Offset((page - 1) * 10).Limit(10).Scan(&cases)
// 	}

// 	fmt.Println("{}", cases)

// 	// json.NewEncoder(w).Encode(tasks)
// 	return cases
// }

// func getCase(caseCode string) []Case {
// 	db, dbErr := gorm.Open("mssql", "sqlserver://PwrQ_DB_adm:f-azLbi4@192.168.1.166:1433?database=PowerQueueDB-2018-10-27-11-29")
// 	if dbErr != nil {
// 		panic("failed to connect database")
// 	}
// 	defer db.Close()

// 	var cs []Case
// 	db.Where("case_code = ?", caseCode).Find(&cs)

// 	return cs
// }

// func newCaseParse(caseCode string, caseType string, name string, shortDesc string, notes string, sla string, assignedTo string, severity string, priority string, completed bool, createdBy string, status string) Case {
// 	// currentTime := time.Now()
// 	return Case{
// 		// ID:       xid.New().String(),
// 		CaseCode:              caseCode,
// 		CaseType:              caseType,
// 		CaseName:              name,
// 		CaseShortDesc:         shortDesc,
// 		CaseNotes:             notes,
// 		CaseTypeSLA:           sla,
// 		CoordinatorAssignedTo: assignedTo,
// 		CaseSeverity:          severity,
// 		CasePriority:          priority,
// 		CaseStatus:            status,
// 		CaseCompleted:         completed,
// 		CaseCreatedBy:         createdBy,
// 	}
// }

// func createCaseCode(caseType string) string {
// 	switch codeType := caseType; codeType {
// 	case "Appeals":
// 		return "APP-" + xid.New().String()
// 	case "ProviderDispute":
// 		return "PDISP-" + xid.New().String()
// 	case "Credentialing":
// 		return "PCRED-" + xid.New().String()
// 	case "Contracting":
// 		return "PCONT-" + xid.New().String()
// 	default:
// 		return "MGRIV-" + xid.New().String()
// 	}
// }

// func newCase(caseType string, name string, shortDesc string, notes string, sla string, assignedTo string, severity string, priority string, completed bool, createdBy string, status string) string {
// 	// fmt.Println("New User Endpoint Hit")

// 	db, err := gorm.Open("mssql", "sqlserver://PwrQ_DB_adm:f-azLbi4@192.168.1.166:1433?database=PowerQueueDB-2018-10-27-11-29")
// 	if err != nil {
// 		panic("failed to connect database")
// 	}
// 	defer db.Close()

// 	c := newCaseParse(createCaseCode(caseType), caseType, name, shortDesc, notes, sla, assignedTo, severity, priority, completed, "PQRS", status)

// 	db.Create(&Case{CaseCode: c.CaseCode, CaseType: c.CaseType, CaseName: c.CaseName, CaseShortDesc: c.CaseShortDesc, CaseNotes: c.CaseNotes, CaseTypeSLA: c.CaseTypeSLA, CoordinatorAssignedTo: c.CoordinatorAssignedTo, CaseSeverity: c.CaseSeverity, CasePriority: c.CasePriority, CaseCompleted: c.CaseCompleted, CaseCreatedBy: c.CaseCreatedBy, CaseStatus: c.CaseStatus})

// 	return c.CaseCode
// }

// func termCase(caseID string) error {
// 	db, dbErr := gorm.Open("mssql", "sqlserver://PwrQ_DB_adm:f-azLbi4@192.168.1.166:1433?database=PowerQueueDB-2018-10-27-11-29")
// 	if dbErr != nil {
// 		panic("failed to connect database")
// 	}
// 	defer db.Close()

// 	var cs Case

// 	err := db.Where("id = ?", caseID).Find(&cs)
// 	if err != nil {
// 		return err.Error
// 	}
// 	db.Delete(&cs)
// 	return nil
// }

// func updateCase(cs Case) error {
// 	db, dbErr := gorm.Open("mssql", "sqlserver://PwrQ_DB_adm:f-azLbi4@192.168.1.166:1433?database=PowerQueueDB-2018-10-27-11-29")
// 	if dbErr != nil {
// 		panic("failed to connect database")
// 	}
// 	defer db.Close()

// 	var updateCase Case
// 	err := db.Where("task_name = ?", cs.CaseName).Find(&updateCase)
// 	if err != nil {
// 		return err.Error
// 	}

// 	updateCase.CoordinatorAssignedTo = cs.CoordinatorAssignedTo

// 	err = db.Save(&cs)
// 	if err != nil {
// 		return err.Error
// 	}

// 	return nil
// }

// //UserLogin Framework Object
// type UserLogin struct {
// 	rfidTag            string `json:"rfidTag"`
// 	UserName           string `json:"UserName"`
// 	UserLoggedIn       bool   `json:"UserLoggedIn"`
// 	UserLoginCreatedBy string `json:"UserLoginCreatedBy"`
// 	gorm.Model
// }

// // LoginUserHandler adds a new todo to the todo list
// func LoginUserHandler(c *gin.Context) {
// 	ul, statusCode, err := convertHTTPBodyToUserLogin(c.Request.Body)
// 	if err != nil {
// 		c.JSON(statusCode, err)
// 		return
// 	}
// 	//(caseType string, name string, shortDesc string, notes string, dueDate string, assignedTo string, severity string, priority string, completed bool)
// 	c.JSON(statusCode, gin.H{"CaseCode": newLogin(ul.rfidTag, ul.UserName, ul.UserLoggedIn)})
// }
