package Controllers

import (
	"bytes"
	"encoding/json"
	"go-todo-app/Config"

	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestCreateATodo(t *testing.T) {
	// Create a new Gin router instance
	db := Config.Database.ConnectToDB()

	defer db.Close()
	Config.NewTable()
	// Create a test request with sample encrypted data
	decryptedData := []byte(`{"Title": "Test Title", "Description": "Test Description"}`)
	req, err := http.NewRequest(http.MethodPost, "/todo", nil)
	req.Header.Set("x-key", "noenonrgkgneroiw")
	req.Header.Set("x-iv", "1461618689689168")
	//req.SetBasicAuth("username", "password")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Set up a test context with the encrypted data
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Set("decryptedText", decryptedData)
	ctx.Request = req

	// Call the CreateATodo function with the test context
	CreateATodo(ctx)

	// Assert that the response is a successful HTTP status and contains the expected message
	assert.Equal(t, http.StatusCreated, rr.Code)
	encrypted := AESEncrypt("Todo created Successfully.....", []byte(ctx.Request.Header.Get("x-key")), ctx.Request.Header.Get("x-iv"))
	actual := rr.Body.String()
	expected, err := json.Marshal(encrypted)
	assert.Equal(t, string(expected), actual)
}

func TestGetTodos(t *testing.T) {
	// create a new gin context
	req, err := http.NewRequest("GET", "/todos", nil)
	req.Header.Set("x-key", "noenonrgkgneroiw")
	req.Header.Set("x-iv", "1461618689689168")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req
	db := Config.Database.ConnectToDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// create a new todo record in the database
	_, err = db.Exec("INSERT INTO todo (title, description) VALUES (?, ?)", "Test Title", "Test Description")
	if err != nil {
		t.Fatal(err)
	}

	// call the GetTodos function
	GetTodos(ctx)

	// check the response
	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", rr.Code)
	}
	assert.NotNil(t, rr, "reponse is nil")
}

func TestGetATodo(t *testing.T) {
	// Create a new Gin router instance
	r := gin.Default()

	// Add a GET route for testing the GetATodo function
	r.GET("/todo/:id", GetATodo)
	db := Config.Database.ConnectToDB()
	defer db.Close()

	// Create a test request with a sample todo ID
	req, err := http.NewRequest(http.MethodGet, "/todo/1", nil)
	req.Header.Set("x-key", "noenonrgkgneroiw")
	req.Header.Set("x-iv", "1461618689689168")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Call the GetATodo function with the test context
	r.ServeHTTP(rr, req)

	// Assert that the response is a successful HTTP status and contains the expected todo data
	assert.Equal(t, http.StatusOK, rr.Code)
	expectedData := `{"id":1,"title":"Test Title","description":"Test Description"}`
	encrypted := AESEncrypt(expectedData, []byte("noenonrgkgneroiw"), "1461618689689168")
	//encrypted := AESEncrypt("Todo created Successfully.....", []byte(ctx.Request.Header.Get("x-key")), ctx.Request.Header.Get("x-iv"))
	actual := rr.Body.String()
	expected, err := json.Marshal(encrypted)
	fmt.Println("........................", string(expected))
	assert.Equal(t, string(expected), actual)

}

func TestUpdateATodo(t *testing.T) {
	// initialize Gin router
	db := Config.Database.ConnectToDB()
	defer db.Close()

	// create a mock HTTP request with a sample encrypted JSON payload
	req, err := http.NewRequest(http.MethodPut, "/todo/2", bytes.NewBufferString(`{"title":"Updated Title","description":"Updated Description"}`))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("x-key", "noenonrgkgneroiw")
	req.Header.Set("x-iv", "1461618689689168")
	// set up mock HTTP response recorder
	resp := httptest.NewRecorder()

	// set up mock Gin context
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "2"}}
	// simulate adding the decrypted data to the Gin context
	c.Set("decryptedText", []byte(`{"title":"Updated Title","description":"Updated Description"}`))

	// call the handler function
	UpdateATodo(c)

	// check the response status code
	if resp.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Code)
	}

	// check the response body
	encrypted := AESEncrypt("Updated Successfully.......", []byte(c.Request.Header.Get("x-key")), c.Request.Header.Get("x-iv"))
	expected, err := json.Marshal(encrypted)
	assert.Equal(t, string(expected), resp.Body.String())

}
func TestDeleteATodo(t *testing.T) {
	// set up test database
	db := Config.Database.ConnectToDB()
	defer db.Close()

	// insert a test record into the database
	_, err := db.Exec("insert into todo (Title, Description) values (?, ?)", "Test Title", "Test Description")
	if err != nil {
		t.Fatalf("failed to insert test record: %v", err)
	}

	// create a mock HTTP request
	req, err := http.NewRequest("DELETE", "/todo/3", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("x-key", "noenonrgkgneroiw")
	req.Header.Set("x-iv", "1461618689689168")
	// set up mock HTTP response recorder
	resp := httptest.NewRecorder()

	// set up mock Gin context
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "3"}}

	// call the handler function
	DeleteATodo(c)

	// check the response status code
	if resp.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Code)
	}

	// check the response body
	// expected := `{"data":"UmVjb3JkIGRlbGV0ZWQgU3VjY2Vzc2Z1bC4uLi4uLi4uLi4sLi4uLi4uLi4uLi4uLi4uLi4sLi4uLi4uLi4="}`
	// if resp.Body.String() != expected {
	// 	t.Errorf("expected body %v; got %v", expected, resp.Body.String())
	// }

	encrypted := AESEncrypt("Record deleted Succesfully.......", []byte(c.Request.Header.Get("x-key")), c.Request.Header.Get("x-iv"))
	expected, err := json.Marshal(encrypted)
	assert.Equal(t, string(expected), resp.Body.String())

	//check that the record was deleted from the database
	var count int
	err = db.QueryRow("select count(*) from todo where ID = ?", 3).Scan(&count)
	if err != nil {
		t.Fatalf("failed to query database: %v", err)
	}
	if count != 0 {
		t.Errorf("expected record to be deleted; got count=%v", count)
	}
}
