import paho.mqtt.client as mqtt
import time
import json
import datetime

mqtt_connected = False

def on_connect(client, userdata, flags, rc):
    global mqtt_connected
    if rc == 0:
        mqtt_connected = True
        print("MQTT Connected successfully")
    else:
        print(f"MQTT Connection failed with code {rc}")

def on_disconnect(client, userdata, rc):
    global mqtt_connected
    mqtt_connected = False
    print(f"MQTT Disconnected with code {rc}")

def wait_for_mqtt_connection(timeout=10):
    start_time = time.time()
    while not mqtt_connected and (time.time() - start_time) < timeout:
        time.sleep(0.1)
    return mqtt_connected


if __name__ == "__main__":
    client = mqtt.Client()
    client.on_connect = on_connect
    client.on_disconnect = on_disconnect
    client.connect("localhost", 1884, 60)
    client.loop_start()
    start = 0

    while True:
        data = []
        for i in range(50):
            data.append(start + i)


        measurementData = {
            'scenario_id': 19,
            'device_id': 1,
            'timestamp': datetime.datetime.now(datetime.UTC).strftime('%Y-%m-%d %H:%M:%S.%f')[:-3],
            "data": {
            "acc_x": data,
            "acc_y": data,
            "acc_z": data,
            "gyro_x": data,
            "gyro_y": data,
            "gyro_z": data,
            "curr_v": data,
            "temp": data
        }
        }
        client.publish("device/1/raw", json.dumps(measurementData))
        start += 50
        print(f"Sent {start} samples")
        time.sleep(0.1)
