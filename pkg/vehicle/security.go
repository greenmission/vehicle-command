package vehicle

import (
	"context"
	"crypto/ecdh"

	"google.golang.org/protobuf/proto"

	"github.com/greenmission/vehicle-command/pkg/connector"
	"github.com/greenmission/vehicle-command/pkg/protocol"
	carserver "github.com/greenmission/vehicle-command/pkg/protocol/protobuf/carserver"
	"github.com/greenmission/vehicle-command/pkg/protocol/protobuf/vcsec"
)

func (v *Vehicle) SetValetMode(ctx context.Context, on bool, valetPassword string) error {
	return v.executeCarServerAction(ctx,
		&carserver.Action_VehicleAction{
			VehicleAction: &carserver.VehicleAction{
				VehicleActionMsg: &carserver.VehicleAction_VehicleControlSetValetModeAction{
					VehicleControlSetValetModeAction: &carserver.VehicleControlSetValetModeAction{
						On:       on,
						Password: valetPassword,
					},
				},
			},
		})
}

func (v *Vehicle) ResetValetPin(ctx context.Context) error {
	return v.executeCarServerAction(ctx,
		&carserver.Action_VehicleAction{
			VehicleAction: &carserver.VehicleAction{
				VehicleActionMsg: &carserver.VehicleAction_VehicleControlResetValetPinAction{
					VehicleControlResetValetPinAction: &carserver.VehicleControlResetValetPinAction{},
				},
			},
		})
}

// ResetPIN clears the saved PIN. You must disable PIN to drive before clearing the PIN. This allows
// setting a new PIN using SetPINToDrive.
func (v *Vehicle) ResetPIN(ctx context.Context) error {
	return v.executeCarServerAction(ctx,
		&carserver.Action_VehicleAction{
			VehicleAction: &carserver.VehicleAction{
				VehicleActionMsg: &carserver.VehicleAction_VehicleControlResetPinToDriveAction{
					VehicleControlResetPinToDriveAction: &carserver.VehicleControlResetPinToDriveAction{},
				},
			},
		})
}

func (v *Vehicle) ActivateSpeedLimit(ctx context.Context, speedLimitPin string) error {
	return v.executeCarServerAction(ctx,
		&carserver.Action_VehicleAction{
			VehicleAction: &carserver.VehicleAction{
				VehicleActionMsg: &carserver.VehicleAction_DrivingSpeedLimitAction{
					DrivingSpeedLimitAction: &carserver.DrivingSpeedLimitAction{
						Activate: true,
						Pin:      speedLimitPin,
					},
				},
			},
		})
}

func (v *Vehicle) DeactivateSpeedLimit(ctx context.Context, speedLimitPin string) error {
	return v.executeCarServerAction(ctx,
		&carserver.Action_VehicleAction{
			VehicleAction: &carserver.VehicleAction{
				VehicleActionMsg: &carserver.VehicleAction_DrivingSpeedLimitAction{
					DrivingSpeedLimitAction: &carserver.DrivingSpeedLimitAction{
						Activate: false,
						Pin:      speedLimitPin,
					},
				},
			},
		})
}

func (v *Vehicle) SpeedLimitSetLimitMPH(ctx context.Context, speedLimitMPH float64) error {
	return v.executeCarServerAction(ctx,
		&carserver.Action_VehicleAction{
			VehicleAction: &carserver.VehicleAction{
				VehicleActionMsg: &carserver.VehicleAction_DrivingSetSpeedLimitAction{
					DrivingSetSpeedLimitAction: &carserver.DrivingSetSpeedLimitAction{
						LimitMph: speedLimitMPH,
					},
				},
			},
		})
}

func (v *Vehicle) ClearSpeedLimitPIN(ctx context.Context, speedLimitPin string) error {
	return v.executeCarServerAction(ctx,
		&carserver.Action_VehicleAction{
			VehicleAction: &carserver.VehicleAction{
				VehicleActionMsg: &carserver.VehicleAction_DrivingClearSpeedLimitPinAction{
					DrivingClearSpeedLimitPinAction: &carserver.DrivingClearSpeedLimitPinAction{
						Pin: speedLimitPin,
					},
				},
			},
		})
}

