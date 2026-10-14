package response_test

import (
	"fmt"
	"net/http/httptest"

	"github.com/mateusz-kacpura/integra-core/pkg/response"
)

func ExampleWriteOK() {
	rec := httptest.NewRecorder()
	_ = response.WriteOK(rec, map[string]string{"status": "healthy"})

	fmt.Println("Status Code:", rec.Code)
	fmt.Println("Body:", rec.Body.String())
	// Output:
	// Status Code: 200
	// Body: {"success":true,"code":"OK","data":{"status":"healthy"}}
}

func ExampleWriteError() {
	rec := httptest.NewRecorder()
	_ = response.WriteError(rec, 404, response.CodeNotFound, "entity not found")

	fmt.Println("Status Code:", rec.Code)
	fmt.Println("Body:", rec.Body.String())
	// Output:
	// Status Code: 404
	// Body: {"success":false,"code":"NOT_FOUND","message":"entity not found"}
}
