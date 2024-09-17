package swaggerui

import (
	"fmt"
	httpx "net/http"

	"github.com/go-kratos/kratos/v2/transport/http"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SwaggerUI(srv *http.Server, staticDir string, serviceName string, version string) {
	fmt.Printf("Absolute path: %s\n", staticDir)
	srv.HandlePrefix("/swagger-ui/", httpx.StripPrefix("/swagger-ui/", httpx.FileServer(httpx.Dir(staticDir))))
	uri := fmt.Sprintf("/swagger-ui/%s/%s.swagger.json", version, serviceName)
	srv.HandlePrefix("/swagger/", httpSwagger.Handler(httpSwagger.URL(uri)))
}
