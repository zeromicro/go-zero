package rabbitmqgen

import (
	"fmt"
	"strings"

	apiSpec "github.com/lerity-yao/go-zero/tools/cztctl/api/spec"
	"github.com/zeromicro/go-zero/core/collection"
)

func generateRabbitmqEtcNames(api *apiSpec.ApiSpec) []string {
	names := collection.NewSet[string]()
	for _, g := range api.Service.Groups {
		for _, h := range g.Routes {
			l := fmt.Sprintf(
				"%sRabbitmqConf: \n  Username:\n  Password:\n  Host:\n  Port:\n  ListenerQueues:\n    - Name: %s\n",
				strings.TrimSuffix(h.Handler, "Handler"), strings.TrimPrefix(h.Path, "/"))
			names.Add(l)
		}
	}
	return names.Keys()
}

func generateRabbitmqConfigNames(api *apiSpec.ApiSpec) []string {
	names := collection.NewSet[string]()
	for _, g := range api.Service.Groups {
		for _, h := range g.Routes {
			l := fmt.Sprintf("%sRabbitmqConf rabbitmq.RabbitListenerConf", strings.TrimSuffix(h.Handler, "Handler"))
			names.Add(l)
		}
	}
	return names.Keys()
}
