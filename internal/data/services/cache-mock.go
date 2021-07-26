package services

import "errors"

type CacheServicesMock struct {
}

func NewCacheServicesMock() *CacheServicesMock {
	return &CacheServicesMock{}
}

func (o *CacheServicesMock) Set(key string, value string) error {
	return nil
}

func (o *CacheServicesMock) Get(key string) (string, error) {
	var output string
	return output, errors.New("")
}

func (o *CacheServicesMock) Delete(key string) error {
	return nil
}
