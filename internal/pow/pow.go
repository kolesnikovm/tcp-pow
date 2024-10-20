package pow

import (
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"
)

type PowShieldFatory struct{}

func NewPowShieldFactory() *PowShieldFatory {
	return &PowShieldFatory{}
}

type PowShield struct {
	difficulty int
	target     *big.Int
	maxNonce   int64

	challenge []byte
}

func (f *PowShieldFatory) NewPowShield(difficulty int) *PowShield {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	return &PowShield{
		difficulty: difficulty,
		target:     target,
		maxNonce:   math.MaxInt64,
	}
}

func (p *PowShield) GetDifficulty() int {
	return p.difficulty
}

func (p *PowShield) SetChallenge(challenge []byte) {
	p.challenge = challenge
}

func (p *PowShield) VerifySolution(nonce int64) bool {
	var hashInt big.Int

	data := p.prepareData(nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(p.target) == -1

	return isValid
}

func (p *PowShield) GetSolution() int64 {
	var hashInt big.Int
	var hash [32]byte
	nonce := int64(0)

	for nonce < p.maxNonce {
		data := p.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(p.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce
}

func (p *PowShield) prepareData(nonce int64) []byte {
	nonceBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(nonceBytes, uint32(nonce))

	data := append(p.challenge[:], nonceBytes...)

	return data
}
