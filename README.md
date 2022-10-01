# Remote rfid-reader service

Use connected rfid communication reader to to read cards and send cards numbers to the server via mqtt broker.

### Configuration

`id` - unique identifier of reader service

`device.driver` - now only `M302` is supported

```json
{
  "id": "work-station-1-reader-1",
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