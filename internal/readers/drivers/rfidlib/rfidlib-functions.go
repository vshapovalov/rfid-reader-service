package rfidlib

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

var (
	rfidlib_reader = syscall.NewLazyDLL("drivers\\rfidlib\\rfidlib_reader.dll")
	rdr_Open       = rfidlib_reader.NewProc("RDR_Open")
	rdr_Close      = rfidlib_reader.NewProc("RDR_Close")
	//rdr_GetLibVersion              = rfidlib_reader.NewProc("RDR_GetLibVersion")
	rdr_GetAntennaInterfaceCount = rfidlib_reader.NewProc("RDR_GetAntennaInterfaceCount")
	//rdr_SetCommuImmeTimeout        = rfidlib_reader.NewProc("RDR_SetCommuImmeTimeout")
	rdr_TagInventory               = rfidlib_reader.NewProc("RDR_TagInventory")
	rdr_GetTagDataReportCount      = rfidlib_reader.NewProc("RDR_GetTagDataReportCount")
	rdr_GetTagDataReport           = rfidlib_reader.NewProc("RDR_GetTagDataReport")
	dnode_Destroy                  = rfidlib_reader.NewProc("DNODE_Destroy")
	rdr_ResetCommuImmeTimeout      = rfidlib_reader.NewProc("RDR_ResetCommuImmeTimeout")
	hid_GetEnumItem                = rfidlib_reader.NewProc("HID_GetEnumItem")
	hid_Enum                       = rfidlib_reader.NewProc("HID_Enum")
	rdr_LoadReaderDrivers          = rfidlib_reader.NewProc("RDR_LoadReaderDrivers")
	rdr_GetLoadedReaderDriverCount = rfidlib_reader.NewProc("RDR_GetLoadedReaderDriverCount")
	rdr_GetLoadedReaderDriverOpt   = rfidlib_reader.NewProc("RDR_GetLoadedReaderDriverOpt")
	comPort_Enum                   = rfidlib_reader.NewProc("COMPort_Enum")
	comPort_GetEnumItem            = rfidlib_reader.NewProc("COMPort_GetEnumItem")
	rdr_TagDisconnect              = rfidlib_reader.NewProc("RDR_TagDisconnect")

	rfid_aip_iso15693            = syscall.NewLazyDLL("drivers\\rfidlib\\rfidlib_aip_iso15693.dll")
	iso15693_ParseTagDataReport  = rfid_aip_iso15693.NewProc("ISO15693_ParseTagDataReport")
	rfid_aip_iso14443a           = syscall.NewLazyDLL("drivers\\rfidlib\\rfidlib_aip_iso14443A.dll")
	iso14443A_ParseTagDataReport = rfid_aip_iso14443a.NewProc("ISO14443A_ParseTagDataReport")
)

func stringToPtr(s string) uintptr {
	return uintptr(unsafe.Pointer(&[]byte(s + "\x00")[0]))
}

func makeSliceToPtr(size int) (uintptr, []byte) {
	slc := make([]byte, size, size)
	return uintptr(unsafe.Pointer(&slc[0])), slc
}

func fillSliceWithNullString(val []byte) {
	for i := 0; i < len(val); i++ {
		val[i] = 0
	}
}

func loadReaderDrivers(path string) error {
	path = path + "\x00"
	r1, _, err := rdr_LoadReaderDrivers.Call(stringToPtr(path))
	if r1 != 0 {
		return fmt.Errorf("cannot load drivers: %d, %w", r1, err)
	}
	return nil
}

func getLoadedReaderDriverCount() int {
	r1, _, _ := rdr_GetLoadedReaderDriverCount.Call()
	return int(r1)
}

func getLoadedReaderDriverOpt(driverNum int, opt string) (string, error) {
	var strBufSize int
	strBufPtr, strBuf := makeSliceToPtr(64)

	strBufSize = len(strBuf)
	sizePtr := uintptr(unsafe.Pointer(&strBufSize))
	r1, _, err := rdr_GetLoadedReaderDriverOpt.Call(uintptr(driverNum), stringToPtr(opt), strBufPtr, sizePtr)
	if r1 != 0 {
		return "", fmt.Errorf("cannot get [%s] of driver [%d]: %d, %w", opt, driverNum, r1, err)
	}
	return strings.TrimRight(string(strBuf), "\x00"), nil
}

func comPortEnum() int {
	r1, _, _ := comPort_Enum.Call()
	return int(r1)
}

func comPortGetEnumItem(comPortItemNum int) (string, error) {
	var strBufSize int

	strBufPtr, strBuf := makeSliceToPtr(32)
	strBufSize = len(strBuf)
	sizePtr := uintptr(unsafe.Pointer(&strBufSize))
	r1, _, err := comPort_GetEnumItem.Call(uintptr(comPortItemNum), strBufPtr, sizePtr)
	if r1 != 0 {
		return "", fmt.Errorf("cannot get COM port item [%d] info: %d, %w", comPortItemNum, r1, err)
	}
	return strings.TrimRight(string(strBuf), "\x00"), nil
}

