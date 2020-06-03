package commercio

import (
	"encoding/json"

	commerciomint "github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	"github.com/commercionetwork/commercionetwork/x/docs"
	id "github.com/commercionetwork/commercionetwork/x/id/types"
	"github.com/commercionetwork/commercionetwork/x/memberships"
	"github.com/commercionetwork/commercionetwork/x/vbr"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// messageEnclosure encloses a Cosmos message into its REST-accepted enclosure.
type messageEnclosure struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

//
// Commercio.network type exports
//

type (
	// Standard Cosmos-sdk messages
	MsgSend bank.MsgSend

	// x/docs messages
	Document                            docs.Document
	DocumentMetadata                    docs.DocumentMetadata
	MetadataSchema                      docs.MetadataSchema
	DocumentMetadataSchema              docs.DocumentMetadataSchema
	DocumentChecksum                    docs.DocumentChecksum
	DocumentReceipt                     docs.DocumentReceipt
	MsgShareDocument                    docs.MsgShareDocument
	MsgSendDocumentReceipt              docs.MsgSendDocumentReceipt
	MsgAddSupportedMetadataSchema       docs.MsgAddSupportedMetadataSchema
	MsgAddTrustedMetadataSchemaProposer docs.MsgAddTrustedMetadataSchemaProposer

	// x/id messages
	DidDocument          id.DidDocument
	PubKeys              id.PubKeys
	PubKey               id.PubKey
	Proof                id.Proof
	Services             id.Services
	MsgSetIdentity       id.MsgSetIdentity
	MsgRequestDidPowerUp id.MsgRequestDidPowerUp

	// x/memberships messages
	MsgInviteUser               memberships.MsgInviteUser
	MsgDepositIntoLiquidityPool memberships.MsgDepositIntoLiquidityPool
	MsgBuyMembership            memberships.MsgBuyMembership

	// x/vbr messages
	MsgIncrementsBlockRewardsPool vbr.MsgIncrementsBlockRewardsPool

	// x/commerciomint messages
	MsgOpenCdp  = commerciomint.MsgOpenCdp
	MsgCloseCdp = commerciomint.MsgCloseCdp
)

// Membership types definition
var (
	MembershipTypeBronze = "bronze"
	MembershipTypeSilver = "silver"
	MembershipTypeGold   = "gold"
	MembershipTypeBlack  = "black"
)
