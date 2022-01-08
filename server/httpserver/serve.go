package httpserver

import (
	"context"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strconv"
	"time"

	"github.com/0w0mewo/budong-apigateway/server"
	"github.com/0w0mewo/budong-apigateway/server/grpcclient"
	"github.com/0w0mewo/budong-apigateway/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ApiServer struct {
	routes     *gin.Engine
	httpServer *http.Server
	service    server.Service
	logger     *logrus.Entry
	routemap   map[string]gin.HandlerFunc
}

func NewRestfulServer(addr string) *ApiServer {
	routes := gin.New()
	logger := logrus.StandardLogger().WithField("server", "rest")
	serve := grpcclient.NewSetuGrpcClient("127.0.0.1:9999")

	return &ApiServer{
		routes: routes,
		httpServer: &http.Server{
			Addr:    addr,
			Handler: routes,
		},
		logger:  logger,
		service: serve,
	}
}

func (r *ApiServer) Init() {
	r.routemap = map[string]gin.HandlerFunc{
		"random":                     r.hello,
		"dofetch/:num":               r.fetchsetu,
		"inventory/:page/:page_size": r.inventory,
		"get/:id":                    r.givemesetu,
		"/":                          r.help,
		"/help":                      r.help,
		"count":                      r.count,
	}
}

func (r *ApiServer) Run() {
	for path, handler := range r.routemap {
		r.routes.GET(path, handler)
	}

	r.httpServer.ListenAndServe()
}

func (r *ApiServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r.service.Shutdown()
	r.httpServer.Shutdown(ctx)

	r.logger.Infoln("http server shutdown")
}

func (r *ApiServer) givemesetu(c *gin.Context) {
	tryOrSendErr(c, func() error {
		// get setu id
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}

		// search it in db
		se, err := r.service.GetSetuFromDB(id)
		if err != nil {
			return err
		}

		imgType := utils.ImageBytesFmt(se)

		c.Data(http.StatusOK, "image/"+imgType, se)
		return nil
	})

}

func (r *ApiServer) help(c *gin.Context) {
	fmt.Fprintf(c.Writer, `
		GET /get/<id> : get a setu from DB by id, id is int
		GET /dofetch/<amount> : fetch amount of setu and store them to DB, amount is int
		GET /random : return a random setu from DB
		GET /inventory/<page>/<page_size>
		GET /count
	
	
	`)
}

func (r *ApiServer) count(c *gin.Context) {
	tryOrSendErr(c, func() error {
		cnt := r.service.Count()
		c.JSON(http.StatusOK, &CountResp{Count: cnt})

		return nil

	})
}

func (r *ApiServer) inventory(c *gin.Context) {
	tryOrSendErr(c, func() error {
		// get page param
		page, err := strconv.ParseUint(c.Param("page"), 10, 64)
		if err != nil {
			return err
		}

		// get page size param
		size, err := strconv.ParseUint(c.Param("page_size"), 10, 64)
		if err != nil {
			return err
		}

		// ensure page size is between 0 and 100
		if size > 50 || size < 1 {
			return ErrPageSize
		}

		// ensure page is in valid range
		if page < 1 || page > r.service.Count()/size+1 {
			return ErrPageRange
		}

		data, err := r.service.GetInventory(page, size)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, data)

		return nil
	})

}

func (r *ApiServer) hello(c *gin.Context) {
	tryOrSendErr(c, func() error {
		img, err := r.service.RandomSetu()
		if err != nil {
			return err
		}

		imgType := utils.ImageBytesFmt(img)
		if err != nil {
			return err

		}

		c.Data(http.StatusOK, "image/"+imgType, img)

		return nil
	})

}

func (r *ApiServer) fetchsetu(c *gin.Context) {
	tryOrSendErr(c, func() error {
		num, err := strconv.Atoi(c.Param("num"))
		if err != nil {
			return err
		}

		err = r.service.RequestSetu(num, false) // 不可以涩色
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, &ErrResp{ErrMsg: "ok"})
		return nil
	})

}
