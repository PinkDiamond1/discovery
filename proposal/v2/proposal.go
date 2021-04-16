// Copyright (c) 2021 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package v2

import (
	"encoding/json"
	"math/big"
)

const Format = "service-proposal/v2"

type Proposal struct {
	Format         string         `json:"format"`
	Compatibility  int            `json:"compatibility"`
	ProviderID     string         `json:"provider_id"`
	ServiceType    string         `json:"service_type"`
	Location       Location       `json:"location"`
	Price          Price          `json:"price"`
	Contacts       []Contact      `json:"contacts"`
	AccessPolicies []AccessPolicy `json:"access_policies,omitempty"`
	Quality        float32        `json:"quality,omitempty"`
}

func NewProposal(providerID, serviceType string) *Proposal {
	return &Proposal{
		Format:      Format,
		ProviderID:  providerID,
		ServiceType: serviceType,
	}
}

func (p Proposal) Key() string {
	return p.ServiceType + ":" + p.ProviderID
}

func (p Proposal) MarshalBinary() (data []byte, err error) {
	marshal, err := json.Marshal(p)
	return marshal, err
}

func (p Proposal) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &p)
}

type IPType string

func (ipt IPType) IsResidential() bool {
	switch ipt {
	case "residential", "cellular":
		return true
	default:
		return false
	}
}

type Location struct {
	Continent string `json:"continent,omitempty"`
	Country   string `json:"country,omitempty"`
	City      string `json:"city,omitempty"`
	ASN       int    `json:"asn,omitempty"`
	ISP       string `json:"isp,omitempty"`
	IPType    IPType `json:"ip_type,omitempty"`
}

type Price struct {
	Currency Currency `json:"currency"`
	PerHour  *big.Int `json:"per_hour" swaggertype:"integer"`
	PerGiB   *big.Int `json:"per_gib" swaggertype:"integer"`
}

type Contact struct {
	Type       string           `json:"type"`
	Definition *json.RawMessage `json:"definition" swaggertype:"object"`
}

// AccessPolicy represents the access controls for proposal
type AccessPolicy struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}
