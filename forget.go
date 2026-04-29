package gomind

import "context"

// Forget removes a specific fact from the knowledge graph.
func (c *Client) Forget(ctx context.Context, subject, predicate, object string) error {
	return c.ForgetWithOptions(ctx, ForgetRequest{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
	})
}

// ForgetWithOptions removes a fact with full control over the request payload,
// including the optional Collection field.
func (c *Client) ForgetWithOptions(ctx context.Context, req ForgetRequest) error {
	req.Collection = c.resolveCollection(req.Collection)

	_, err := c.post(ctx, "/v1/forget", req)
	if err != nil {
		c.logger.Error("Gomind Forget failed",
			"error", err,
			"subject", req.Subject,
			"predicate", req.Predicate,
			"object", req.Object,
		)
		return err
	}

	c.logger.Info("Gomind Forget success",
		"subject", req.Subject,
		"predicate", req.Predicate,
		"object", req.Object,
	)

	return nil
}

// ForgetEntity removes all facts about a specific entity from the knowledge graph.
func (c *Client) ForgetEntity(ctx context.Context, entity string) error {
	return c.ForgetEntityWithOptions(ctx, ForgetEntityRequest{
		Entity: entity,
	})
}

// ForgetEntityWithOptions removes an entity with full control over the
// request payload, including the optional Collection field.
func (c *Client) ForgetEntityWithOptions(ctx context.Context, req ForgetEntityRequest) error {
	req.Collection = c.resolveCollection(req.Collection)

	_, err := c.post(ctx, "/v1/forget_entity", req)
	if err != nil {
		c.logger.Error("Gomind ForgetEntity failed", "error", err, "entity", req.Entity)
		return err
	}

	c.logger.Info("Gomind ForgetEntity success", "entity", req.Entity)
	return nil
}
