package serializer

import (
	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
)

type EasyJsonSerializer struct{}

var defaultSerializer = echo.DefaultJSONSerializer{}

func (e EasyJsonSerializer) Serialize(c echo.Context, i any, indent string) error {
	obj, ok := i.(easyjson.Marshaler)
	if !ok {
		return defaultSerializer.Serialize(c, i, indent)
	}
	_, err := easyjson.MarshalToWriter(obj, c.Response())
	return err
}

func (e EasyJsonSerializer) Deserialize(c echo.Context, i any) error {
	obj, ok := i.(easyjson.Unmarshaler)
	if !ok {
		return defaultSerializer.Deserialize(c, i)
	}
	return easyjson.UnmarshalFromReader(c.Request().Body, obj)
}
