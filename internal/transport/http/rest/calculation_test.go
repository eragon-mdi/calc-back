package resttransport

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eragon-mdi/calc-back/internal/domain"
	"github.com/eragon-mdi/calc-back/internal/transport/http/rest/mocks"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func newEchoCtx(method, uri string, args ...string) echo.Context {
	e := echo.New()

	var body io.Reader
	if len(args) > 0 {
		body = strings.NewReader(args[0])
	} else {
		body = nil
	}

	req := httptest.NewRequest(method, uri, body)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

func newLogger() *zap.SugaredLogger {
	return zap.NewNop().Sugar()
}

var logger = newLogger()

func Test_transport_GetLastCalculations(t *testing.T) {
	ctx := newEchoCtx(http.MethodGet, "/calculations")
	logger := newLogger()

	type fields struct {
		s func() Service
		l *zap.SugaredLogger
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			fields: fields{
				s: func() Service {
					ms := mocks.NewService(t)
					ms.EXPECT().GetLastCalculations().
						Return([]domain.Calculation{{ID: "1", Expression: "1+1", Result: "2"}}, nil)
					return ms
				},
				l: logger,
			},
			args: args{
				c: ctx,
			},
			wantErr: false,
		},
		{
			name: "failde service case: notFound",
			fields: fields{
				s: func() Service {
					ms := mocks.NewService(t)
					ms.EXPECT().GetLastCalculations().
						Return(nil, domain.ErrNotFound)
					return ms
				},
				l: logger,
			},
			args: args{
				c: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := transport{
				s: tt.fields.s(),
				l: tt.fields.l,
			}
			err := tr.GetLastCalculations(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("transport.GetLastCalculations() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				httpErr, ok := err.(*echo.HTTPError)
				if !ok {
					t.Errorf("error is not *echo.HTTPError: %v", err)
					return
				}

				if tt.name == "failde service case: notFound" {
					if httpErr.Code != http.StatusNotFound {
						t.Errorf("expected HTTP status %d, got %d", http.StatusNotFound, httpErr.Code)
					}
					if httpErr.Message != errRespNotFound {
						t.Errorf("expected error message %v, got %v", errRespNotFound, httpErr.Message)
					}
				}
			}
		})
	}
}

func Test_transport_GetCalculationById(t *testing.T) {
	ctx := newEchoCtx(http.MethodGet, "/calculations/a8098c1a-f86e-11da-bd1a-00112444be1e")
	ctx.SetParamNames("id")
	ctx.SetParamValues("a8098c1a-f86e-11da-bd1a-00112444be1e")

	ctxBadId := newEchoCtx(http.MethodGet, "/calculations/a8098c1a-f86e-11da-bd1a-")
	ctxBadId.SetParamNames("id")
	ctxBadId.SetParamValues("a8098c1a-f86e-11da-bd1a-")

	type fields struct {
		s func() Service
		l *zap.SugaredLogger
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			fields: fields{
				s: func() Service {
					ms := mocks.NewService(t)
					ms.EXPECT().GetCalculationById(domain.CalcID{ID: "a8098c1a-f86e-11da-bd1a-00112444be1e"}).
						Return(domain.Calculation{ID: "a8098c1a-f86e-11da-bd1a-00112444be1e", Expression: "1+1", Result: "2"}, nil)
					return ms
				},
				l: logger,
			},
			args: args{
				c: ctx,
			},
			wantErr: false,
		},
		{
			name: "failed validate id",
			fields: fields{
				s: func() Service { return nil },
				l: logger,
			},
			args: args{
				c: ctxBadId,
			},
			wantErr: true,
		},
		{
			name: "failde service case: notFound",
			fields: fields{
				s: func() Service {
					ms := mocks.NewService(t)
					ms.EXPECT().GetCalculationById(domain.CalcID{ID: "a8098c1a-f86e-11da-bd1a-00112444be1e"}).
						Return(domain.Calculation{}, domain.ErrNotFound)
					return ms
				},
				l: logger,
			},
			args: args{
				c: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := transport{
				s: tt.fields.s(),
				l: tt.fields.l,
			}
			if err := tr.GetCalculationById(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("transport.GetCalculationById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_transport_PostCalculation(t *testing.T) {
	type fields struct {
		s func() Service
		l *zap.SugaredLogger
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			fields: fields{
				s: func() Service {
					ms := mocks.NewService(t)
					ms.EXPECT().CreateCalculation(domain.CalcExpr{Expr: "1+1"}).
						Return(domain.Calculation{ID: "1", Expression: "1+1", Result: "2"}, nil)
					return ms
				},
				l: logger,
			},
			args: args{
				c: func() echo.Context {
					e := echo.New()
					req := httptest.NewRequest(http.MethodPost, "/calculations", strings.NewReader(`{"expression":"1+1"}`))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					rec := httptest.NewRecorder()
					return e.NewContext(req, rec)
				}(),
			},
			wantErr: false,
		},
		{
			name: "bad request - invalid JSON",
			fields: fields{
				s: func() Service { return nil },
				l: logger,
			},
			args: args{
				c: func() echo.Context {
					e := echo.New()
					req := httptest.NewRequest(http.MethodPost, "/calculations", strings.NewReader(`{"expression":`))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					rec := httptest.NewRecorder()
					return e.NewContext(req, rec)
				}(),
			},
			wantErr: true,
		},
		{
			name: "service error",
			fields: fields{
				s: func() Service {
					ms := mocks.NewService(t)
					ms.EXPECT().CreateCalculation(domain.CalcExpr{Expr: "1+1"}).
						Return(domain.Calculation{}, errors.New("service error"))
					return ms
				},
				l: logger,
			},
			args: args{
				c: func() echo.Context {
					e := echo.New()
					req := httptest.NewRequest(http.MethodPost, "/calculations", strings.NewReader(`{"expression":"1+1"}`))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					rec := httptest.NewRecorder()
					return e.NewContext(req, rec)
				}(),
			},
			wantErr: true,
		},
		{
			name: "validation error from service",
			fields: fields{
				s: func() Service {
					ms := mocks.NewService(t)
					ms.EXPECT().
						CreateCalculation(domain.CalcExpr{Expr: "1+1"}).
						Return(domain.Calculation{}, domain.ErrValidation)
					return ms
				},
				l: logger,
			},
			args: args{
				c: func() echo.Context {
					e := echo.New()
					req := httptest.NewRequest(http.MethodPost, "/calculations", strings.NewReader(`{"expression":"1+1"}`))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					rec := httptest.NewRecorder()
					return e.NewContext(req, rec)
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := transport{
				s: tt.fields.s(),
				l: tt.fields.l,
			}
			if err := tr.PostCalculation(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("transport.PostCalculation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_transport_DeleteCalcById(t *testing.T) {
	type fields struct {
		s func() Service
		l *zap.SugaredLogger
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			fields: fields{
				s: func() Service {
					ms := mocks.NewService(t)
					ms.EXPECT().DeleteCalcById(domain.CalcID{ID: "a8098c1a-f86e-11da-bd1a-00112444be1e"}).
						Return(nil)
					return ms
				},
				l: logger,
			},
			args: args{
				c: func() echo.Context {
					ctx := newEchoCtx(http.MethodDelete, "/calculations/a8098c1a-f86e-11da-bd1a-00112444be1e")
					ctx.SetParamNames("id")
					ctx.SetParamValues("a8098c1a-f86e-11da-bd1a-00112444be1e")
					return ctx
				}(),
			},
			wantErr: false,
		},
		{
			name: "invalid UUID",
			fields: fields{
				s: func() Service { return nil },
				l: logger,
			},
			args: args{
				c: func() echo.Context {
					ctx := newEchoCtx(http.MethodDelete, "/calculations/invalid-uuid")
					ctx.SetParamNames("id")
					ctx.SetParamValues("invalid-uuid")
					return ctx
				}(),
			},
			wantErr: true,
		},
		{
			name: "service error - not found",
			fields: fields{
				s: func() Service {
					ms := mocks.NewService(t)
					ms.EXPECT().DeleteCalcById(domain.CalcID{ID: "a8098c1a-f86e-11da-bd1a-00112444be1e"}).
						Return(domain.ErrNotFound)
					return ms
				},
				l: logger,
			},
			args: args{
				c: func() echo.Context {
					ctx := newEchoCtx(http.MethodDelete, "/calculations/a8098c1a-f86e-11da-bd1a-00112444be1e")
					ctx.SetParamNames("id")
					ctx.SetParamValues("a8098c1a-f86e-11da-bd1a-00112444be1e")
					return ctx
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := transport{
				s: tt.fields.s(),
				l: tt.fields.l,
			}
			if err := tr.DeleteCalcById(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("transport.DeleteCalcById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_transport_PatchCalculationById(t *testing.T) {
	type fields struct {
		s func() Service
		l *zap.SugaredLogger
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		check   func(t *testing.T, err error)
	}{
		{
			name: "successful case",
			fields: fields{
				s: func() Service {
					ms := mocks.NewService(t)
					ms.EXPECT().
						UpdateCalculationById(domain.Calculation{ID: "a8098c1a-f86e-11da-bd1a-00112444be1e", Expression: "1+2"}).
						Return(domain.Calculation{ID: "a8098c1a-f86e-11da-bd1a-00112444be1e", Expression: "1+2", Result: "3"}, nil)
					return ms
				},
				l: logger,
			},
			args: args{
				c: func() echo.Context {
					ctx := newEchoCtx(http.MethodPatch, "/calculations/a8098c1a-f86e-11da-bd1a-00112444be1e", `{"expression":"1+2"}`)
					ctx.SetParamNames("id")
					ctx.SetParamValues("a8098c1a-f86e-11da-bd1a-00112444be1e")
					ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					return ctx
				}(),
			},
			wantErr: false,
		},
		{
			name: "invalid UUID",
			fields: fields{
				s: func() Service { return nil },
				l: logger,
			},
			args: args{
				c: func() echo.Context {
					ctx := newEchoCtx(http.MethodPatch, "/calculations/invalid-uuid", `{"expression":"1+2"}`)
					ctx.SetParamNames("id")
					ctx.SetParamValues("invalid-uuid")
					ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					return ctx
				}(),
			},
			wantErr: true,
			check: func(t *testing.T, err error) {
				httpErr, ok := err.(*echo.HTTPError)
				if !ok || httpErr.Code != http.StatusBadRequest {
					t.Errorf("expected 400 Bad Request, got: %v", err)
				}
			},
		},
		{
			name: "invalid JSON",
			fields: fields{
				s: func() Service { return nil },
				l: logger,
			},
			args: args{
				c: func() echo.Context {
					ctx := newEchoCtx(http.MethodPatch, "/calculations/a8098c1a-f86e-11da-bd1a-00112444be1e", `{"expression":`)
					ctx.SetParamNames("id")
					ctx.SetParamValues("a8098c1a-f86e-11da-bd1a-00112444be1e")
					ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					return ctx
				}(),
			},
			wantErr: true,
			check: func(t *testing.T, err error) {
				httpErr, ok := err.(*echo.HTTPError)
				if !ok || httpErr.Code != http.StatusBadRequest {
					t.Errorf("expected 400 Bad Request, got: %v", err)
				}
			},
		},
		{
			name: "not found from service",
			fields: fields{
				s: func() Service {
					ms := mocks.NewService(t)
					ms.EXPECT().
						UpdateCalculationById(domain.Calculation{ID: "a8098c1a-f86e-11da-bd1a-00112444be1e", Expression: "1+2"}).
						Return(domain.Calculation{}, domain.ErrNotFound)
					return ms
				},
				l: logger,
			},
			args: args{
				c: func() echo.Context {
					ctx := newEchoCtx(http.MethodPatch, "/calculations/a8098c1a-f86e-11da-bd1a-00112444be1e", `{"expression":"1+2"}`)
					ctx.SetParamNames("id")
					ctx.SetParamValues("a8098c1a-f86e-11da-bd1a-00112444be1e")
					ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					return ctx
				}(),
			},
			wantErr: true,
			check: func(t *testing.T, err error) {
				httpErr, ok := err.(*echo.HTTPError)
				if !ok || httpErr.Code != http.StatusNotFound {
					t.Errorf("expected 404 Not Found, got: %v", err)
				}
			},
		},
		{
			name: "internal error from service",
			fields: fields{
				s: func() Service {
					ms := mocks.NewService(t)
					ms.EXPECT().
						UpdateCalculationById(domain.Calculation{ID: "a8098c1a-f86e-11da-bd1a-00112444be1e", Expression: "1+2"}).
						Return(domain.Calculation{}, errors.New("unexpected error"))
					return ms
				},
				l: logger,
			},
			args: args{
				c: func() echo.Context {
					ctx := newEchoCtx(http.MethodPatch, "/calculations/a8098c1a-f86e-11da-bd1a-00112444be1e", `{"expression":"1+2"}`)
					ctx.SetParamNames("id")
					ctx.SetParamValues("a8098c1a-f86e-11da-bd1a-00112444be1e")
					ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					return ctx
				}(),
			},
			wantErr: true,
			check: func(t *testing.T, err error) {
				httpErr, ok := err.(*echo.HTTPError)
				if !ok || httpErr.Code != http.StatusInternalServerError {
					t.Errorf("expected 500 Internal Server Error, got: %v", err)
				}
			},
		},
		{
			name: "validation error from service",
			fields: fields{
				s: func() Service {
					ms := mocks.NewService(t)
					ms.EXPECT().
						UpdateCalculationById(domain.Calculation{
							ID:         "a8098c1a-f86e-11da-bd1a-00112444be1e",
							Expression: "1+2",
						}).
						Return(domain.Calculation{}, domain.ErrValidation)
					return ms
				},
				l: logger,
			},
			args: args{
				c: func() echo.Context {
					ctx := newEchoCtx(http.MethodPatch, "/calculations/a8098c1a-f86e-11da-bd1a-00112444be1e", `{"expression":"1+2"}`)
					ctx.SetParamNames("id")
					ctx.SetParamValues("a8098c1a-f86e-11da-bd1a-00112444be1e")
					ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					return ctx
				}(),
			},
			wantErr: true,
			check: func(t *testing.T, err error) {
				httpErr, ok := err.(*echo.HTTPError)
				if !ok || httpErr.Code != http.StatusBadRequest {
					t.Errorf("expected 400 Bad Request, got: %v", err)
				}
				if httpErr.Message != errRespValidation {
					t.Errorf("expected error message %v, got %v", errRespValidation, httpErr.Message)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := transport{
				s: tt.fields.s(),
				l: tt.fields.l,
			}
			err := tr.PatchCalculationById(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("transport.PatchCalculationById() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.check != nil {
				tt.check(t, err)
			}
		})
	}
}
