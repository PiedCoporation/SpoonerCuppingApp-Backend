package abstractions

import (
	"backend/internal/contracts/common"
	"backend/internal/contracts/sample"
	"context"
)

type ISampleService interface {
	GetAll(ctx context.Context, pageSize int, pageNumber int) (*common.Result[common.PageResult[sample.SampleRes]])
	Create(ctx context.Context, req sample.SampleReq) (*common.Result[sample.SampleRes])
}
