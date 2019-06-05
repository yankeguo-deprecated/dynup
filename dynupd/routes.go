package main

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
	"strings"
)

var (
	illegalProjectNameCharacters = regexp.MustCompile("[^a-zA-Z0-9_-]")
)

type Project struct {
	Name     string              `json:"name"`
	Rules    []Rule              `json:"rules"`
	Backends map[string][]string `json:"backends"`
}

type Rule struct {
	Type    string   `json:"type"`
	Fields  []string `json:"fields"`
	Pattern string   `json:"pattern"`
	Target  string   `json:"target"`
}

func mountRoutes(e *echo.Echo) {
	e.GET("/dynup/api/projects", routeProjects)
	e.GET("/dynup/api/projects/:name", routeProjectDetail)
	e.POST("/dynup/api/projects/create", routeProjectCreate)
	e.POST("/dynup/api/projects/:name/update", routeProjectUpdate)
	e.POST("/dynup/api/projects/:name/destroy", routeProjectDestroy)
}

func routeProjects(c echo.Context) (err error) {
	var ret []string
	if ret, err = storage.GetProjectNames(); err != nil {
		return
	}
	return c.JSON(http.StatusOK, ret)
}

func routeProjectCreate(c echo.Context) (err error) {
	p := Project{}
	if err = c.Bind(&p); err != nil {
		return
	}
	name := strings.TrimSpace(p.Name)
	if len(name) == 0 {
		return errors.New("project name cannot be empty")
	}
	name = illegalProjectNameCharacters.ReplaceAllString(name, "-")
	if err = storage.CreateProject(name); err != nil {
		return
	}
	return routeProjects(c)
}

func routeProjectDetail(c echo.Context) (err error) {
	name := c.Param("name")
	if len(name) == 0 {
		return errors.New("project name cannot be empty")
	}
	var p Project
	if p, err = storage.GetProject(name); err != nil {
		return
	}
	return c.JSON(http.StatusOK, p)
}

func routeProjectDestroy(c echo.Context) (err error) {
	name := strings.TrimSpace(c.Param("name"))
	if len(name) == 0 {
		return errors.New("project name cannot be empty")
	}
	if err = storage.DestroyProject(name); err != nil {
		return
	}
	return routeProjects(c)
}

func routeProjectUpdate(c echo.Context) (err error) {
	name := c.Param("name")
	if len(name) == 0 {
		return errors.New("project name cannot be empty")
	}
	p := Project{}
	if err = c.Bind(&p); err != nil {
		return
	}
	p.Name = name
	if err = storage.UpdateProject(p); err != nil {
		return
	}
	return routeProjectDetail(c)
}
