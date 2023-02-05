package rfidlib

const (
	/*
		Air protocol type
	*/
	RfidAplUnknownId   = 0
	RfidAplIso15693Id  = 1
	RfidAplIso14443aId = 2

	/*
		ISO15693 Tag type id
	*/
	RfidUnknownPiccId               = 0
	RfidIso15693PiccIcodeSliId      = 1
	RfidIso15693PiccTiHfiPlusId     = 2
	RfidIso15693PiccStId            = 3
	RfidIso15693PiccFujMb89r118cId  = 4
	RfidIso15693PiccStM24lr64Id     = 5
	RfidIso15693PiccStM24lr16eId    = 6
	RfidIso15693PiccIcodeSlixId     = 7
	RfidIso15693PiccTihfiStandardId = 8
	RfidIso15693PiccTihfiProId      = 9

	/*
		Describes ISO14443A Tag type id
	*/
	RfidIso14443aPiccNxpUltralightId = 1
	RfidIso14443aPiccNxpMifareS50Id  = 2
	RfidIso14443aPiccNxpMifareS70Id  = 3

	/*
		Inventory type
	*/
	AiTypeNew      = 1 // new antenna inventory  (reset RF power)
	AiTypeContinue = 2 // continue antenna inventory

	/*
		Move position
	*/
	RfidNoSeek    = 0 // No seeking
	RfidSeekFirst = 1 // Seek first
	RfidSeekNext  = 2 // Seek next
	RfidSeekLast  = 3 // Seek last

	/*
		usb enum information types
	*/
	HidEnumInfTypeSerialnum  = 1
	HidEnumInfTypeDriverpath = 2

	/*
		Get loaded reader driver option
	*/
	LoadedRdrdvrOptCatalog           = "CATALOG"
	LoadedRdrdvrOptName              = "NAME"
	LoadedRdrdvrOptId                = "ID"
	LoadedRdrdvrOptCommtypesupported = "COMM_TYPE_SUPPORTED"

	/*
		Reader driver type
	*/
	RdrdvrTypeReader = "Reader" // general reader
	RdrdvrTypeMtgate = "MTGate" // meeting gate
	RdrdvrTypeLsgate = "LSGate" // Library secure gate

	/*
		Open connection string
	*/

	ConnstrNameRdtype   = "RDType"
	ConnstrNameCommtype = "CommType"

	ConnstrNameCommtypeCom = "COM"
	ConnstrNameCommtypeUsb = "USB"
	ConnstrNameCommtypeNet = "NET"

	/*
		HID connection type params
	*/
	ConnstrNameHidaddrmode = "AddrMode"
	ConnstrNameHidsernum   = "SerNum"

	/*
		COM connection params
	*/
	ConnstrNameComname  = "COMName"
	ConnstrNameCombarud = "BaudRate"
	ConnstrNameComframe = "Frame"
	ConnstrNameBusaddr  = "BusAddr"

	/*
		TCP,UDP connection params
	*/
	ConnstrNameRemoteip   = "RemoteIP"
	ConnstrNameRemoteport = "RemotePort"
	ConnstrNameLocalip    = "LocalIP"

	/*
		supported comm type bits
	*/
	CommtypeComEn = 1
	CommtypeUsbEn = 2
	CommtypeNetEn = 4
	CommtypeBltEn = 8
)
