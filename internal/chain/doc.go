// Package chain provides ordered stage-based validation of environment variable sets.
//
// A Chain is a sequence of named deployment stages (e.g. dev, staging, prod),
// each associated with an EnvSet that declares which environment variables are
// required or optional for that stage.
//
// Usage:
//
//	c := chain.New()
//	c.AddStage("dev",     envset.New([]string{"APP_HOST"}, nil))
//	c.AddStage("staging", envset.New([]string{"APP_HOST", "DB_URL"}, nil))
//	c.AddStage("prod",    envset.New([]string{"APP_HOST", "DB_URL", "SECRET_KEY"}, nil))
//
//	// Validate all stages:
//	if errs := c.Validate(); len(errs) > 0 {
//		for _, e := range errs {
//			fmt.Println(e)
//		}
//	}
//
//	// Validate only up to (and including) a specific stage:
//	if errs := c.ValidateUpTo("staging"); len(errs) > 0 {
//		fmt.Println(errs[0])
//	}
package chain
