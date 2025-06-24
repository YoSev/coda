package coda

import "fmt"

func (c *Coda) findEntrypoint() (string, error) {
	var k string
	var found = false
	for key, op := range c.Operations {
		if op.Entrypoint {
			if found {
				return "", fmt.Errorf("multiple entrypoints found: %s and %s", k, key)
			}
			k = key
			found = true
		}
	}
	if !found {
		return "", fmt.Errorf("missing entrypoint")
	}
	return k, nil
}

func (c *Coda) validateLinks() error {
	for uid, op := range c.Operations {
		if err := c.isValidLink(uid, op.OnSuccess); err != nil {
			return err
		}

		if err := c.isValidLink(uid, op.OnFail); err != nil {
			return err
		}
	}
	return nil
}

func (c *Coda) isValidLink(sourceUid string, targetUid string) error {
	if targetUid == "" {
		return nil // No link to validate
	}

	if sourceUid == targetUid {
		return fmt.Errorf("self-links are not allowed: %s -> %s", sourceUid, targetUid)
	}

	_, exists := c.Operations[targetUid]
	if !exists {
		return fmt.Errorf("target operation %s does not exist for link: %s -> %s", targetUid, sourceUid, targetUid)
	}
	return nil
}