func (v *Vehicle) SetSentryMode(ctx context.Context, state bool) error {
	return v.executeCarServerAction(ctx,
		&carserver.Action_VehicleAction{
			VehicleAction: &carserver.VehicleAction{
				VehicleActionMsg: &carserver.VehicleAction_VehicleControlSetSentryModeAction{
					VehicleControlSetSentryModeAction: &carserver.VehicleControlSetSentryModeAction{
						On: state,
					},
				},
			},
		})
}

// SetGuestMode enables or disables the vehicle's guest mode.
//
// We recommend users avoid this command unless they are managing a fleet of vehicles and understand
// the implications of enabling the mode. See official API documentation at
// https://developer.tesla.com/docs/fleet-api#guest_mode.
func (v *Vehicle) SetGuestMode(ctx context.Context, enabled bool) error {
	return v.executeCarServerAction(ctx,
		&carserver.Action_VehicleAction{
			VehicleAction: &carserver.VehicleAction{
				VehicleActionMsg: &carserver.VehicleAction_GuestModeAction{
					GuestModeAction: &carserver.VehicleState_GuestMode{
						GuestModeActive: enabled,
					},
				},
			},
		})
}

// SetPINToDrive controls whether the PIN to Drive feature is enabled or not. It is also used to set
// the PIN.
//
// Once a PIN is set, the vehicle remembers its value even when PIN to Drive is disabled and
// discards any new PIN provided using this method. To change an existing PIN, first call
// v.ResetPIN.
func (v *Vehicle) SetPINToDrive(ctx context.Context, enabled bool, pin string) error {
	if _, ok := v.conn.(connector.FleetAPIConnector); !ok {
		return protocol.ErrRequiresEncryption
	}

	return v.executeCarServerAction(ctx,
		&carserver.Action_VehicleAction{
			VehicleAction: &carserver.VehicleAction{
				VehicleActionMsg: &carserver.VehicleAction_VehicleControlSetPinToDriveAction{
					VehicleControlSetPinToDriveAction: &carserver.VehicleControlSetPinToDriveAction{
						On:       enabled,
						Password: pin,
					},
				},
			},
		})
}

func (v *Vehicle) TriggerHomelink(ctx context.Context, latitude float32, longitude float32) error {
	return v.executeCarServerAction(ctx,
		&carserver.Action_VehicleAction{
			VehicleAction: &carserver.VehicleAction{
				VehicleActionMsg: &carserver.VehicleAction_VehicleControlTriggerHomelinkAction{
					VehicleControlTriggerHomelinkAction: &carserver.VehicleControlTriggerHomelinkAction{
						Location: &carserver.LatLong{
							Latitude:  latitude,
							Longitude: longitude,
						},
					},
				},
			},
		})
}

// AddKey adds a public key to the vehicle's whitelist. If isOwner is true, the new key can
// authorize changes to vehicle access controls, such as adding/removing other keys.
func (v *Vehicle) AddKey(ctx context.Context, publicKey *ecdh.PublicKey, isOwner bool, formFactor vcsec.KeyFormFactor) error {
	if publicKey.Curve() != ecdh.P256() {
		return protocol.ErrInvalidPublicKey
	}
	payload := addKeyPayload(publicKey, isOwner, formFactor)
	encodedPayload, err := proto.Marshal(payload)
	if err != nil {
		return err
	}
	return v.executeWhitelistOperation(ctx, encodedPayload)
}

func (v *Vehicle) RemoveKey(ctx context.Context, publicKey *ecdh.PublicKey) error {
	if publicKey.Curve() != ecdh.P256() {
		return protocol.ErrInvalidPublicKey
	}
	payload := vcsec.UnsignedMessage{
		SubMessage: &vcsec.UnsignedMessage_WhitelistOperation{
			WhitelistOperation: &vcsec.WhitelistOperation{
				SubMessage: &vcsec.WhitelistOperation_RemovePublicKeyFromWhitelist{
					RemovePublicKeyFromWhitelist: &vcsec.PublicKey{
						PublicKeyRaw: publicKey.Bytes(),
					},
				},
			},
		},
	}
	encodedPayload, err := proto.Marshal(&payload)
	if err != nil {
		return err
	}
	return v.executeWhitelistOperation(ctx, encodedPayload)
}

