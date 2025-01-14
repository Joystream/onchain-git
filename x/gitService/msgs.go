package gitService

import (
	"encoding/json"

	"gopkg.in/src-d/go-git.v4/plumbing/protocol/packp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// UpdateReferenceCommand encodes how to update a reference
type UpdateReferenceCommand struct {
	Name plumbing.ReferenceName
	Old  plumbing.Hash
	New  plumbing.Hash
}

// UpdateReferenceAction represents the type of reference update action
type UpdateReferenceAction string

const (
	// CreateAction says to create reference
	CreateAction UpdateReferenceAction = "create"
	// UpdateAction says to update reference
	UpdateAction = "update"
	// DeleteAction says to delete reference
	DeleteAction = "delete"
	// InvalidAction means reference is invalid
	InvalidAction = "invalid"
)

// Action says which kind of reference update should be undertaken
func (c *UpdateReferenceCommand) Action() UpdateReferenceAction {
	if c.Old == plumbing.ZeroHash && c.New == plumbing.ZeroHash {
		return InvalidAction
	}

	if c.Old == plumbing.ZeroHash {
		return CreateAction
	}

	if c.New == plumbing.ZeroHash {
		return DeleteAction
	}

	return UpdateAction
}

func (c *UpdateReferenceCommand) validate() error {
	if c.Action() == InvalidAction {
		return errors.Errorf("Malformed command")
	}

	return nil
}

// MsgUpdateReferences defines the UpdateReferences message
type MsgUpdateReferences struct {
	URI      string
	Author   sdk.AccAddress
	Commands []*UpdateReferenceCommand
	Shallow  *plumbing.Hash
	Packfile []byte
}

// NewMsgUpdateReferences is the constructor function for MsgUpdateReferences
func NewMsgUpdateReferences(uri string, req *packp.ReferenceUpdateRequest,
	packfile []byte, author sdk.AccAddress) (*MsgUpdateReferences, sdk.Error) {
	cmds := make([]*UpdateReferenceCommand, 0, len(req.Commands))
	for _, cmd := range req.Commands {
		cmds = append(cmds, &UpdateReferenceCommand{
			Name: cmd.Name,
			Old:  cmd.Old,
			New:  cmd.New,
		})
	}
	msg := &MsgUpdateReferences{
		URI:      uri,
		Commands: cmds,
		Packfile: packfile,
		Shallow:  req.Shallow,
		Author:   author,
	}

	return msg, msg.ValidateBasic()
}

// Route implements Msg.
func (msg MsgUpdateReferences) Route() string { return "gitService" }

// Type implements Msg.
func (msg MsgUpdateReferences) Type() string { return "push" }

// ValidateBasic Implements Msg.
func (msg MsgUpdateReferences) ValidateBasic() sdk.Error {
	if msg.Author.Empty() {
		log.Debug().Msgf("MsgUpdateReferences author empty")
		return sdk.ErrInvalidAddress(msg.Author.String())
	}
	if len(msg.URI) == 0 {
		log.Debug().Msgf("MsgUpdateReferences URI empty")
		return sdk.ErrUnknownRequest("URI cannot be empty")
	}
	if len(msg.Commands) == 0 {
		log.Debug().Msgf("MsgUpdateReferences commands empty")
		return sdk.ErrUnknownRequest("Commands cannot be empty")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUpdateReferences) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgUpdateReferences) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Author}
}

// MsgRemoveRepository defines the DeleteRepository message
type MsgRemoveRepository struct {
	URI    string
	Author sdk.AccAddress
}

// NewMsgRemoveRepository is the constructor function for MsgRemoveRepository
func NewMsgRemoveRepository(uri string, author sdk.AccAddress) (*MsgRemoveRepository, sdk.Error) {
	msg := &MsgRemoveRepository{
		URI:    uri,
		Author: author,
	}

	return msg, msg.ValidateBasic()
}

// Route implements Msg.
func (msg MsgRemoveRepository) Route() string { return "gitService" }

// Type implements Msg.
func (msg MsgRemoveRepository) Type() string { return "removeRepository" }

// ValidateBasic Implements Msg.
func (msg MsgRemoveRepository) ValidateBasic() sdk.Error {
	if msg.Author.Empty() {
		log.Debug().Msgf("MsgRemoveRepository author empty")
		return sdk.ErrInvalidAddress(msg.Author.String())
	}
	if len(msg.URI) == 0 {
		log.Debug().Msgf("MsgRemoveRepository URI empty")
		return sdk.ErrUnknownRequest("URI cannot be empty")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRemoveRepository) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgRemoveRepository) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Author}
}
