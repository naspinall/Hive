export type Permission = "READ" | "CREATE" | "UPDATE" | "DELETE" | "NONE";

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

export interface User {
  ID: number;
  email: string;
  displayName: string;
  role: {
    alarms: Permission;
    users: Permission;
    measurements: Permission;
    devices: Permission;
    subscriptions: Permission;
  };
}

export interface Role {
  alarms: number;
  users: number;
  measurements: number;
  devices: number;
  subscriptions: number;
}
