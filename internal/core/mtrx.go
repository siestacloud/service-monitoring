package core

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
)

//Metric Метрика
type Metric struct {
	ID    string   `json:"id" db:"name"`               // имя метрики
	MType string   `json:"type" db:"type"`             // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty" db:"delta"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty" db:"value"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`             // значение хеш-функции
}

//СОздает метрику
func NewMetric() *Metric {
	return &Metric{}
}

//Устанавливаем ID
func (m *Metric) SetID(id string) error {
	if err := m.checkID(id); err != nil {
		return err
	}
	m.ID = id
	return nil
}

//Получить ID
func (m *Metric) GetID() string {
	return m.ID
}

//Проверка названия метрики
func (m *Metric) checkID(id string) error {
	if id == "" {
		return errors.New("Unable set id for metric <" + id + ">: empty id ")
	}
	if len(id) > 50 {
		return errors.New("Unable set id for metric <" + id + ">: len too big ")
	}
	return nil
}

//Устанавливаем Тип
func (m *Metric) SetType(t string) error {
	if err := m.checkType(t); err != nil {
		return err
	}
	m.MType = t
	return nil
}

//Получить Тип
func (m *Metric) GetType() string {
	return m.MType
}

//Проверка типа метрики
func (m *Metric) checkType(t string) error {
	if t == "" {
		return errors.New("Unable set type for metric <" + m.GetID() + ">: empty type ")
	}
	if t == "counter" {
		return nil
	}
	if t == "gauge" {
		return nil
	}
	return errors.New("Unable set type for metric <" + m.GetID() + ">: undefined type ")
}

//Устанавливаем Значение
func (m *Metric) SetValue(v interface{}) error {
	if m.ID == "" || m.MType == "" {
		return errors.New("unable set value for metric ID or type for mtrx undefined")
	}
	if err := m.checkValue(v); err != nil {
		return err
	}

	switch m.GetType() {
	case "gauge":
		sval, ok := v.(string)
		if ok {
			val, err := strconv.ParseFloat(sval, 64)
			if err != nil {
				return err
			}
			m.Value = &val
			return nil
		}
		val, ok := v.(float64)
		if ok {
			m.Value = &val
			return nil
		}
		return errors.New("unable set value for metric <" + m.GetID() + ">: incorrect type value for gauge mtrx ")
	case "counter":

		sval, ok := v.(string)
		if ok {
			val, err := strconv.ParseInt(sval, 10, 64)
			if err != nil {
				return err
			}

			m.Delta = &val
			return nil
		}
		val, ok := v.(int64)
		if ok {
			m.Delta = &val
			return nil
		}
		return errors.New("Unable set value for metric <" + m.GetID() + ">: incorrect type value for gauge mtrx ")

	default:
		return errors.New("Unable set value for metric <" + m.GetID() + ">: mtrx have incorect type")
	}
}

//Получить Значание
func (m *Metric) GetValue() (float64, error) {
	if m.Value != nil {
		return *m.Value, nil
	}
	return 0, errors.New("unable get value for mtrx <" + m.GetID() + ">: nil pointer to value")
}

//Получить Значание
func (m *Metric) GetDelta() (int64, error) {
	if m.Delta != nil {
		return *m.Delta, nil
	}
	return 0, errors.New("unable get delta for mtrx <" + m.GetID() + ">: nil pointer to value")
}

//Проверка Значения метрики
func (m *Metric) checkValue(v interface{}) error {
	switch v.(type) {
	case float64:
		return nil
	case int64:
		return nil
	case string:
		return nil
	case nil:
		return errors.New("Unable set value for metric <" + m.GetID() + ">: value type empty ")
	default:
		return errors.New("Unable set value for metric <" + m.GetID() + ">: value type bug ")
	}
}

// WriteMetricJSON сериализует структуру Metric в JSON, и если всё отрабатывает
// успешно, то вызывается соответствующий метод Write() из io.Writer.
func (m *Metric) MarshalMetricsinJSON(w io.Writer) error {
	js, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = w.Write(js)
	return err
}

