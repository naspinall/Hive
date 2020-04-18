export interface Device {
  ID: number;
  name: string;
  imei: number;
  CreatedAt: string;
  latitude: number;
  longitude: number;
}

export interface DeviceLike {
  name: string;
}

export interface Alarm {
  Type: string;
  Status: string;
  Severity: string;
  DeviceID: number;
}
