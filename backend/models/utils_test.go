package models

import "testing"

func TestSearchAllByName(t *testing.T) {
	var cas CategoriesAndSounds
	if err := cas.SearchAllByName("Test"); err != nil {
		t.Error(err)
	}

	if len(cas.Categories) != 2 {
		t.Errorf("CategoriesAndSounds.Categories length should be 2:%d", len(cas.Categories))
	}

	if len(cas.Sounds) != 2 {
		t.Errorf("CategoriesAndSounds.Sounds length should be 2:%d", len(cas.Sounds))
	}

	if err := cas.SearchAllByName("Non-exist name"); err != nil {
		t.Error(err)
	}

	if len(cas.Categories) != 0 || len(cas.Sounds) != 0 {
		t.Errorf("If it receive non-exist name, should return 0 size slice:%d:%d", len(cas.Categories), len(cas.Sounds))
	}
}
