package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	mocks "github.com/sgahlot/go-mongo-service-quickstart/mocks/pkg/common"
	"github.com/sgahlot/go-mongo-service-quickstart/pkg/common"
	"io/ioutil"

	"github.com/sgahlot/go-mongo-service-quickstart/pkg/db"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	INITIAL_4_FRUITS = `
        {
            "description": "Good for health",
            "id": "1",
            "name": "Banana"
        },
        {
            "description": "Keeps the doctor away",
            "id": "2",
            "name": "Apple"
        },
        {
            "description": "Antioxidant Superfood",
            "id": "3",
            "name": "Blueberry"
        },
        {
            "description": "Healing fruit",
            "id": "4",
            "name": "Peach"
        }`

	INSERT_FRUIT_RESPONSE = `
		{
			"id":      "1010101",
			"message": "success",
		}
	`
)

func performRequest(handler http.Handler, body interface{}, method, path string) *httptest.ResponseRecorder {
	req := prepareRequestBody(body, method, path)
	resWriter := httptest.NewRecorder()

	handler.ServeHTTP(resWriter, req)
	return resWriter
}

func prepareRequestBody(body interface{}, method, path string) *http.Request {
	reqBytes := new(bytes.Buffer)
	json.NewEncoder(reqBytes).Encode(body)

	req, _ := http.NewRequest(method, path, ioutil.NopCloser(bytes.NewBuffer(reqBytes.Bytes())))
	return req
}

func TestRetrieveFruits(t *testing.T) {
	mockService := new(mocks.Service)
	router := mongoRoute(mockService)

	t.Run("Invalid query params", func(t *testing.T) {
		writer := performRequest(router, nil, "GET", "/api/v1/fruits?name1="+db.ALL_ROWS)
		assert.Equal(t, http.StatusInternalServerError, writer.Code)

		var response interface{}
		err := json.Unmarshal([]byte(writer.Body.String()), &response)

		expectedResponse := make(map[string]interface{})
		expectedResponse["error"] = "bad request. Could not find any of (id or name or desc) query params"

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, response)
	})

	t.Run("Non-existent fruit", func(t *testing.T) {
		RESPONSE_MESSAGE_NO_FRUITS_FOUND := "No fruits found for given query"

		fruitRequest := common.Fruit{Name: "BLAH"}

		responseJson := `{
            "message": "%s"
        }`
		responseJson = fmt.Sprintf(responseJson, RESPONSE_MESSAGE_NO_FRUITS_FOUND)

		var responseAsObj common.FruitResponse
		json.Unmarshal([]byte(responseJson), &responseAsObj)
		mockService.On("GetFruits", &fruitRequest).Return(responseAsObj)

		// fmt.Printf("Expected response: %+v\n", responseAsObj)

		response, err := invokeApiAndVerifyResponse(t, router, nil, "GET", "/api/v1/fruits?name=BLAH", http.StatusOK)
		fruitResponse := response.(common.FruitResponse)

		assert.Nil(t, err)
		assert.Nil(t, fruitResponse.Err)
		assert.Equal(t, fruitResponse.Message, RESPONSE_MESSAGE_NO_FRUITS_FOUND)
		assert.Nil(t, fruitResponse.Fruits)
	})
}

func TestRetrieveAllFruits(t *testing.T) {
	mockService := new(mocks.Service)
	router := mongoRoute(mockService)

	fruitRequest := common.Fruit{Name: db.ALL_ROWS}

	RESPONSE_MESSAGE := "Found 4 fruits"

	responseJson := getFruitResponseJson("["+INITIAL_4_FRUITS+"]", RESPONSE_MESSAGE)

	var responseAsObj common.FruitResponse
	json.Unmarshal([]byte(responseJson), &responseAsObj)
	mockService.On("GetFruits", &fruitRequest).Return(responseAsObj)

	// fmt.Printf("Expected response: %+v\n", responseAsObj)

	response, err := invokeApiAndVerifyResponse(t, router, nil, "GET", "/api/v1/fruits?name="+db.ALL_ROWS, http.StatusOK)
	fruitResponse := response.(common.FruitResponse)

	assert.Nil(t, err)
	assert.Nil(t, fruitResponse.Err)
	assert.Equal(t, fruitResponse.Message, RESPONSE_MESSAGE)
	assert.Equal(t, responseAsObj.Fruits, fruitResponse.Fruits)
}

func TestInsertFruit(t *testing.T) {
	mockService := new(mocks.Service)
	router := mongoRoute(mockService)

	fruitRequest := common.Fruit{Description: "Full of fibre adn Vitamin C", Name: "Pear"}
	responseAsObj := getFruitResponseAsObject(INSERT_FRUIT_RESPONSE)
	mockService.On("InsertFruit", &fruitRequest).Return(responseAsObj)
	response, err := invokeApiAndVerifyResponse(t, router,
		fruitRequest,
		"POST",
		"/api/v1/fruits",
		http.StatusOK)

	RESPONSE_MESSAGE := "Found 5 fruits"
	responseJson := getFruitResponseJson(fmt.Sprintf(`[%s, {"description": "%s", "name": "%s"}]`, INITIAL_4_FRUITS, fruitRequest.Description, fruitRequest.Name),
		RESPONSE_MESSAGE)

	responseAsObj = getFruitResponseAsObject(responseJson)
	fruitRequest = common.Fruit{Name: db.ALL_ROWS}
	mockService.On("GetFruits", &fruitRequest).Return(responseAsObj)

	// fmt.Printf("Expected response: %+v\n", responseAsObj)

	response, err = invokeApiAndVerifyResponse(t, router, nil, "GET", "/api/v1/fruits?name="+db.ALL_ROWS, http.StatusOK)
	fruitResponse := response.(common.FruitResponse)

	assert.Nil(t, err)
	assert.Nil(t, fruitResponse.Err)
	assert.Equal(t, fruitResponse.Message, RESPONSE_MESSAGE)
	assert.Equal(t, responseAsObj.Fruits, fruitResponse.Fruits)
}

func getFruitResponseAsObject(responseJson string) common.FruitResponse {
	var responseAsObj common.FruitResponse
	json.Unmarshal([]byte(responseJson), &responseAsObj)

	return responseAsObj
}

func getFruitResponseJson(fruits, message string) string {
	responseJson := `{
        "fruits": %s,
        "message": "%s"
    }`

	responseJson = fmt.Sprintf(responseJson, fruits, message)
	return responseJson
}

func invokeApiAndVerifyResponse(t *testing.T, router http.Handler, body interface{}, method, path string, httpStatus int) (interface{}, error) {
	writer := performRequest(router, body, method, path)
	assert.Equal(t, httpStatus, writer.Code)

	var response common.FruitResponse
	err := json.Unmarshal([]byte(writer.Body.String()), &response)

	// fmt.Printf("Error: %v\n", err)
	// fmt.Printf("response: %#v\n", response)

	return response, err
}
