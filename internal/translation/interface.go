package translation

import "context"

// Service defines the translation service interface
type Service interface {
	// TranslateToRussian translates text to Russian
	TranslateToRussian(ctx context.Context, text string) (string, error)
	
	// TranslateBatch translates multiple texts to Russian
	TranslateBatch(ctx context.Context, texts []string) ([]string, error)
	
	// IsHealthy checks if the translation service is working
	IsHealthy(ctx context.Context) error
}