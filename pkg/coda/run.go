package coda

import (
	"fmt"
	"strings"
	"time"

	"github.com/yosev/coda/pkg/metrics"
)

func (c *Coda) Run() error {
	start := time.Now()
	c.debug(fmt.Sprintf("found %d operations", len(c.Operations)))
	lastOp, ops, err := c.runOperations(c.Operations)
	defer func() {
		since := time.Since(start)
		c.debug(fmt.Sprintf("executed %d operations after %s", ops, since))
		metrics.Inc("coda_total")
		metrics.IncValue("coda_runtime_total", float64(since.Milliseconds()))
	}()

	if err != nil {
		metrics.Inc("coda_failed_total")
		return fmt.Errorf("failed to execute node '%s': %s", lastOp.Action, err)
	}

	metrics.Inc("coda_successful_total")
	return nil
}

func (c *Coda) runOperations(ops []Operation) (Operation, int64, error) {
	var count int64 = 0
	for _, op := range ops {
		count++
		if err := c.executeOperation(op); err != nil {
			metrics.Inc("operations_failed_total")
			// TODO check why wrong operation action/name is being returned
			// failed to run coda: failed to execute node 'io.stdout': category of operation 'file.size' is disabled (File)

			// check for blacklisted error and if so, do not apply onFail
			if strings.Contains(err.Error(), "is blacklisted") {
				return op, count, err
			}

			// check for on-fail operations
			if len(op.OnFail) > 0 {
				n, c, err := c.runOperations(op.OnFail)
				count += c
				if err != nil {
					return n, count, err
				}
			} else {
				return ops[len(ops)-1], count, err
			}
		} else {
			metrics.Inc("operations_successful_total")
		}
	}
	if len(ops) > 0 {
		return ops[len(ops)-1], count, nil
	}
	return Operation{}, count, nil
}

func (c *Coda) executeOperation(op Operation) error {
	if action, ok := operations[op.Action]; !ok {
		return fmt.Errorf("unknown action: %s", op.Action)
	} else {
		if c.isBlacklisted(action.Category) {
			metrics.Inc("operations_blacklisted_total")
			return fmt.Errorf("category of operation '%s' is disabled (%s)", op.Action, action.Category)
		}
		start := time.Now()
		defer func() {
			since := time.Since(start)
			metrics.Inc("operations_total")
			metrics.IncValue("operations_runtime_total", float64(since.Milliseconds()))
			c.debug(fmt.Sprintf("executed operation '%s' after %s", action.Name, since))
		}()
		p, err := c.resolveVariables(op.Params)
		if err != nil {
			metrics.Inc("variables_failed_total")
			return fmt.Errorf("failed to resolve variables: %v", err)
		}
		op.Params = p
		result, err := action.Fn(c, op.Params)
		if err != nil {
			return err
		}

		if op.Store != "" {
			if len(result) != 0 {
				c.Store[op.Store] = result
			}
		}
	}
	return nil
}

func (c *Coda) isBlacklisted(category OperationCategory) bool {
	if len(c.Blacklist) == 0 {
		return false
	}
	for _, bl := range c.Blacklist {
		if bl == category {
			return true
		}
	}
	return false
}