//Создает метрики и чекает ее поля из JSON
func (m *Metric) UnmarshalMetricJSON(data []byte) error {

	mtrxCheck := NewMetric() // Промежуточный обьект, поля которого будут проверены
	// if err := json.NewDecoder(r).Decode(&mtrxCheck); err != nil {
	// 	return err
	// }
	if err := json.Unmarshal(data, mtrxCheck); err != nil {
		return err
	}

	err := m.SetID(mtrxCheck.GetID()) // Проверка имени метрики
	if err != nil {
		return err // Проверка не пройдена
	}
	err = m.SetType(mtrxCheck.GetType()) //Проверка типа метрики
	if err != nil {
		return err //Проверка не пройдена
	}

	m.Hash = mtrxCheck.Hash

	if m.GetType() == "counter" { //Если у новой метрики тип counter
		d, err := mtrxCheck.GetDelta()
		if err != nil {
			return err
		}
		m.SetValue(d) //Присваиваем дельту int64
		return nil
	}
	if m.GetType() == "gauge" {
		v, err := mtrxCheck.GetValue() //Если у новой метрики тип gauge
		if err != nil {
			return err
		}
		m.SetValue(v) //Присваиваем value float64/
		return nil
	}
	return errors.New("something incorrect")
}

//Создает метрики и чекает ее поля из JSON
func UnmarshalMetricCaseJSON(data []byte) ([]Metric, error) {
	mtrxCaseCheck := []Metric{} // пул полученных от клиента метрик, не проверенных
	mtrxCase := []Metric{}      // пул полученных от клиента метрик, проверенных

	// if err := json.NewDecoder(r).Decode(&mtrxCheck); err != nil {
	// 	return err
	// }

	if err := json.Unmarshal(data, &mtrxCaseCheck); err != nil {
		return nil, err
	}

	for _, mtrxCheck := range mtrxCaseCheck {

		mtrx := NewMetric()

		err := mtrx.SetID(mtrxCheck.GetID()) // Проверка имени метрики
		if err != nil {
			return nil, err // Проверка не пройдена
		}
		err = mtrx.SetType(mtrxCheck.GetType()) //Проверка типа метрики
		if err != nil {
			return nil, err //Проверка не пройдена
		}

		mtrx.Hash = mtrxCheck.Hash

		if mtrx.GetType() == "counter" { //Если у новой метрики тип counter
			d, err := mtrxCheck.GetDelta()
			if err != nil {
				return nil, err
			}
			if err := mtrx.SetValue(d); err != nil {
				return nil, err
			}
			mtrxCase = append(mtrxCase, *mtrx)
			continue
		}
		if mtrx.GetType() == "gauge" {
			fmt.Println("OKS")
			v, err := mtrxCheck.GetValue() //Если у новой метрики тип gauge
			if err != nil {
				return nil, err
			}
			if err := mtrx.SetValue(v); err != nil {
				return nil, err
			} //Присваиваем value float64/
			mtrxCase = append(mtrxCase, *mtrx)
			continue
		}

	}
	return mtrxCase, nil

}

// Вычисляем hash метрики
func (m *Metric) SetHash(key string) error {
	if key != "" {
		var sign string
		switch m.GetType() {
		case "counter":
			sign = hash(fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta), key)

		case "gauge":
			sign = hash(fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value), key)
		default:
			return errors.New("unable set hash invalid type mtrx")
		}
		m.Hash = sign
		return nil
	}
	return nil
}

//Получить ID
func (m *Metric) GetHash() string {
	return m.Hash
}

func hash(str, key string) string {

	var secretkey = []byte(key)
	var data = []byte(str)

	h := hmac.New(sha256.New, secretkey)
	h.Write(data)
	sign := h.Sum(nil)
	return fmt.Sprintf("%x", sign)
}
