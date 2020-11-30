package data

import (
	"github.com/google/uuid"
)

// Account represents a bank account, its structure follows https://api-docs.form3.tech/api.html#organisation-accounts-resource.
// Fake account service does not implement "private_identification" and "relationships".
type Account struct {
	Type           RecordType `json:"type"`
	ID             uuid.UUID  `json:"id"`
	OrganisationID uuid.UUID  `json:"organisation_id"`
	Version        int        `json:"version"`
	Attributes     Attributes `json:"attributes"`
}
