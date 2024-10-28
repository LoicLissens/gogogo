package test_utils

import (
	"jiva-guildes/domain/models/guilde"
	"jiva-guildes/domain/ports/repositories"
	"time"
)

func CreateGuilde(repo repositories.GuildeRepository, opts guilde.GuildeOptions) guilde.Guilde {

	g, err := guilde.New(opts)
	if err != nil {
		panic(err)
	}
	savedGuilde, err := repo.Save(*g)
	if err != nil {
		panic(err)
	}
	return savedGuilde
}
func SaveBatchSamples(repo repositories.GuildeRepository, samples []guilde.GuildeOptions) []guilde.GuildeOptions {
	for _, sample := range samples {
		CreateGuilde(repo, sample)
	}
	return samples
}
func SaveBasicSamples(repo repositories.GuildeRepository) []guilde.GuildeOptions {
	active := true
	notActive := false
	today := time.Now().UTC()
	date := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	samples := []guilde.GuildeOptions{
		{Name: "GUnit", Page_url: "https://www.googleimage.com", Exists: true, Validated: true, Active: nil, Creation_date: nil},
		{Name: "D12", Page_url: "img1", Exists: true, Validated: true, Active: nil, Creation_date: nil},
		{Name: "eminem", Exists: true, Validated: true, Active: &active, Creation_date: &today},
		{Name: "AS$AP", Page_url: "img3", Exists: false, Validated: false, Active: &notActive, Creation_date: &date},
	}
	return SaveBatchSamples(repo, samples)
}
