//go:build !solution

package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Evaluator struct {
	keyWords   map[string][]string
	operations map[string]func(*Evaluator) error
	stack      []int
}

func NewEvaluator() *Evaluator {
	e := &Evaluator{
		keyWords:   make(map[string][]string),
		operations: make(map[string]func(*Evaluator) error),
		stack:      []int{},
	}

	e.operations["+"] = e.addOperation
	e.operations["-"] = e.subOperation
	e.operations["*"] = e.mulOperation
	e.operations["/"] = e.divOperation
	e.operations["dup"] = e.dupOperation
	e.operations["drop"] = e.dropOperation
	e.operations["swap"] = e.swapOperation
	e.operations["over"] = e.overOperation
	return e
}

func (e *Evaluator) execute(token string) error {
	token = strings.ToLower(token)

	if def, exists := e.keyWords[token]; exists {
		return e.executeDefinition(def)
	}

	if operation, exists := e.operations[token]; exists {
		return operation(e)
	}

	if value, err := strconv.Atoi(token); err == nil {
		e.stack = append(e.stack, value)
		return nil
	}

	return fmt.Errorf("undefined word: %s", token)
}

func (e *Evaluator) freezeOperations(def []string) []string {
	var result []string
	for _, token := range def {
		if _, exists := e.keyWords[token]; exists {
			result = append(result, e.keyWords[token]...)
		} else if _, err := strconv.Atoi(token); err == nil {
			result = append(result, token)
		} else if _, exists := e.operations[token]; exists {
			result = append(result, "default_"+token)
		}
	}
	return result
}

func (e *Evaluator) executeFrozen(token string) error {
	baseToken := token[8:]
	if operation, exists := e.operations[baseToken]; exists {
		return operation(e)
	}
	return fmt.Errorf("undefined frozen operation: %s", token)
}

func (e *Evaluator) executeDefinition(def []string) error {
	for _, subToken := range def {
		if strings.HasPrefix(subToken, "default_") {
			err := e.executeFrozen(subToken)
			if err != nil {
				return err
			}
		} else {
			err := e.execute(subToken)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Evaluator) Process(row string) ([]int, error) {
	tokens := strings.Fields(row)
	for i := 0; i < len(tokens); i++ {
		token := strings.ToLower(tokens[i])
		if token == ":" {
			if i+1 >= len(tokens) {
				return nil, errors.New("invalid definition")
			}
			word := strings.ToLower(tokens[i+1])
			if _, err := strconv.Atoi(word); err == nil {
				return nil, fmt.Errorf("cannot redefine number: %s", word)
			}

			endIndex := -1
			for j := i + 2; j < len(tokens); j++ {
				if tokens[j] == ";" {
					endIndex = j
					break
				}
			}
			if endIndex == -1 {
				return nil, errors.New("missing ';' in definition")
			}
			expandedDefinition := tokens[i+2 : endIndex]
			for j := 0; j < len(expandedDefinition); j++ {
				expandedDefinition[j] = strings.ToLower(expandedDefinition[j])
			}
			expandedDefinitions := e.freezeOperations(expandedDefinition)
			e.keyWords[word] = expandedDefinitions
			i = endIndex
		} else {
			err := e.execute(token)
			if err != nil {
				return nil, err
			}
		}
	}
	return e.stack, nil
}

func (e *Evaluator) addOperation(*Evaluator) error {
	if len(e.stack) < 2 {
		return errors.New("stack underflow")
	}
	b, a := e.pop(), e.pop()
	e.stack = append(e.stack, a+b)
	return nil
}

func (e *Evaluator) subOperation(*Evaluator) error {
	if len(e.stack) < 2 {
		return errors.New("stack underflow")
	}
	b, a := e.pop(), e.pop()
	e.stack = append(e.stack, a-b)
	return nil
}

func (e *Evaluator) mulOperation(*Evaluator) error {
	if len(e.stack) < 2 {
		return errors.New("stack underflow")
	}
	b, a := e.pop(), e.pop()
	e.stack = append(e.stack, a*b)
	return nil
}

func (e *Evaluator) divOperation(*Evaluator) error {
	if len(e.stack) < 2 {
		return errors.New("stack underflow")
	}
	b, a := e.pop(), e.pop()
	if b == 0 {
		return errors.New("division by zero")
	}
	e.stack = append(e.stack, a/b)
	return nil
}

func (e *Evaluator) dupOperation(*Evaluator) error {
	if len(e.stack) < 1 {
		return errors.New("stack underflow")
	}
	e.stack = append(e.stack, e.stack[len(e.stack)-1])
	return nil
}

func (e *Evaluator) dropOperation(*Evaluator) error {
	if len(e.stack) < 1 {
		return errors.New("stack underflow")
	}
	e.pop()
	return nil
}

func (e *Evaluator) swapOperation(*Evaluator) error {
	if len(e.stack) < 2 {
		return errors.New("stack underflow")
	}
	e.stack[len(e.stack)-1], e.stack[len(e.stack)-2] = e.stack[len(e.stack)-2], e.stack[len(e.stack)-1]
	return nil
}

func (e *Evaluator) overOperation(*Evaluator) error {
	if len(e.stack) < 2 {
		return errors.New("stack underflow")
	}
	e.stack = append(e.stack, e.stack[len(e.stack)-2])
	return nil
}

func (e *Evaluator) pop() int {
	value := e.stack[len(e.stack)-1]
	e.stack = e.stack[:len(e.stack)-1]
	return value
}
