package models

// CategoriesAndSounds is struct of Categories and Sounds
type CategoriesAndSounds struct {
	Categories Categories
	Sounds     Sounds
}

// SearchAllByName searches sound and category by name
func (cas *CategoriesAndSounds) SearchAllByName(name string) (err error) {
	if err = cas.Sounds.SearchByName(name); err != nil {
		return
	}
	if err = cas.Categories.SearchByName(name); err != nil {
		return
	}
	return
}
