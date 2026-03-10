package interfaces

import (
	"context"
	"opg-reports/report/packages/args"
)

// ImportGetterF is a function that will return data - optionally accepts previous data to allow pass along
//
// For example, fetching all code repositories and then passing that to the command to generate code
// ownership; rather than repeating the function
type ImportGetterF[T any, O any, C any] func(ctx context.Context, client C, opts *args.Import, in ...O) (data []T, err error)

// ImportFilterF allows for data filtering based on struct values - so after api has been called.
//
// Normally used to reduce the data set before insert on simple values such as names
type ImportFilterF[T any] func(ctx context.Context, data []T, opts *args.Filters) (filtered []T)

// ImportTransformF converts the raw data ([]T) into the data to write to database ([]R)
//
// Used to simplify data from
type ImportTransformF[R Insertable, T any] func(ctx context.Context, data []T, opts *args.Import) (result []R, err error)
