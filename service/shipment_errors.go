package service

import "errors"

var (
	ErrMarkShippedNotAllowed = errors.New("cannot mark shipped for this auction")
	ErrInvalidShipmentCarrier = errors.New("invalid carrier")
	ErrInvalidTrackingNumber  = errors.New("invalid tracking number")
	ErrShipmentNotDelivered     = errors.New("shipment not delivered yet")
	ErrShipmentAccessDenied     = errors.New("shipment access denied")
	ErrShipmentNotRegistered    = errors.New("shipment not registered")
)