func (v *Vehicle) KeySummary(ctx context.Context) (*vcsec.WhitelistInfo, error) {
	payload := vcsec.UnsignedMessage{
		SubMessage: &vcsec.UnsignedMessage_InformationRequest{
			InformationRequest: &vcsec.InformationRequest{
				InformationRequestType: vcsec.InformationRequestType_INFORMATION_REQUEST_TYPE_GET_WHITELIST_INFO,
			},
		},
	}
	encodedPayload, err := proto.Marshal(&payload)
	if err != nil {
		return nil, err
	}
	done := func(v *vcsec.FromVCSECMessage) (bool, error) { return true, nil }
	reply, err := v.getVCSECResult(ctx, encodedPayload, connector.AuthMethodNone, done)
	if err != nil {
		return nil, err
	}
	return reply.GetWhitelistInfo(), err
}

func (v *Vehicle) KeyInfoBySlot(ctx context.Context, slot uint32) (*vcsec.WhitelistEntryInfo, error) {
	payload := vcsec.UnsignedMessage{
		SubMessage: &vcsec.UnsignedMessage_InformationRequest{
			InformationRequest: &vcsec.InformationRequest{
				InformationRequestType: vcsec.InformationRequestType_INFORMATION_REQUEST_TYPE_GET_WHITELIST_ENTRY_INFO,
				Key: &vcsec.InformationRequest_Slot{
					Slot: slot,
				},
			},
		},
	}
	encodedPayload, err := proto.Marshal(&payload)
	if err != nil {
		return nil, err
	}
	done := func(v *vcsec.FromVCSECMessage) (bool, error) { return true, nil }
	reply, err := v.getVCSECResult(ctx, encodedPayload, connector.AuthMethodNone, done)
	if err != nil {
		return nil, err
	}
	return reply.GetWhitelistEntryInfo(), err
}

func (v *Vehicle) Lock(ctx context.Context) error {
	return v.executeRKEAction(ctx, vcsec.RKEAction_E_RKE_ACTION_LOCK)
}

func (v *Vehicle) Unlock(ctx context.Context) error {
	return v.executeRKEAction(ctx, vcsec.RKEAction_E_RKE_ACTION_UNLOCK)
}

// SendAddKeyRequest sends an add-key request to the vehicle over BLE. The user must approve the
// request by tapping their NFC card on the center console and then confirming their intent on the
// vehicle UI.
//
// If isOwner is true, the new key can authorize changes to vehicle access controls, such as
// adding/removing other keys.
//
// This function returns nil as soon as the request is transmitted. A nil return value does not
// guarantee the user has approved the request.
//
// Clients can check if publicKey has been enrolled and synchronized with the infotainment system by
// attempting to call v.SessionInfo with the domain argument set to
// [universal.Domain_DOMAIN_INFOTAINMENT].
func (v *Vehicle) SendAddKeyRequest(ctx context.Context, publicKey *ecdh.PublicKey, isOwner bool, formFactor vcsec.KeyFormFactor) error {
	if publicKey.Curve() != ecdh.P256() {
		return protocol.ErrInvalidPublicKey
	}
	if _, ok := v.conn.(connector.FleetAPIConnector); ok {
		return protocol.ErrRequiresBLE
	}
	encodedPayload, err := proto.Marshal(addKeyPayload(publicKey, isOwner, formFactor))
	if err != nil {
		return err
	}
	envelope := &vcsec.ToVCSECMessage{
		SignedMessage: &vcsec.SignedMessage{
			ProtobufMessageAsBytes: encodedPayload,
			SignatureType:          vcsec.SignatureType_SIGNATURE_TYPE_PRESENT_KEY,
		},
	}
	encodedEnvelope, err := proto.Marshal(envelope)
	if err != nil {
		return err
	}
	return v.conn.Send(ctx, encodedEnvelope)
}
