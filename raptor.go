package raptor

import (
	"log"

	"github.com/fatih/color"
	"github.com/valyala/fasthttp"
)

func init() {
	color.Blue(`
'########:::::'###::::'########::'########::'#######::'########::
##.... ##:::'## ##::: ##.... ##:... ##..::'##.... ##: ##.... ##:
##:::: ##::'##:. ##:: ##:::: ##:::: ##:::: ##:::: ##: ##:::: ##:
########::'##:::. ##: ########::::: ##:::: ##:::: ##: ########::
##.. ##::: #########: ##.....:::::: ##:::: ##:::: ##: ##.. ##:::
##::. ##:: ##.... ##: ##::::::::::: ##:::: ##:::: ##: ##::. ##::
##:::. ##: ##:::: ##: ##::::::::::: ##::::. #######:: ##:::. ##:
..:::::..::..:::::..::..::::::::::::..::::::.......:::..:::::..::
	
`)
}

// Raptor :
type Raptor struct {
	*Router
}

// HandlerFunc :
type HandlerFunc func(*Context) error

// MiddlewareFunc :
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// New :
func New() *Raptor {
	return &Raptor{
		Router: NewRouter(),
	}
}

// Start :
func (r *Raptor) Start(port string, handler ...HandlerFunc) {
	cb := r.router.Handler
	if len(handler) > 0 {
		cb = func(ctx *fasthttp.RequestCtx) {
			if err := handler[0](&Context{
				RequestCtx: ctx,
			}); err != nil {
				return
			}
		}
	}
	log.Println("fasthttp server started on", port)
	log.Fatal(fasthttp.ListenAndServe(port, cb))
}

// Close :
func (r *Raptor) Close() {
	return
}
