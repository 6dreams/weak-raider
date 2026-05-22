package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Translation struct {
	Data map[string]string `json:"translation"`
}

func (t *Translation) Get(locale string) string {
	r, ok := t.Data[locale]

	if !ok {
		return t.Default()
	}

	return r
}

func (t *Translation) Default() string {
	return t.Data["en_US"]
}

func (t *Translation) Value() (driver.Value, error) {
	d, err := json.Marshal(t.Data)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (t *Translation) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal translations data")
	}
	return json.Unmarshal(bytes, &t)
}

func NewTranslation(data map[string]string) *Translation {
	return &Translation{Data: data}
}
