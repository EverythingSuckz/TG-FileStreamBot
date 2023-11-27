package routes

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Route struct {
	Name   string
	Engine *gin.Engine
}

func (r *Route) Init(engine *gin.Engine) {
	r.Engine = engine
}

type allRoutes struct {
	log *zap.Logger
}

func Load(log *zap.Logger, r *gin.Engine) {
	log = log.Named("routes")
	defer log.Sugar().Info("Loaded all API Routes")
	route := &Route{Name: "/", Engine: r}
	route.Init(r)
	Type := reflect.TypeOf(&allRoutes{log})
	Value := reflect.ValueOf(&allRoutes{log})
	for i := 0; i < Type.NumMethod(); i++ {
		Type.Method(i).Func.Call([]reflect.Value{Value, reflect.ValueOf(route)})
	}
}
