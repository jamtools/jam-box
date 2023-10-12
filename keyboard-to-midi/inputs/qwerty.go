package inputs

import (
	"github.com/pkg/errors"
)

type QwertyInput struct {
	Name      string
	EventPath string
	c         chan string
}

func (q *QwertyInput) SubscribeToInput() (chan string, error) {
	if q.c != nil {
		return nil, errors.Errorf("there is already a subscriber to this qwerty input: %s", q.Name)
	}

	q.c = make(chan string)
	return q.c, nil
}

func NewQwertyInput(name string) *QwertyInput {
	q := &QwertyInput{Name: name}

}
