package response

type ICell interface {
	SetName(name string)
	GetName() string
}
type IRow[C ICell] interface {
	SetCells(cells []C)
	AddCells(cells ...C)
	GetCells() []C
}

type ITableData[C ICell, R IRow[C]] interface {
	SetRows(rows []R)
	AddRows(rows ...R)
	GetRows() []R
	SetHeadings(h R)
	GetHeadings() R
}

// ITimings handles simple start, end and duration elements of the interface.
type ITimings interface {
	Start()
	End()
}

// IStatus handles tracking the http status of the api response.
// Its value should be used with IApi.Write call at the end
type IStatus interface {
	SetStatus(status int)
	GetStatus() int
}

// IErrors allows tracking of server side errors such as validation
// and will be included in the IApi.Write
type IErrors interface {
	SetErrors(errors []error)
	AddError(err error)
	GetErrors() []error
}

// IBase is a merge interface that wuld be typical of an http response.
// This version excludes any result data / handling for simplicty on errors or
// empty results
type IBase interface {
	ITimings
	IStatus
	IErrors
	AddErrorWithStatus(err error, status int)
}

// IResult providers a response interface whose result type can vary between
// slice, a map or a map of slices.
// This allows api respsones to adapt to the most useful data type for the endpoint
type IResult[C ICell, R IRow[C], D ITableData[C, R]] interface {
	IBase
	SetResult(result D)
	GetResult() D
}
