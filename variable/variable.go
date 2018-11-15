package variable

const (

	//gateway to Platform topics
	G2P_Connect    = "v1/gateway/connect" // -m {"serialNumber":}
	G2P_Disconnect = "v1/gateway/disconnect"
	G2P_Telemetry  = "v1/gateway/telemetry"
	G2P_RPC        = "v1/gateway/rpc"

	//device to gateway Topic
	D2G_Sensor           = "v1/sensors"           // -t  "v1/sensors" -m {"serialNumber":"SN-001", "model":"T1000", "temperature":36.6}
	D2G_Sensor_embed     = "v1/sensors/+/+"       //sensor/SN-004/temperature
	D2G_Connect          = "v1/sensors/connect"   // -t "v1/sensors/connect" -m '{"serialNumber":"SN-001"}'
	D2G_Connect_embed    = "v1/sensors/+/connect" // -t "v1/sensors/SN-001/connect" -m ''
	D2G_Disconnect       = "v1/sensors/disconnect"
	D2G_Disconnect_embed = "v1/sensors/+/disconnect"
)
