package payment

import "context"

type Provider interface {
	Initialize(
		ctx context.Context,
		payment *Payment,
		email string,
	) (*PaymentInitialization, error)

	Verify(
		ctx context.Context,
		providerRef string,
	) (*PaymentVerification, error)
}

type ProviderRegistry struct {
	providers map[string]Provider
}

func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers: make(map[string]Provider),
	}
}

func (r *ProviderRegistry) Register(
	name string,
	provider Provider,
) {
	r.providers[name] = provider
}

func (r *ProviderRegistry) Get(name string) (Provider, bool) {
	provider, ok := r.providers[name]
	return provider, ok
}
