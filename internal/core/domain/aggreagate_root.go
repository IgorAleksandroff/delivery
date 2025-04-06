package domain

type AggregateRoot interface {
	GetDomainEvents() []Event
	ClearDomainEvents()
}
