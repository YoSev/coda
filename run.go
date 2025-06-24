package coda

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func (c *Coda) run() error {
	start := time.Now()
	defer func() {
		since := time.Since(start)
		c.Stats.CodaRuntimeTotalMs += float64(since.Milliseconds())
	}()

	if startUid, err := c.findEntrypoint(); err != nil {
		return err
	} else {
		if err := c.validateLinks(); err != nil {
			return fmt.Errorf("failed to validate links: %s", err)
		} else {
			c.debug(fmt.Sprintf("found %d operations", len(c.Operations)))
			lastUid, err := c.runOperations(startUid)
			if err != nil {
				return fmt.Errorf("failed to execute node '%s': %s", lastUid, err)
			}
		}

	}

	return nil
}

func (c *Coda) runOperations(uid string) (string, error) {
	for uid != "" {
		op, ok := c.Operations[uid]
		if !ok {
			return uid, fmt.Errorf("operation with UID %s not found", uid)
		}

		err := c.executeOperation(op)
		if err != nil {
			c.Stats.OperationsFailedTotal++
			if op.OnFail == "" {
				return uid, err
			}
			uid = op.OnFail
		} else {
			c.Stats.OperationsSuccessfulTotal++
			if op.OnSuccess == "" {
				return uid, nil
			}
			uid = op.OnSuccess
		}
	}

	return uid, nil
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
			result, err := action.Fn(c, op.Params)
			if err != nil {
				return err
			}

			// delay locking to make this routine non-blocking during the actual execution
			c.mutex.Lock()
			defer c.mutex.Unlock()

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
