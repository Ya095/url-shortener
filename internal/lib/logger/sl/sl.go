package sl

import "log/slog"


// Err - возвращает атрибут с ошибкой в виде строки.
func Err(err error) slog.Attr {
	return slog.Attr{
		Key: "error",
		Value: slog.StringValue(err.Error()),
	}
}