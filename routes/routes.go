package routes

import (
	"log"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Name   string
	Engine *gin.Engine
}

func (r *Route) Init(engine *gin.Engine) {
	r.Engine = engine
}

type allRoutes struct {
}

func Load(r *gin.Engine) {
	defer log.Println("Loaded all API Routes")
	route := &Route{Name: "/", Engine: r}
	route.Init(r)
	Type := reflect.TypeOf(&allRoutes{})
	Value := reflect.ValueOf(&allRoutes{})
	for i := 0; i < Type.NumMethod(); i++ {
		Type.Method(i).Func.Call([]reflect.Value{Value, reflect.ValueOf(route)})
	}
}
