# Remote rfid-reader service

Use connected rfid communication reader to read and send cards numbers to the server via mqtt broker.

### Configuration

`id` - unique identifier of reader service

`device.driver` - now only `M302` is supported

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

When card is read, service sends message to the broker with topic `reader/{readerId}/card` and payload with card number.

```json
{
  "cardNumber": "1234567890",
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
