package mongo

import (
    "context"
    "encoding/json"
    "errors"
    "net/http"
)

const (
    KEY_ID = "id"
    KEY_NAME = "name"
    KEY_DESCRIPTION = "desc"
)

func decodeFruitRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var fruitRequest FruitRequest
    if err := json.NewDecoder(r.Body).Decode(&fruitRequest); err != nil {
        return nil, err
    }

    return fruitRequest, nil
}

func decodeInsertFruitRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    return decodeFruitRequest(ctx, r)
}

func decodeGetFruitsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    request := FruitRequest{}

    keys := r.URL.Query()
    id := keys.Get(KEY_ID)
    name := keys.Get(KEY_NAME)
    description := keys.Get(KEY_DESCRIPTION)

    if id == "" && name == "" && description == "" {
        return nil, errors.New("bad request. Could not find any of (id or name or desc) query params")
    }

    request.Id = id
    request.Name = name
    request.Description = description

    return request, nil
}

func encodeFruitResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
    setContentType(w)

    return json.NewEncoder(w).Encode(response)
}
