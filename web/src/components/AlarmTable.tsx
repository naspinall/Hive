import React from "react";
import { Table } from "bloomer";
import { Alarm } from "../module";

interface Props {
  alarms: Alarm[];
}
const AlarmTable: React.FunctionComponent<Props> = ({ alarms }: Props) => {
  return (
    <Table isBordered isStriped>
      <thead>
        <tr>
          <th>Type</th>
          <th>Status</th>
          <th>Severity</th>
          <th>DeviceID</th>
        </tr>
      </thead>
      <tbody>
        {alarms.map(alarm => (
          <tr>
            <td>{alarm.Type}</td>
            <td>{alarm.Status}</td>
            <td>{alarm.Severity}</td>
            <td>{alarm.DeviceID}</td>
          </tr>
        ))}
      </tbody>
    </Table>
  );
};

export { AlarmTable };
