package models

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestCreateCategory(t *testing.T) {
	prepareTestDatabase()
	c := &Category{
		Name: "test",
	}
	if err := c.Create(); err != nil {
		t.Error(err)
	}

	if c.ID == 0 {
		t.Error("Category.ID should not to be 0")
	}

	if c.Name != "test" {
		t.Errorf("Category.Name should be test: %s", c.Name)
	}
}

func TestFindCategoryByID(t *testing.T) {
	prepareTestDatabase()
	var category Category
	if err := category.FindByID(1); err != nil {
		t.Error(err)
	}

	if category.ID != 1 {
		t.Errorf("Category.ID should be 1:%d", category.ID)
	}

	if category.Name != "Test Category 1" {
		t.Errorf("Category.Name should be Test Category 1:%s", category.Name)
	}

	if len(category.Sounds) != 2 {
		t.Errorf("Category.Category.Sounds length should be 2%d", len(category.Sounds))
	}

	if category.Sounds[0].ID != 1 {
		t.Errorf("Category.Sounds[0].ID sould be 1:%d", category.Sounds[0].ID)
	}

	if category.Sounds[0].Name != "Test Sound 1" {
		t.Errorf("Category.Sounds[0].Name sould be Test Sound 1:%s", category.Sounds[0].Name)
	}

	if category.Sounds[0].Path != "data/test/sound1.dca" {
		t.Errorf("Category.Sounds[0].Path sould be data/test/sound1.dca:%s", category.Sounds[0].Path)
	}

	if category.Sounds[0].CategoryID != 1 {
		t.Errorf("Category.Sounds[0].CategoryID should be 1:%d", category.Sounds[0].CategoryID)
	}

	if category.Sounds[1].ID != 2 {
		t.Errorf("Category.Sounds[1].ID sould be 2:%d", category.Sounds[0].ID)
	}

	if category.Sounds[1].Name != "Test Sound 2" {
		t.Errorf("Category.Sounds[1].Name sould be Test Sound 2:%s", category.Sounds[0].Name)
	}

	if category.Sounds[1].Path != "data/test/sound2.dca" {
		t.Errorf("Category.Sounds[1].Path sould be data/test/sound2.dca:%s", category.Sounds[0].Path)
	}

	if category.Sounds[1].CategoryID != 1 {
		t.Errorf("Category.Sounds[1].CategoryID should be 1:%d", category.Sounds[0].CategoryID)
	}

	if err := category.FindByID(100); err == nil {
		t.Error("If send non-exist id, should return error.")
	}

}

func TestFindCategoryByName(t *testing.T) {
	prepareTestDatabase()
	var category Category
	if err := category.FindByName("Test Category 1"); err != nil {
		t.Error(err)
	}

	if category.ID != 1 {
		t.Errorf("Category.ID should be 1:%d", category.ID)
	}

	if category.Name != "Test Category 1" {
		t.Errorf("Category.Name should be Test Category 1:%s", category.Name)
	}

	if len(category.Sounds) != 2 {
		t.Errorf("Category.Category.Sounds length should be 2%d", len(category.Sounds))
	}

	if category.Sounds[0].ID != 1 {
		t.Errorf("Category.Sounds[0].ID sould be 1:%d", category.Sounds[0].ID)
	}

	if category.Sounds[0].Name != "Test Sound 1" {
		t.Errorf("Category.Sounds[0].Name sould be Test Sound 1:%s", category.Sounds[0].Name)
	}

	if category.Sounds[0].Path != "data/test/sound1.dca" {
		t.Errorf("Category.Sounds[0].Path sould be data/test/sound1.dca:%s", category.Sounds[0].Path)
	}

	if category.Sounds[0].CategoryID != 1 {
		t.Errorf("Category.Sounds[0].CategoryID should be 1:%d", category.Sounds[0].CategoryID)
	}

	if category.Sounds[1].ID != 2 {
		t.Errorf("Category.Sounds[1].ID sould be 2:%d", category.Sounds[0].ID)
	}

	if category.Sounds[1].Name != "Test Sound 2" {
		t.Errorf("Category.Sounds[1].Name sould be Test Sound 2:%s", category.Sounds[0].Name)
	}

	if category.Sounds[1].Path != "data/test/sound2.dca" {
		t.Errorf("Category.Sounds[1].Path sould be data/test/sound2.dca:%s", category.Sounds[0].Path)
	}

	if category.Sounds[1].CategoryID != 1 {
		t.Errorf("Category.Sounds[1].CategoryID should be 1:%d", category.Sounds[0].CategoryID)
	}

	if err := category.FindByName("Non-exist category name"); err == nil {
		t.Error("If send non-exist category name, should return error.")
	}

}

func TestFindAllCategories(t *testing.T) {
	prepareTestDatabase()
	var categories Categories
	if err := categories.FindAll(); err != nil {
		t.Error(err)
	}

	if len(categories) != 2 {
		t.Errorf("Categories length sould be 2:%d", len(categories))
	}

	if categories[0].ID != 1 {
		t.Errorf("Categories[0].ID sould be 1:%d", categories[0].ID)
	}

	if categories[0].Name != "Test Category 1" {
		t.Errorf("Categories[0].Name sould be Test Category 1:%s", categories[0].Name)
	}

	if categories[1].ID != 2 {
		t.Errorf("Categories[1].ID sould be 2:%d", categories[1].ID)
	}

	if categories[1].Name != "Test Category 2" {
		t.Errorf("Categories[2].Name sould be Test Category 2:%s", categories[1].Name)
	}
}

func TestSearchCategoryByName(t *testing.T) {
	prepareTestDatabase()
	var categories Categories
	if err := categories.SearchByName("Test Category"); err != nil {
		t.Error(err)
	}

	if len(categories) != 2 {
		t.Errorf("Categories length sould be 2:%d", len(categories))
	}

	if categories[0].ID != 1 {
		t.Errorf("Categories[0].ID sould be 1:%d", categories[0].ID)
	}

	if categories[0].Name != "Test Category 1" {
		t.Errorf("Categories[0].Name sould be Test Category 1:%s", categories[0].Name)
	}

	if categories[1].ID != 2 {
		t.Errorf("Categories[1].ID sould be 2:%d", categories[1].ID)
	}

	if categories[1].Name != "Test Category 2" {
		t.Errorf("Categories[2].Name sould be Test Category 2:%s", categories[1].Name)
	}

	if err := categories.SearchByName("Non-exist category name"); err != nil {
		t.Error(err)
	}

	if len(categories) != 0 {
		t.Errorf("If it send non-exist category name, should receive 0 size slice:%d", len(categories))

	}
}
