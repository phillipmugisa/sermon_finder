package server

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

type AppError struct {
	msg  string
	code int
}

func NewError(msg string, code int) *AppError {
	return &AppError{
		msg:  msg,
		code: code,
	}
}

func (err *AppError) Error() string {
	return err.msg
}

type requestHandler func(context.Context, http.ResponseWriter, *http.Request) error

func (server *Server) MakeRequestHandler(handler requestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := context.Background()
		access_err := server.monitor.Monitor(r)
		if access_err != nil {
			server.HandlerRequestError(access_err)
			return
		}

		if err := handler(ctx, w, r); err != nil {
			server.HandlerRequestError(err)
			return
		}
	}
}

func (server *Server) HandlerRequestError(err error) {
	server.logger.Error(fmt.Sprintf("%s", err))
}

func (server Server) Render(c context.Context, w io.Writer, name string, data interface{}) error {
	return server.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *template.Template {
	return template.Must(template.ParseGlob("templates/*.html"))
}
