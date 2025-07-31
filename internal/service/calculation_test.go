package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/eragon-mdi/calc-back/internal/domain"
	"github.com/eragon-mdi/calc-back/internal/service/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_service_GetLastCalculations(t *testing.T) {
	type fields struct {
		r Repository
	}

	mockCalcs := []domain.Calculation{
		{ID: "1", Expression: "1+1", Result: "2"},
	}

	tests := []struct {
		name    string
		fields  fields
		want    []domain.Calculation
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				r: func() Repository {
					m := mocks.NewRepository(t)
					m.On("GetCalculations", Max_Calcs).Return(mockCalcs, nil).Once()
					return m
				}(),
			},
			want:    mockCalcs,
			wantErr: false,
		},
		{
			name: "empty slice returns ErrNotFound",
			fields: fields{
				r: func() Repository {
					m := mocks.NewRepository(t)
					m.On("GetCalculations", Max_Calcs).Return([]domain.Calculation{}, nil).Once()
					return m
				}(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "repository error",
			fields: fields{
				r: func() Repository {
					m := mocks.NewRepository(t)
					m.On("GetCalculations", Max_Calcs).Return(nil, errors.New("some error")).Once()
					return m
				}(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				r: tt.fields.r,
			}
			got, err := s.GetLastCalculations()
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetLastCalculations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetLastCalculations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetCalculationById(t *testing.T) {
	mockCalc := domain.Calculation{ID: "1", Expression: "1+1", Result: "2"}

	type fields struct {
		r Repository
	}
	type args struct {
		id domain.CalcID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.Calculation
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				r: func() Repository {
					m := mocks.NewRepository(t)
					m.On("GetCalculation", "1").Return(mockCalc, nil).Once()
					return m
				}(),
			},
			args:    args{id: domain.CalcID{ID: "1"}},
			want:    mockCalc,
			wantErr: false,
		},
		{
			name: "repo error",
			fields: fields{
				r: func() Repository {
					m := mocks.NewRepository(t)
					m.On("GetCalculation", "1").Return(domain.Calculation{}, errors.New("not found")).Once()
					return m
				}(),
			},
			args:    args{id: domain.CalcID{ID: "1"}},
			want:    domain.Calculation{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				r: tt.fields.r,
			}
			got, err := s.GetCalculationById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetCalculationById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetCalculationById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_CreateCalculation(t *testing.T) {
	exprValid := domain.CalcExpr{Expr: "1+2"}
	exprInvalid := domain.CalcExpr{Expr: "1++2"}

	mockSavedCalc := domain.Calculation{
		ID:         "uuid-generated", // этот ID вернём из моковой функции
		Expression: "1+2",
		Result:     "3",
	}

	type fields struct {
		r Repository
	}
	type args struct {
		expr domain.CalcExpr
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.Calculation
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				r: func() Repository {
					m := mocks.NewRepository(t)
					m.On("SaveTask", mock.MatchedBy(func(calc domain.Calculation) bool {
						return calc.Expression == "1+2" && calc.Result == "3"
					})).Return(mockSavedCalc, nil).Once()
					return m
				}(),
			},
			args:    args{expr: exprValid},
			want:    mockSavedCalc,
			wantErr: false,
		},
		{
			name: "calculate returns validation error",
			fields: fields{
				r: mocks.NewRepository(t),
			},
			args:    args{expr: exprInvalid},
			want:    domain.Calculation{},
			wantErr: true,
		},
		{
			name: "repository SaveTask error",
			fields: fields{
				r: func() Repository {
					m := mocks.NewRepository(t)
					m.On("SaveTask", mock.MatchedBy(func(calc domain.Calculation) bool {
						return calc.Expression == "1+2" && calc.Result == "3"
					})).Return(domain.Calculation{}, errors.New("db error")).Once()
					return m
				}(),
			},
			args:    args{expr: exprValid},
			want:    domain.Calculation{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				r: tt.fields.r,
			}
			got, err := s.CreateCalculation(tt.args.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.CreateCalculation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.CreateCalculation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_DeleteCalcById(t *testing.T) {
	type fields struct {
		r Repository
	}
	type args struct {
		id domain.CalcID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				r: func() Repository {
					m := mocks.NewRepository(t)
					m.On("DeleteCalculation", "1").Return(nil).Once()
					return m
				}(),
			},
			args:    args{id: domain.CalcID{ID: "1"}},
			wantErr: false,
		},
		{
			name: "repository error",
			fields: fields{
				r: func() Repository {
					m := mocks.NewRepository(t)
					m.On("DeleteCalculation", "1").Return(errors.New("delete error")).Once()
					return m
				}(),
			},
			args:    args{id: domain.CalcID{ID: "1"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				r: tt.fields.r,
			}
			if err := s.DeleteCalcById(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("service.DeleteCalcById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_UpdateCalculationById(t *testing.T) {
	calcValid := domain.Calculation{
		ID:         "1",
		Expression: "1+2",
		Result:     "3",
	}

	calcInvalid := domain.Calculation{
		ID:         "1",
		Expression: "1++2", // некорректно для calculate
	}

	mockUpdatedCalc := domain.Calculation{
		ID:         "1",
		Expression: "1+2",
		Result:     "3",
	}

	type fields struct {
		r Repository
	}
	type args struct {
		calc domain.Calculation
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.Calculation
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				r: func() Repository {
					m := mocks.NewRepository(t)
					m.On("UpdateTaskInfo", calcValid).Return(mockUpdatedCalc, nil).Once()
					return m
				}(),
			},
			args:    args{calc: calcValid},
			want:    mockUpdatedCalc,
			wantErr: false,
		},
		{
			name: "calculate validation error",
			fields: fields{
				r: mocks.NewRepository(t), // UpdateTaskInfo не вызывается
			},
			args:    args{calc: calcInvalid},
			want:    domain.Calculation{},
			wantErr: true,
		},
		{
			name: "repository update error",
			fields: fields{
				r: func() Repository {
					m := mocks.NewRepository(t)
					m.On("UpdateTaskInfo", calcValid).Return(domain.Calculation{}, errors.New("update error")).Once()
					return m
				}(),
			},
			args:    args{calc: calcValid},
			want:    domain.Calculation{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				r: tt.fields.r,
			}
			got, err := s.UpdateCalculationById(tt.args.calc)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.UpdateCalculationById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.UpdateCalculationById() = %v, want %v", got, tt.want)
			}
		})
	}
}
