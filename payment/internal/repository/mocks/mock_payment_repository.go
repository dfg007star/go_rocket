// Code generated for dosin service
// © dosin 2025.

// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/dfg007star/go_rocket/payment/internal/repository/model"
	mock "github.com/stretchr/testify/mock"
)

// PaymentRepository is an autogenerated mock type for the PaymentRepository type
type PaymentRepository struct {
	mock.Mock
}

type PaymentRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *PaymentRepository) EXPECT() *PaymentRepository_Expecter {
	return &PaymentRepository_Expecter{mock: &_m.Mock}
}

// PayOrder provides a mock function with given fields: ctx, payment
func (_m *PaymentRepository) PayOrder(ctx context.Context, payment model.Payment) (string, error) {
	ret := _m.Called(ctx, payment)

	if len(ret) == 0 {
		panic("no return value specified for PayOrder")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.Payment) (string, error)); ok {
		return rf(ctx, payment)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.Payment) string); ok {
		r0 = rf(ctx, payment)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.Payment) error); ok {
		r1 = rf(ctx, payment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PaymentRepository_PayOrder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PayOrder'
type PaymentRepository_PayOrder_Call struct {
	*mock.Call
}

// PayOrder is a helper method to define mock.On call
//   - ctx context.Context
//   - payment model.Payment
func (_e *PaymentRepository_Expecter) PayOrder(ctx interface{}, payment interface{}) *PaymentRepository_PayOrder_Call {
	return &PaymentRepository_PayOrder_Call{Call: _e.mock.On("PayOrder", ctx, payment)}
}

func (_c *PaymentRepository_PayOrder_Call) Run(run func(ctx context.Context, payment model.Payment)) *PaymentRepository_PayOrder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.Payment))
	})
	return _c
}

func (_c *PaymentRepository_PayOrder_Call) Return(_a0 string, _a1 error) *PaymentRepository_PayOrder_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PaymentRepository_PayOrder_Call) RunAndReturn(run func(context.Context, model.Payment) (string, error)) *PaymentRepository_PayOrder_Call {
	_c.Call.Return(run)
	return _c
}

// NewPaymentRepository creates a new instance of PaymentRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPaymentRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *PaymentRepository {
	mock := &PaymentRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
