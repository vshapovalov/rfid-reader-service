package rfidlib

import (
	"fmt"
	"strconv"
	"strings"
)

type Driver struct {
	Id       string
	Name     string
	Type     string
	CommType string
}

type HIDItem struct {
	SerialNum  string
	DriverPath string
}

type Reader struct {
	driver       *Driver
	antennaCount int
	connString   string
	handler      int
}

const ReaderUSBAddrModeNone = 0
const ReaderUSBAddrModeSerial = 1

type ReaderConnOptions interface {
	ReaderCOMOptions | ReaderUSBOptions | ReaderNETOptions
}

type ReaderCOMOptions struct {
	ComPort, ComBand, Frame string
}

type ReaderUSBOptions struct {
	AddrMode     int
	SerialNumber string
}

type ReaderNETOptions struct {
	RemoteIp, RemotePort, LocalIp string
}

func GetReaderConnectionString[T ReaderConnOptions](driverName string, options T) string {
	builder := strings.Builder{}
	builder.WriteString(ConnstrNameRdtype + "=" + driverName + ";")
	switch v := any(options).(type) {
	case ReaderCOMOptions:
		builder.WriteString(ConnstrNameCommtype + "=" + ConnstrNameCommtypeCom + ";")
		builder.WriteString(ConnstrNameComname + "=" + v.ComPort + ";")
		builder.WriteString(ConnstrNameCombarud + "=" + v.ComBand + ";")
		builder.WriteString(ConnstrNameComframe + "=" + v.Frame + ";")
		builder.WriteString(ConnstrNameBusaddr + "=255")
	case ReaderUSBOptions:
		builder.WriteString(ConnstrNameCommtype + "=" + ConnstrNameCommtypeUsb + ";")
		// None Addressed
		// Serial Number
		builder.WriteString(ConnstrNameHidaddrmode + "=" + strconv.Itoa(v.AddrMode) + ";")
		builder.WriteString(ConnstrNameHidsernum + "=" + v.SerialNumber)
	case ReaderNETOptions:
		builder.WriteString(ConnstrNameCommtype + "=" + ConnstrNameCommtypeNet + ";")
		builder.WriteString(ConnstrNameRemoteip + "=" + v.RemoteIp + ";")
		builder.WriteString(ConnstrNameRemoteport + "=" + v.RemotePort + ";")
		builder.WriteString(ConnstrNameLocalip + "=" + v.LocalIp)
	}
	return builder.String()
}

func NewReader[T ReaderConnOptions](driver *Driver, options T) *Reader {
	reader := Reader{
		driver:       driver,
		antennaCount: 0,
		connString:   GetReaderConnectionString(driver.Name, options),
		handler:      0,
	}
	return &reader
}

func (reader *Reader) Open() error {
	var err error
	reader.handler, err = rdrOpen(reader.connString)
	if err != nil {
		return err
	}
	reader.antennaCount = rdrGetAntennaInterfaceCount(reader.handler)
	return nil
}

func (reader *Reader) Close() error {
	if reader.isClosed() {
		return fmt.Errorf("reader is closed")
	}
	rh := reader.handler
	reader.handler = 0
	return rdrClose(rh)
}

func (reader *Reader) isClosed() bool {
	return reader.handler == 0
}
func (reader *Reader) ReadCards() ([][]byte, error) {
	if reader.isClosed() {
		return nil, fmt.Errorf("reader is closed")
	}
	cards := make([]*ReportedTag, 0)
	aiType := byte(1) // all tags
	//aiType := byte(2) // only new tags

	antennasIds := make([]byte, 0, reader.antennaCount)

	for i := 0; i < reader.antennaCount; i++ {
		antennasIds = append(antennasIds, byte(i+1))
	}

	dataNodeHandler := 0
	err := tagInventory(reader.handler, aiType, antennasIds, dataNodeHandler)
	if err != nil {
		return nil, err
	}
	_ = getTagDataReportCount(reader.handler)
	tagReportId := getTagDataReport(reader.handler, 1) // first record
	for tagReportId != 0 {
		card, err := getTagFromReport(tagReportId)
		if err == nil {
			cards = append(cards, card)
		}
		tagReportId = getTagDataReport(reader.handler, 2) // next record
	}
	dNodeDestroy(dataNodeHandler)
	resetCommuImmeTimeout(reader.handler)
	resultCards := make([][]byte, 0, len(cards))

	for _, card := range cards {
		resultCards = append(resultCards, card.Uid)
		tagDisconnect(card.TagId)
	}

	return resultCards, nil
}

func (reader *Reader) Buzz() {
	operationNumber := rdrCreateSetOutputOperations()
	rdrAddOneOutputOperation(operationNumber, 1, 1, 1, 1)
	rdrSetOutput(reader.handler, operationNumber)
	dNodeDestroy(operationNumber)
}

func getTagFromReport(tagReportId int) (*ReportedTag, error) {
	reportISO15693, err := parseTagDataReportISO15693(tagReportId)
	if err != nil {
		reportISO14443A, err := parseTagDataReportISO14443A(tagReportId)
		if err != nil {
			return nil, err
		}
		return reportISO14443A, nil
	}
	return reportISO15693, nil
}

func LoadDrivers(driversPath string) ([]*Driver, error) {

	err := loadReaderDrivers(driversPath)
	if err != nil {
		return nil, err
	}

	drvCount := getLoadedReaderDriverCount()
	drivers := make([]*Driver, 0, drvCount)

	for i := 0; i < drvCount; i++ {
		driver, err := GetDriver(i)
		if err != nil {
			continue
		}
		drivers = append(drivers, driver)
	}

	return drivers, nil
}

func GetDriver(driverNum int) (*Driver, error) {
	var err error

	driver := Driver{}

	driver.Id, err = getLoadedReaderDriverOpt(driverNum, LoadedRdrdvrOptId)
	if err != nil {
		return nil, err
	}
	driver.Name, err = getLoadedReaderDriverOpt(driverNum, LoadedRdrdvrOptName)
	if err != nil {
		return nil, err
	}
	driver.Type, err = getLoadedReaderDriverOpt(driverNum, LoadedRdrdvrOptCatalog)
	if err != nil {
		return nil, err
	}
	driver.CommType, err = getLoadedReaderDriverOpt(driverNum, LoadedRdrdvrOptCommtypesupported)
	if err != nil {
		return nil, err
	}
	return &driver, nil
}

func GetCOMPorts() ([]string, error) {
	comPortCount := comPortEnum()
	comPorts := make([]string, 0, comPortCount)
	for i := 0; i < comPortCount; i++ {
		comPort, err := comPortGetEnumItem(i)
		if err != nil {
			continue
		}
		comPorts = append(comPorts, comPort)
	}
	return comPorts, nil
}

func GetHIDItem(hidItemNum int) (*HIDItem, error) {
	var err error

	hidItem := HIDItem{}

	hidItem.SerialNum, err = hidGetEnumItemOpt(hidItemNum, HidEnumInfTypeSerialnum)
	if err != nil {
		return nil, err
	}
	hidItem.DriverPath, err = hidGetEnumItemOpt(hidItemNum, HidEnumInfTypeDriverpath)
	if err != nil {
		return nil, err
	}
	return &hidItem, nil
}

func GetHIDItems(driverName string) ([]*HIDItem, error) {
	itemsCount := hidEnum(driverName)
	items := make([]*HIDItem, 0, itemsCount)

	for i := 0; i < itemsCount; i++ {
		item, err := GetHIDItem(i)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	return items, nil
}
