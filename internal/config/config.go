package config

import (
	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
)

type DigList struct {
	List     []Host
}

type Host struct {
	Name string `json:"name"`
}

func Load(filename string) (DigList, error) {
	var digList DigList
	if err := configor.Load(&digList, filename); err != nil {
		return DigList{}, errors.WithMessage(err, "failed to load digList")
	}
	return digList, nil
}

