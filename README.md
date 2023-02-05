# Remote rfid-reader service

Use connected rfid communication reader to read and send cards numbers to the server via mqtt broker.

### Configuration

`id` - unique identifier of reader service

`device.driver` -  `M302` or `RFIDLib`

```json
{
  "id": "workstation-1-reader-1",
  "isDebugModeEnabled": true,
  "reverseCardNumber": false,
  "useBuzzerOnRead": true,
  "device": {
    "driver": "M302",
    "M302Settings": {
      "port": "COM1",
      "baud": 9600,
      "readTimeout": 1,
      "size": 8
    }
  },
  "mqttBroker": {
    "uri": "tcp://localhost:1883",
    "username": "user",
    "password": "secret"
  }
}
```

### Broker communications

When cards are read, service sends message to the broker with topic `reader/{readerId}/card` and payload with cards numbers.

```json
{
  "cardNumbers": ["1234567890"],
  "readerId": "workstation-1-reader-1"
}
```

When service starts, it sends message to the broker with topic `reader/{readerId}/status` and payload with status.

```json
{
  "status": "online",
  "readerId": "workstation-1-reader-1"
}
```

When service stops, it sends message to the broker with topic `reader/{readerId}/status` and payload with status.

```json
{
  "status": "offline",
  "readerId": "workstation-1-reader-1"
}
```

To use buzzer on reader module, send message to the broker with topic `reader/{readerId}/buzzer` and payload with
status.

```json
{
  "count": 1
}
```

### RFIDLib driver

Can be used only on windows x64, because uses windows dynamic libraries

This drivers library supports list of readers:

- M201
- MR113R
- RD120M
- RD201
- RD242
- SSR100
- RD5100
- RD5200
- RL8000
- RPAN

in device config section should be specified RFIDLIB driver and RFIDLibSettings settings, last one describes library devices driver, `RPAN` for  example, and comunication settings

```
{
  "driver": "RFIDLIB",
  "RFIDLibSettings": {
    "libDriver": "RPAN",
    "communication": {
      "type": "USB",
      "settings": {
        "serialNumber": "2072190002"
      }
    }
  }
}
```

All allowed settings you can find in `config.example.json` file
