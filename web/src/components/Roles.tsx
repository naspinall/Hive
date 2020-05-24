import React from "react";
import { Box } from "bloomer";
import { useRequest } from "../hooks/useRequest";
import { Role, Permission } from "../module";

interface Props {
  id: number;
}

const mapPermission = (role: number): Permission => {
  switch (role) {
    case 0:
      return "NONE";
    case 1:
      return "READ";
    case 2:
      return "CREATE";
    case 3:
      return "UPDATE";
    case 4:
      return "DELETE";
    default:
      return "NONE";
  }
};

const Roles = ({ id }: Props) => {
  const { response, isLoading } = useRequest<Role>(`/api/users/${id}/roles`);
  return (
    <Box>
      <p>Alarms: {mapPermission(response.alarms)}</p>
      <p>Users: {mapPermission(response.users)}</p>
      <p>Measurements: {mapPermission(response.measurements)}</p>
      <p>Devices: {mapPermission(response.devices)}</p>
      <p>Subscriptions: {mapPermission(response.subscriptions)}</p>
    </Box>
  );
};

export default Roles;
