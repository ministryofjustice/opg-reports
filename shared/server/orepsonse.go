package server

type IApiResponse interface {
	Start()
	End()

	SetResults(results interface{})
	GetResults() interface{}

	SetStatus(status int)
	GetStatus() int

	SetErrors(errors []error)
	AddError(err error)
	AddStatusError(status int, err error)
	GetErrors() []error

	Body() []byte
}

// type ApiTimeData struct {
// 	Start    time.Time     `json:"start"`
// 	End      time.Time     `json:"end"`
// 	Duration time.Duration `json:"duration"`
// }

// type ApiResponse[T ApiListResult] struct {
// 	Times  ApiTimeData `json:"timing"`
// 	Errors []error     `json:"errors"`
// 	Status int         `json:"status"`
// 	Result T           `json:"result"`
// }

// func (r *ApiResponse) Body() []byte {
// 	body, _ := json.Marshal(r)
// 	return body
// }

// func (r *ApiResponse) Start() {
// 	r.RequestStart = time.Now().UTC()
// }

// func (r *ApiResponse) End() {
// 	r.RequestEnd = time.Now().UTC()
// 	r.RequestDuration = r.RequestEnd.Sub(r.RequestStart)
// }

// func (r *ApiResponse) SetResults(results interface{}) {
// 	r.Result = results
// }

// func (r *ApiResponse) GetResults() interface{} {
// 	return r.Result
// }

// func (r *ApiResponse) SetStatus(status int) {
// 	r.Status = status
// }
// func (r *ApiResponse) GetStatus() int {
// 	return r.Status
// }

// func (r *ApiResponse) SetErrors(errors []error) {
// 	r.Errors = errors
// }

// func (r *ApiResponse) AddError(err error) {
// 	r.Errors = append(r.Errors, err)
// }

// func (r *ApiResponse) AddStatusError(status int, err error) {
// 	r.Errors = append(r.Errors, err)
// 	r.SetStatus(status)
// }

// func (r *ApiResponse) GetErrors() []error {
// 	return r.Errors
// }
