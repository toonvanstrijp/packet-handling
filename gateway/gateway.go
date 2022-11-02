// Copyright 2022 Stichting ThingsIX Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package gateway

import (
	"crypto/ecdsa"
	"crypto/sha256"

	"github.com/ThingsIXFoundation/packet-handling/utils"
	"github.com/brocaar/lorawan"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Gateway struct {
	LocalGatewayID           lorawan.EUI64
	NetworkGatewayID         lorawan.EUI64
	PrivateKey               *ecdsa.PrivateKey
	PublicKey                *ecdsa.PublicKey
	PublicKeyBytes           []byte
	CompressedPublicKeyBytes []byte
	Owner                    common.Address
}

// ID is the identifier as which the gateway is registered in the gateway registry.
func (gw Gateway) ID() [32]byte {
	var id [32]byte
	copy(id[:], gw.CompressedPublicKeyBytes)
	return id
}

// CompressedPubKeyBytes returns the compressed public key including 0x02 prefix
func (gw Gateway) CompressedPubKeyBytes() []byte {
	return append([]byte{0x2}, gw.CompressedPublicKeyBytes...)
}

func NewGateway(localGatewayID lorawan.EUI64, priv *ecdsa.PrivateKey) (*Gateway, error) {
	return &Gateway{
		LocalGatewayID:           localGatewayID,
		NetworkGatewayID:         CalculateNetworkGatewayID(priv),
		PrivateKey:               priv,
		PublicKey:                &priv.PublicKey,
		CompressedPublicKeyBytes: CalculateCompressedPublicKeyBytes(&priv.PublicKey),
		PublicKeyBytes:           CalculatePublicKeyBytes(&priv.PublicKey),
		Owner:                    common.Address{}, // TODO
	}, nil
}

func GenerateNewGateway(localID lorawan.EUI64) (*Gateway, error) {
	priv, err := GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	return &Gateway{
		LocalGatewayID:           localID,
		NetworkGatewayID:         CalculateNetworkGatewayID(priv),
		PrivateKey:               priv,
		PublicKey:                &priv.PublicKey,
		CompressedPublicKeyBytes: CalculateCompressedPublicKeyBytes(&priv.PublicKey),
		Owner:                    common.Address{}, // TODO
	}, nil
}

func CalculateNetworkGatewayID(priv *ecdsa.PrivateKey) lorawan.EUI64 {
	pub := priv.PublicKey
	pubBytes := CalculateCompressedPublicKeyBytes(&pub)
	h := sha256.Sum256(pubBytes)

	gatewayID, _ := utils.BytesToGatewayID(h[0:8])

	return gatewayID
}

func CalculateCompressedPublicKeyBytes(pub *ecdsa.PublicKey) []byte {
	return crypto.CompressPubkey(pub)[1:]
}

func CalculatePublicKeyBytes(pub *ecdsa.PublicKey) []byte {
	return crypto.FromECDSAPub(pub)
}

func (gw *Gateway) Address() common.Address {
	return crypto.PubkeyToAddress(*gw.PublicKey)
}