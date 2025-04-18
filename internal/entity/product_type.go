package entity

import "errors"

type ProductType string

const (
	ShoesProductType       ProductType = "обувь"
	ClothesProductType     ProductType = "одежда"
	ElectronicsProductType ProductType = "электроника"
)

var validProductTypeMap = map[ProductType]struct{}{
	ShoesProductType:       {},
	ClothesProductType:     {},
	ElectronicsProductType: {},
}

func (r ProductType) IsValidProductType() bool {
	_, exists := validProductTypeMap[r]
	return exists
}

func (r ProductType) ValidateProductType() error {
	if !r.IsValidProductType() {
		return errors.New("invalid product type")
	}
	return nil
}
