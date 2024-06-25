package cost

import (
	"opg-reports/shared/server"
)

func ResponseAsStringsMap(content []byte) (*server.ApiResponse[*Cost, map[string]*Cost], error) {
	return server.NewApiResponseFromJson[*Cost, map[string]*Cost](content)
}
func ResponseAsStringsMapSlice(content []byte) (*server.ApiResponse[*Cost, map[string][]*Cost], error) {
	return server.NewApiResponseFromJson[*Cost, map[string][]*Cost](content)
}
func ResponseAsStringsMapMapSlice(content []byte) (*server.ApiResponse[*Cost, map[string]map[string][]*Cost], error) {
	return server.NewApiResponseFromJson[*Cost, map[string]map[string][]*Cost](content)
}
func ResponseAsStringsSlice(content []byte) (*server.ApiResponse[*Cost, []*Cost], error) {
	return server.NewApiResponseFromJson[*Cost, []*Cost](content)
}
