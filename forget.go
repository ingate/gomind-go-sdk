package gomind

import "context"

// Forget removes a specific fact from the knowledge graph.
func (c *Client) Forget(ctx context.Context, subject, predicate, object string) error {
	req := ForgetRequest{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
	}

	_, err := c.post(ctx, "/v1/forget", req)
	if err != nil {
		c.logger.Error("Gomind Forget failed",
			"error", err,
			"subject", subject,
			"predicate", predicate,
			"object", object,
		)
		return err
	}

	c.logger.Info("Gomind Forget success",
		"subject", subject,
		"predicate", predicate,
		"object", object,
	)

	return nil
}

// ForgetEntity removes all facts about a specific entity from the knowledge graph.
func (c *Client) ForgetEntity(ctx context.Context, entity string) error {
	req := ForgetEntityRequest{
		Entity: entity,
	}

	_, err := c.post(ctx, "/v1/forget_entity", req)
	if err != nil {
		c.logger.Error("Gomind ForgetEntity failed", "error", err, "entity", entity)
		return err
	}

	c.logger.Info("Gomind ForgetEntity success", "entity", entity)
	return nil
}
