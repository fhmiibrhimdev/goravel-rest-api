package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/mysql"
)

type AppServiceProvider struct {
}

func (receiver *AppServiceProvider) Register(app foundation.Application) {
    (&mysql.ServiceProvider{}).Register(app)
}

func (receiver *AppServiceProvider) Boot(app foundation.Application) {

}
