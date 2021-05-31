package party

import (
	"errors"
	"fmt"

	"github.com/taurusgroup/cmp-ecdsa/pkg/math/curve"
	"github.com/taurusgroup/cmp-ecdsa/pkg/paillier"
	"github.com/taurusgroup/cmp-ecdsa/pkg/params"
	"github.com/taurusgroup/cmp-ecdsa/pkg/pedersen"
)

type Public struct {
	// ID of the party this data is associated with
	ID ID

	// SSID is the hash of the session ID this data is associated to
	// Should never be nil
	SSID []byte

	// ECDSA public key, may be nil if the keygen has not run yet
	ECDSA *curve.Point

	// Paillier public key, may be nil if the keygen has not run yet
	Paillier *paillier.PublicKey

	// Pedersen auxiliary parameters, may be nil if the keygen has not run yet
	Pedersen *pedersen.Parameters
}

func (p Public) preKeygen() bool {
	return p.ECDSA == nil && p.Paillier == nil && p.Pedersen == nil
}

// KeygenDone returns true if all fields resulting from a keygen are non nil
func (p Public) KeygenDone() bool {
	return p.ECDSA != nil && p.Paillier != nil && p.Pedersen != nil
}

func (p *Public) Clone() *Public {
	p2 := &Public{
		ID:   p.ID,
		SSID: append([]byte{}, p.SSID...),
	}

	if p.Paillier != nil {
		p2.Paillier = paillier.NewPublicKey(p.Paillier.N)
	}
	if p.ECDSA != nil {
		p2.ECDSA = curve.NewIdentityPoint().Set(p.ECDSA)
	}
	if p.Pedersen != nil {
		p2.Pedersen = p.Pedersen.Clone()
	}
	return p2
}

// Validate returns an error if Public is invalid. Otherwise return nil.
func (p *Public) Validate() error {
	if p.ID == "" {
		return errors.New("party.Public: ID cannot be empty")
	}

	if p.preKeygen() {
		return nil
	}

	// check SSID length
	if len(p.SSID) != params.HashBytes {
		return errors.New("party.Public: SSID has wrong length")
	}

	// nil checks
	if p.ECDSA == nil {
		return errors.New("party.Public: ECDSA public share cannot be nil")
	}
	if p.Paillier == nil {
		return errors.New("party.Public: Paillier public key cannot be nil")
	}
	if p.Pedersen == nil {
		return errors.New("party.Public: Pedersen parameters cannot be nil")
	}

	// ECDSA is not identity
	if p.ECDSA.IsIdentity() {
		return errors.New("party.Public: ECDSA public key is identity")
	}

	// Paillier check
	if err := p.Paillier.Validate(); err != nil {
		return fmt.Errorf("party.Public: %w", err)
	}

	// Pedersen check
	if err := p.Pedersen.Validate(); err != nil {
		return fmt.Errorf("party.Public: %w", err)
	}

	// Both N's are the same
	if p.Paillier.N.Cmp(p.Pedersen.N) != 0 {
		return errors.New("party.Public: Pedersen and Paillier should share the same N")
	}

	return nil
}
