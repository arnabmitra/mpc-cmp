package sign

import (
	"crypto/ecdsa"
	"errors"

	"github.com/taurusgroup/cmp-ecdsa/pkg/math/curve"
	"github.com/taurusgroup/cmp-ecdsa/pkg/message"
	"github.com/taurusgroup/cmp-ecdsa/pkg/party"
	"github.com/taurusgroup/cmp-ecdsa/pkg/round"
	"github.com/taurusgroup/cmp-ecdsa/pkg/types"
)

type output struct {
	*round4
	// Signature wraps (R,S)
	Signature *Signature
}

// ProcessMessage implements round.Round
//
// - σⱼ != 0
func (r *output) ProcessMessage(from party.ID, content message.Content) error {
	body := content.(*SignOutput)
	partyJ := r.Parties[from]

	if body.SigmaShare.IsZero() {
		return ErrRoundOutputSigmaZero
	}
	partyJ.SigmaShare = body.SigmaShare
	return nil
}

// Finalize implements round.Round
//
// - compute σ = ∑ⱼ σⱼ
// - verify signature
func (r *output) Finalize(out chan<- *message.Message) (round.Round, error) {
	// compute σ = ∑ⱼ σⱼ
	S := curve.NewScalar()
	for _, partyJ := range r.Parties {
		S.Add(S, partyJ.SigmaShare)
	}

	r.Signature = &Signature{
		R: r.BigR,
		S: S,
	}

	RInt, SInt := r.Signature.ToRS()
	// Verify signature using Go's ECDSA lib
	if !ecdsa.Verify(r.PublicKey, r.Message, RInt, SInt) {
		return nil, ErrRoundOutputValidateSigFailedECDSA
	}
	pk := curve.FromPublicKey(r.PublicKey)
	if !r.Signature.Verify(pk, r.Message) {
		return nil, ErrRoundOutputValidateSigFailed
	}

	return nil, nil
}

func (r *output) Result() interface{} {
	// This could be used to handle pre-signatures
	if r.Signature != nil {
		return &Result{Signature: r.Signature}
	}
	return nil
}

func (r *output) MessageContent() message.Content {
	return &SignOutput{}
}

func (m *SignOutput) Validate() error {
	if m == nil {
		return errors.New("sign.round4: message is nil")
	}
	if m.SigmaShare == nil {
		return errors.New("sign.round4: message contains nil fields")
	}
	return nil
}

func (m *SignOutput) RoundNumber() types.RoundNumber {
	return 5
}
