package coda

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func (c *Coda) run() error {
	start := time.Now()
	c.debug(fmt.Sprintf("found %d operations", len(c.Operations)))
	lastOp, ops, err := c.runOperations(c.Operations)
	defer func() {
		since := time.Since(start)
		c.debug(fmt.Sprintf("executed %d operations after %s", ops, since))
		c.Stats.CodaRuntimeTotalMs += float64(since.Milliseconds())
	}()

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
			c.Stats.OperationsFailedTotal++
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
			c.Stats.OperationsSuccessfulTotal++
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
			c.Stats.OperationsBlacklistedTotal++
			return fmt.Errorf("category of operation '%s' is disabled (%s)", op.Action, action.Category)
		}
		start := time.Now()
		defer func() {
			since := time.Since(start)
			c.Stats.OperationsTotal++
			c.Stats.OperationsRuntimeTotalMs += float64(since.Milliseconds())
			c.debug(fmt.Sprintf("executed operation '%s' after %s", action.Name, since))
		}()
		p, err := c.resolveVariables(op.Params)
		if err != nil {
			c.Stats.VariablesFailedTotal++
			return fmt.Errorf("failed to resolve variables: %v", err)
		}
		op.Params = p

		execWithLock := func() error {
			c.storeMutex.Lock()
			defer c.storeMutex.Unlock()

			result, err := action.Fn(c, op.Params)
			if err != nil {
				return err
			}

			if op.Store != "" && len(result) != 0 {
				// check if the result should be stored in a JSON path
				if strings.Contains(op.Store, ".") {
					err := c.storeNestedJSONValue(op.Store, result)
					if err != nil {
						return fmt.Errorf("failed to store nested value: %v", err)
					}
				} else {
					c.Store[op.Store] = result
				}
			}
			return nil
		}

		if op.Async {
			go execWithLock()
			return nil
		}
		return execWithLock()
	}
}

func (c *Coda) storeNestedJSONValue(path string, value json.RawMessage) error {
	parts := strings.Split(path, ".")
	rootKey := parts[0]

	// Create or update the nested structure
	var rootObj map[string]interface{}

	// If the root key exists, start with that data
	if existingData, exists := c.Store[rootKey]; exists {
		if err := json.Unmarshal(existingData, &rootObj); err != nil {
			rootObj = make(map[string]interface{})
		}
	} else {
		rootObj = make(map[string]interface{})
	}

	// Navigate to the right location and set the value
	var valueObj interface{}
	if err := json.Unmarshal(value, &valueObj); err != nil {
		// If we can't unmarshal as object, use the raw value
		valueObj = string(value)
	}

	// Navigate and build the nested structure
	current := rootObj
	for i := 1; i < len(parts)-1; i++ {
		key := parts[i]

		// Check if this path exists
		if _, exists := current[key]; !exists {
			current[key] = make(map[string]interface{})
		}

		// If it's not a map, convert it
		if nextMap, ok := current[key].(map[string]interface{}); ok {
			current = nextMap
		} else {
			// Replace with a new map
			newMap := make(map[string]interface{})
			current[key] = newMap
			current = newMap
		}
	}

	// Set the final value
	if len(parts) > 1 {
		lastKey := parts[len(parts)-1]
		current[lastKey] = valueObj
	}

	// Convert back to JSON
	updatedData, err := json.Marshal(rootObj)
	if err != nil {
		return err
	}

	// Store the updated structure
	c.Store[rootKey] = json.RawMessage(updatedData)
	return nil
}

func (c *Coda) isBlacklisted(category OperationCategory) bool {
	if len(c.blacklist) == 0 {
		return false
	}
	for _, bl := range c.blacklist {
		if bl == category {
			return true
		}
	}
	return false
}
