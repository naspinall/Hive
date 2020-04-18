import React from "react";
import { Table, Button } from "bloomer";
import { Device } from "../module";
import moment from "moment";
import { Link, useRouteMatch } from "react-router-dom";

interface Props {
  devices: Device[];
}
const DeviceTable: React.FunctionComponent<Props> = ({ devices }: Props) => {
  const match = useRouteMatch();
  return (
    <Table isBordered isStriped>
      <thead>
        <tr>
          <th>ID</th>
          <th>Name</th>
          <th>IMEI</th>
          <th>Registration Date</th>
          <th>Latitude</th>
          <th>Longitude</th>
        </tr>
      </thead>
      <tbody>
        {devices.map(device => (
          <tr>
            <td>{device.ID}</td>
            <td>{device.name}</td>
            <td>{device.imei}</td>
            <td>{moment(device.CreatedAt).format("MMMM Do YYYY, h:mm a")}</td>
            <td>{device.latitude}</td>
            <td>{device.longitude}</td>
            <td>
              <Button>
                <Link to={`${match.path}/${device.ID}`}>
                  See More
                </Link>
              </Button>
            </td>
          </tr>
        ))}
      </tbody>
    </Table>
  );
};

export { DeviceTable };
