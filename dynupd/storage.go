package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/virtcanhead/health"
	"sort"
)

type Storage interface {
	health.Resource
	GetProjectNames() ([]string, error)
	GetProject(name string) (Project, error)
	CreateProject(name string) error
	DestroyProject(name string) error
	UpdateProject(p Project) error
}

func checkRedisError(err *error) bool {
	if err == nil {
		return false
	}
	if *err == nil {
		return false
	}
	if *err == redis.Nil {
		*err = nil
		return false
	}
	return true
}

type redisStorage struct {
	client *redis.Client
}

func (s *redisStorage) getJSON(key string, out interface{}) (err error) {
	var val string
	if val, err = s.client.Get(key).Result(); checkRedisError(&err) {
		return
	}
	if len(val) == 0 {
		return
	}
	return json.Unmarshal([]byte(val), out)
}

func (s *redisStorage) setJSON(key string, in interface{}) (err error) {
	var val []byte
	if val, err = json.Marshal(in); err != nil {
		return
	}
	if err = s.client.Set(key, string(val), 0).Err(); checkRedisError(&err) {
		return
	}
	return
}

func (s *redisStorage) HealthCheck() error {
	return s.client.Ping().Err()
}

func (s *redisStorage) GetProjectNames() (names []string, err error) {
	if names, err = s.client.SMembers("dynup.projects").Result(); checkRedisError(&err) {
		return
	}
	if names == nil {
		names = make([]string, 0, 0)
	} else {
		sort.StringSlice(names).Sort()
	}
	return
}

func (s *redisStorage) GetProject(name string) (p Project, err error) {
	p.Name = name
	if err = s.getJSON("gateway-rules-"+name, &p.Rules); err != nil {
		return
	}
	if p.Rules == nil {
		p.Rules = make([]Rule, 0, 0)
	}
	if err = s.getJSON("gateway-backends-"+name, &p.Backends); err != nil {
		return
	}
	if p.Backends == nil {
		p.Backends = make(map[string][]string, 0)
	}
	return
}

func (s *redisStorage) CreateProject(name string) (err error) {
	if err = s.client.SAdd("dynup.projects", name).Err(); checkRedisError(&err) {
		return
	}
	return
}

func (s *redisStorage) DestroyProject(name string) (err error) {
	if err = s.client.SRem("dynup.projects", name).Err(); checkRedisError(&err) {
		return
	}
	return
}

func (s *redisStorage) UpdateProject(p Project) (err error) {
	if err = s.setJSON("gateway-rules-"+p.Name, p.Rules); err != nil {
		return
	}
	if err = s.setJSON("gateway-backends-"+p.Name, p.Backends); err != nil {
		return
	}
	return
}

func NewRedisStorage(addr string) Storage {
	return &redisStorage{client: redis.NewClient(&redis.Options{Addr: addr})}
}
