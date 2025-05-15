package table

import "go.uber.org/zap/zapcore"

type Table struct {
	Name  string
	Count int
}

func (t *Table) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", t.Name)
	enc.AddInt("count", t.Count)
	return nil
}

func MarshalTables(tables []Table) zapcore.ArrayMarshaler {
	return zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, t := range tables {
			if err := enc.AppendObject(&t); err != nil {
				return err
			}
		}
		return nil
	})
}