func hidEnum(driverName string) int {
	targetDriverName := driverName
	r1, _, _ := hid_Enum.Call(stringToPtr(targetDriverName))
	return int(r1)
}

func rdrOpen(connString string) (int, error) {
	var reader int
	r1, _, err := rdr_Open.Call(stringToPtr(connString), uintptr(unsafe.Pointer(&reader)))
	if r1 != 0 {
		return 0, fmt.Errorf("cannot open reader: %d, %w", r1, err)
	}
	return reader, nil
}

func rdrGetAntennaInterfaceCount(readerHandler int) int {
	r1, _, _ := rdr_GetAntennaInterfaceCount.Call(uintptr(readerHandler))
	return int(r1)
}

func rdrClose(readerHandler int) error {
	r1, _, err := rdr_Close.Call(uintptr(readerHandler))
	if r1 != 0 {
		return fmt.Errorf("cannot close reader: %d, %w", r1, err)
	}
	return nil
}

func hidGetEnumItemOpt(itemNum int, opt byte) (string, error) {
	var strBufSize int

	strBufPtr, strBuf := makeSliceToPtr(64)

	strBufSize = len(strBuf)
	sizePtr := uintptr(unsafe.Pointer(&strBufSize))
	r1, _, err := hid_GetEnumItem.Call(uintptr(itemNum), uintptr(opt), strBufPtr, sizePtr)
	if r1 != 0 {
		return "", fmt.Errorf("cannot get [%d] of hid item [%d]: %d, %w", opt, itemNum, r1, err)
	}
	return strings.TrimRight(string(strBuf), "\x00"), nil
}

func tagInventory(readerHandler int, aiType byte, antennasIds []byte, dataNodeHandler int) error {
	r1, _, err := rdr_TagInventory.Call(
		uintptr(readerHandler),
		uintptr(aiType),
		uintptr(len(antennasIds)),
		uintptr(unsafe.Pointer(&antennasIds[0])),
		uintptr(dataNodeHandler),
	)
	if int(r1) == 0 || int(r1) == -21 {
		return nil
	}
	return fmt.Errorf("cannot create inventory: %d, %w", r1, err)
}

func getTagDataReportCount(readerHandler int) int {
	r1, _, _ := rdr_GetTagDataReportCount.Call(uintptr(readerHandler))
	return int(r1)
}

func getTagDataReport(readerHandler int, seekType byte) int {
	// seek on current record
	// seekType - 1 first
	// seekType - 2 next
	// seekType - 3 last
	r1, _, _ := rdr_GetTagDataReport.Call(uintptr(readerHandler), uintptr(seekType))
	return int(r1)
}

func dNodeDestroy(dataNodeHandler int) {
	_, _, _ = dnode_Destroy.Call(uintptr(dataNodeHandler))
}

func resetCommuImmeTimeout(readerHandler int) {
	_, _, _ = rdr_ResetCommuImmeTimeout.Call(uintptr(readerHandler))
}

type ReportedTag struct {
	AipId, TagId, AntId int
	DsfId               byte
	Uid                 []byte
}

func parseTagDataReportISO15693(tagReportId int) (*ReportedTag, error) {
	var aipId, tagId, antId int
	var dsfId byte
	uid := make([]byte, 10, 10)
	r1, _, err := iso15693_ParseTagDataReport.Call(
		uintptr(tagReportId),
		uintptr(unsafe.Pointer(&aipId)),
		uintptr(unsafe.Pointer(&tagId)),
		uintptr(unsafe.Pointer(&antId)),
		uintptr(unsafe.Pointer(&dsfId)),
		uintptr(unsafe.Pointer(&uid[0])),
	)
	if r1 != 0 {
		return nil, fmt.Errorf("cannot parse tag data report: %d, %w", r1, err)
	}
	return &ReportedTag{
		AipId: aipId,
		TagId: tagId,
		AntId: antId,
		DsfId: dsfId,
		Uid:   uid[:8],
	}, nil
}

func parseTagDataReportISO14443A(tagReportId int) (*ReportedTag, error) {
	var aipId, tagId, antId int
	var dsfId, uidLen byte
	uid := make([]byte, 10, 10)
	r1, _, err := iso14443A_ParseTagDataReport.Call(
		uintptr(tagReportId),
		uintptr(unsafe.Pointer(&aipId)),
		uintptr(unsafe.Pointer(&tagId)),
		uintptr(unsafe.Pointer(&antId)),
		uintptr(unsafe.Pointer(&dsfId)),
		uintptr(unsafe.Pointer(&uid[0])),
		uintptr(unsafe.Pointer(&uidLen)),
	)
	if r1 != 0 {
		return nil, fmt.Errorf("cannot parse tag data report: %d, %w", r1, err)
	}
	return &ReportedTag{
		AipId: aipId,
		TagId: tagId,
		AntId: antId,
		DsfId: dsfId,
		Uid:   uid[:uidLen],
	}, nil
}

func tagDisconnect(tagId int) {
	_, _, _ = rdr_TagDisconnect.Call(uintptr(tagId))
}
