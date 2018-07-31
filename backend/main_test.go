package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lon9/discord-generalized-sound-bot/backend/config"
	"github.com/lon9/discord-generalized-sound-bot/backend/database"
	"github.com/lon9/discord-generalized-sound-bot/backend/models"
	"github.com/lon9/discord-generalized-sound-bot/backend/server"
	testfixtures "gopkg.in/testfixtures.v2"
)

var (
	r        *gin.Engine
	fixtures *testfixtures.Context
)

func TestHealthStatus(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code should be 200:%d", w.Code)
	}

	if w.Body.String() != "Working!" {
		t.Errorf("Body should be Working!:%s", w.Body.String())
	}
}

func TestCategoriesIndex(t *testing.T) {
	prepareTestDatabase()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code should be 200:%d", w.Code)
	}

	var response struct {
		Status  int               `json:"status"`
		Message string            `json:"message"`
		Result  models.Categories `json:"result"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.Status != http.StatusOK {
		t.Errorf("Response.Status should be 200:%d", response.Status)
	}

	if len(response.Result) != 2 {
		t.Errorf("Categories length should be 2:%d", len(response.Result))
	}

	if response.Result[0].ID != 1 {
		t.Errorf("Categories[0].ID should be 1:%d", response.Result[0].ID)
	}

	if response.Result[0].Name != "Test Category 1" {
		t.Errorf("Categories[0].Name should be Test Category 1:%s", response.Result[0].Name)
	}

	if response.Result[1].ID != 2 {
		t.Errorf("Categories[1].ID should be 2:%d", response.Result[1].ID)
	}

	if response.Result[1].Name != "Test Category 2" {
		t.Errorf("Categories[1].Name should be Test Category 2:%s", response.Result[1].Name)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/categories/?query=Category%201", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code should be 200:%d", w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.Status != http.StatusOK {
		t.Errorf("Response.Status should be 200:%d", response.Status)
	}

	if len(response.Result) != 1 {
		t.Errorf("Response.Result length should be 1:%d", len(response.Result))
	}

	if response.Result[0].ID != 1 {
		t.Errorf("Response.Result[0].ID should be 1:%d", response.Result[0].ID)
	}

	if response.Result[0].Name != "Test Category 1" {
		t.Errorf("Response.Result[0].Name should be Test Category 1%s", response.Result[0].Name)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/categories/?query=UnknownQuery", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code should be 200:%d", w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.Status != http.StatusOK {
		t.Errorf("Response.Status should be 200:%d", response.Status)
	}

	if len(response.Result) != 0 {
		t.Errorf("Response.Result length should be 0:%d", len(response.Result))
	}
}

func TestCategoriesShow(t *testing.T) {
	prepareTestDatabase()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code should be 200:%d", w.Code)
	}

	var response struct {
		Status  int             `json:"status"`
		Message string          `json:"message"`
		Result  models.Category `json:"result"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.Status != http.StatusOK {
		t.Errorf("Response.Status should be 200:%d", response.Status)
	}

	if response.Result.ID != 1 {
		t.Errorf("Category.ID should be 1:%d", response.Result.ID)
	}

	if response.Result.Name != "Test Category 1" {
		t.Errorf("Category.Name should be Test Category 1:%s", response.Result.Name)
	}

	if len(response.Result.Sounds) != 2 {
		t.Errorf("Category.Category.Sounds length should be 2%d", len(response.Result.Sounds))
	}

	if response.Result.Sounds[0].ID != 1 {
		t.Errorf("Category.Sounds[0].ID sould be 1:%d", response.Result.Sounds[0].ID)
	}

	if response.Result.Sounds[0].Name != "Test Sound 1" {
		t.Errorf("Category.Sounds[0].Name sould be Test Sound 1:%s", response.Result.Sounds[0].Name)
	}

	if response.Result.Sounds[0].Path != "data/test/sound1.dca" {
		t.Errorf("Category.Sounds[0].Path sould be data/test/sound1.dca:%s", response.Result.Sounds[0].Path)
	}

	if response.Result.Sounds[0].CategoryID != 1 {
		t.Errorf("Category.Sounds[0].CategoryID should be 1:%d", response.Result.Sounds[0].CategoryID)
	}

	if response.Result.Sounds[1].ID != 2 {
		t.Errorf("Category.Sounds[1].ID sould be 2:%d", response.Result.Sounds[0].ID)
	}

	if response.Result.Sounds[1].Name != "Test Sound 2" {
		t.Errorf("Category.Sounds[1].Name sould be Test Sound 2:%s", response.Result.Sounds[0].Name)
	}

	if response.Result.Sounds[1].Path != "data/test/sound2.dca" {
		t.Errorf("Category.Sounds[1].Path sould be data/test/sound2.dca:%s", response.Result.Sounds[0].Path)
	}

	if response.Result.Sounds[1].CategoryID != 1 {
		t.Errorf("Category.Sounds[1].CategoryID should be 1:%d", response.Result.Sounds[0].CategoryID)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/categories/100", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Status code should be 404:%d", w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.Status != http.StatusNotFound {
		t.Errorf("Response.Status should be 404:%d", response.Status)
	}

	if response.Message == "" {
		t.Errorf("Response.Message should not to be empty:%s", response.Message)
	}
}

func TestSoundsIndex(t *testing.T) {
	prepareTestDatabase()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/sounds/?query=Sound%201", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code should be 200:%d", w.Code)
	}

	var response struct {
		Status  int           `json:"status"`
		Message string        `json:"message"`
		Result  models.Sounds `json:"result"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.Status != http.StatusOK {
		t.Errorf("Response.Status should be 200:%d", response.Status)
	}

	if len(response.Result) != 1 {
		t.Errorf("Response.Result length should be 1:%d", len(response.Result))
	}

	if response.Result[0].ID != 1 {
		t.Errorf("Response.Result[0].ID should be 1:%d", response.Result[0].ID)
	}

	if response.Result[0].Name != "Test Sound 1" {
		t.Errorf("Response.Result[0].Name should be Test Sound 1%s", response.Result[0].Name)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/sounds/?query=UnknownQuery", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code should be 200:%d", w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.Status != http.StatusOK {
		t.Errorf("Response.Status should be 200:%d", response.Status)
	}

	if len(response.Result) != 0 {
		t.Errorf("Response.Result length should be 0:%d", len(response.Result))
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/sounds/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code should be 400:%d", w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.Status != http.StatusBadRequest {
		t.Errorf("Response.Status should be 400:%d", response.Status)
	}
}

func readSoundFile(p, categoryName, uri string) (*http.Request, error) {
	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(p))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}

	writer.WriteField("name", filepath.Base(p[:len(p)-len(filepath.Ext(p))]))
	writer.WriteField("categoryName", categoryName)
	if err = writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func login() (token string, err error) {
	c := config.GetConfig()
	authInfo := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		c.GetString("auth.username"),
		"password",
	}
	jsonByte, err := json.Marshal(authInfo)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonByte))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response struct {
		Token string `json:"token"`
	}
	if err = json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		return
	}
	token = response.Token
	return
}

func TestSoundsCreate(t *testing.T) {

	token, err := login()
	if err != nil {
		t.Error(err)
	}

	// Test for wav and exist category
	prepareTestDatabase()

	req, err := readSoundFile("fixtures/files/sample1.wav", "Test Category 1", "/admin/sounds")
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Status code should be 201:%d", w.Code)
	}

	var response struct {
		Status  int           `json:"status"`
		Message string        `json:"message"`
		Result  *models.Sound `json:"result"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.Status != http.StatusCreated {
		t.Errorf("Response.Status should be 201:%d", response.Status)
	}

	if response.Result == nil {
		t.Error("Response.Result should not to be nil")
	}

	if response.Result.ID == 3 {
		t.Errorf("Response.Result.ID should not to be 0:%d", response.Result.ID)
	}

	if response.Result.Name != "sample1" {
		t.Errorf("Response.Result.Name should be sample1:%s", response.Result.Name)
	}

	if response.Result.Path != "sounds_dca/Test Category 1/sample1.dca" {
		t.Errorf("Response.Result.Path should be sounds_dca/Test Category 1/sample1.dca:%s", response.Result.Path)
	}

	if response.Result.CategoryID == 0 {
		t.Errorf("Response.Result.CategoryID should not to be 0:%d", response.Result.CategoryID)
	}

	if response.Result.Category == nil {
		t.Error("Response.Result.Category should not to be nil")
	}

	if response.Result.Category.ID == 0 {
		t.Errorf("Response.Result.Category.ID should not to be 0:%d", response.Result.Category.ID)
	}

	if response.Result.Category.Name != "Test Category 1" {
		t.Errorf("Response.Result.Category.Name should be Test Category 1:%s", response.Result.Category.Name)
	}

	if _, err := os.Stat("data/test/Test Category 1/sample1.dca"); err != nil {
		t.Error("File should be exist on data/test/Test Category 1/sample1.dca")
	}

	// Test for mp3 and non-exist-category
	req, err = readSoundFile("fixtures/files/sample2.mp3", "newcategory", "/admin/sounds")
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Status code should be 201:%d", w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.Status != http.StatusCreated {
		t.Errorf("Response.Status should be 201:%d", response.Status)
	}

	if response.Result == nil {
		t.Error("Response.Result should not to be nil")
	}

	if response.Result.ID == 0 {
		t.Errorf("Response.Result.ID should not to be 0:%d", response.Result.ID)
	}

	if response.Result.Name != "sample2" {
		t.Errorf("Response.Result.Name should be sample2:%s", response.Result.Name)
	}

	if response.Result.Path != "sounds_dca/newcategory/sample2.dca" {
		t.Errorf("Response.Result.Path should be sounds_dca/newcategory/sample2.dca:%s", response.Result.Path)
	}

	if response.Result.CategoryID == 0 {
		t.Errorf("Response.Result.CategoryID should not to be 0:%d", response.Result.CategoryID)
	}

	if response.Result.Category == nil {
		t.Error("Response.Result.Category should not to be nil")
	}

	if response.Result.Category.ID == 0 {
		t.Errorf("Response.Result.Category.ID should not to be 0%d", response.Result.Category.ID)
	}

	if response.Result.Category.Name != "newcategory" {
		t.Errorf("Response.Result.Category.Name should be newcategory:%s", response.Result.Category.Name)
	}

	if _, err := os.Stat("data/test/newcategory/sample2.dca"); err != nil {
		t.Error("File should be exist on data/test/newcategory/sample2.dca")
	}

	// Test for flac
	req, err = readSoundFile("fixtures/files/sample3.flac", "newcategory", "/admin/sounds")
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Status code should be 201:%d", w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.Status != http.StatusCreated {
		t.Errorf("Response.Status should be 201:%d", response.Status)
	}

	if response.Result == nil {
		t.Error("Response.Result should not to be nil")
	}

	if response.Result.ID == 0 {
		t.Errorf("Response.Result.ID should not to be 0:%d", response.Result.ID)
	}

	if response.Result.Name != "sample3" {
		t.Errorf("Response.Result.Name should be sample3:%s", response.Result.Name)
	}

	if response.Result.Path != "sounds_dca/newcategory/sample3.dca" {
		t.Errorf("Response.Result.Path should be sounds_dca/newcategory/sample3.dca:%s", response.Result.Path)
	}

	if response.Result.CategoryID == 0 {
		t.Errorf("Response.Result.CategoryID should not to be 0:%d", response.Result.CategoryID)
	}

	if response.Result.Category == nil {
		t.Error("Response.Result.Category should not to be nil")
	}

	if response.Result.Category.ID == 0 {
		t.Errorf("Response.Result.Category.ID should not to be 0%d", response.Result.Category.ID)
	}

	if response.Result.Category.Name != "newcategory" {
		t.Errorf("Response.Result.Category.Name should be newcategory:%s", response.Result.Category.Name)
	}

	if _, err := os.Stat("data/test/newcategory/sample3.dca"); err != nil {
		t.Error("File should be exist on data/test/newcategory/sample3.dca")
	}

	// Test for ogg
	req, err = readSoundFile("fixtures/files/sample4.ogg", "newcategory", "/admin/sounds")
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Status code should be 201:%d", w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.Status != http.StatusCreated {
		t.Errorf("Response.Status should be 201:%d", response.Status)
	}

	if response.Result == nil {
		t.Error("Response.Result should not to be nil")
	}

	if response.Result.ID == 0 {
		t.Errorf("Response.Result.ID should not to be 0:%d", response.Result.ID)
	}

	if response.Result.Name != "sample4" {
		t.Errorf("Response.Result.Name should be sample4:%s", response.Result.Name)
	}

	if response.Result.Path != "sounds_dca/newcategory/sample4.dca" {
		t.Errorf("Response.Result.Path should be sounds_dca/newcategory/sample4.dca:%s", response.Result.Path)
	}

	if response.Result.CategoryID == 0 {
		t.Errorf("Response.Result.CategoryID should not to be 0:%d", response.Result.CategoryID)
	}

	if response.Result.Category == nil {
		t.Error("Response.Result.Category should not to be nil")
	}

	if response.Result.Category.ID == 0 {
		t.Errorf("Response.Result.Category.ID should not to be 0%d", response.Result.Category.ID)
	}

	if response.Result.Category.Name != "newcategory" {
		t.Errorf("Response.Result.Category.Name should be newcategory:%s", response.Result.Category.Name)
	}

	if _, err := os.Stat("data/test/newcategory/sample4.dca"); err != nil {
		t.Error("File should be exist on data/test/newcategory/sample4.dca")
	}

	// Duplicate file
	req, err = readSoundFile("fixtures/files/sample1.wav", "newcategory", "/admin/sounds")
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code should be 400:%d", w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.Status != http.StatusBadRequest {
		t.Errorf("Response.Status should be 400:%d", response.Status)
	}

	if response.Message == "" {
		t.Error("Response.Message should not empty")
	}

	// Invalid mime type
	req, err = readSoundFile("fixtures/files/sample5.txt", "newcategory", "/admin/sounds")
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code should be 400:%d", w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.Status != http.StatusBadRequest {
		t.Errorf("Response.Status should be 400:%d", response.Status)
	}

	if response.Message == "" {
		t.Error("Response.Message should not empty")
	}

	db := database.GetDB()
	if !db.Where("name = ?", "sample5").Find(new(models.Sound)).RecordNotFound() {
		t.Error("should be not found")
	}

}

func TestMain(m *testing.M) {
	config.Init("test")
	database.Init(true, &models.Sound{}, &models.Category{})
	r = server.NewRouter()
	dbURL := config.GetConfig().GetString("db.url")
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		panic(err)
	}
	fixtures, err = testfixtures.NewFolder(db, &testfixtures.SQLite{}, "fixtures")
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		panic(err)
	}
}
