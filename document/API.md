##MQTT topics which is related to the Gateway 
* v1/gateway/connect // To indicate an specific device is connected to the gateway
* v1/gateway/disconnect // to indicate an specific device is disconnected to the gateway
* v1/gateway/telemetry // this topic is revolving around the telemetry data of devices. 
{
  "MAC address": [
    {
      "ts": 1483228800000,
      "values": {
        "temperature": 42,
        "humidity": 80
      }
    },
    {
      "ts": 1483228801000,
      "values": {
        "temperature": 43,
        "humidity": 82
      }
    }
  ],
  "Device B": [
    {
      "ts": 1483228800000,
      "values": {
        "temperature": 42,
        "humidity": 80
      }
    }
  ]
}
