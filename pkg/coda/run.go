package coda

import (
	"fmt"
	"time"
)

func (c *Coda) Run() error {
	start := time.Now()
	c.debug(fmt.Sprintf("found %d operations", len(c.Operations)))
	lastOp, ops, err := c.runOperations(c.Operations)
	c.debug(fmt.Sprintf("executed %d operations after %s", ops, time.Since(start)))
	if err != nil {
		return fmt.Errorf("failed to execute node '%s': %s", lastOp.Action, err)
	}
	return nil
}

func (c *Coda) runOperations(ops []Operation) (Operation, int64, error) {
	var count int64 = 0
	for _, op := range ops {
		count++
		if err := c.executeOperation(op); err != nil {
			if len(op.OnFail) > 0 {
				n, c, err := c.runOperations(op.OnFail)
				count += c
				if err != nil {
					return n, count, err
				}
			} else {
				return ops[len(ops)-1], count, err
			}
		}
	}
	return ops[len(ops)-1], count, nil
}

func (c *Coda) executeOperation(op Operation) error {
	if action, ok := operations[op.Action]; !ok {
		return fmt.Errorf("unknown action: %s", op.Action)
	} else {
		start := time.Now()
		defer func() {
			c.debug(fmt.Sprintf("executed operation '%s' after %s", action.Name, time.Since(start)))
		}()
		p, err := c.resolveVariables(op.Params)
		if err != nil {
			return fmt.Errorf("failed to resolve variables: %v", err)
		}
		op.Params = p
		result, err := action.Fn(c, op.Params)
		if err != nil {
			return err
		}

		if op.Store != "" {
			if len(result) != 0 {
				if c.Coda.Strict {
					if _, ok := c.Store[op.Store]; !ok {
						return fmt.Errorf("store '%s' does not exist", op.Store)
					}
				}
				c.Store[op.Store] = result
			}
		}
	}
	return nil
}
