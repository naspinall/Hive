import React, { useState } from "react";
import { Table, Column, Columns } from "bloomer";
import { User } from "../module";
import Roles from "./Roles";

interface Props {
  users: User[];
}
const UsersTable: React.FunctionComponent<Props> = ({ users }: Props) => {
  const [userID, setUserID] = useState(0);
  return (
    <Columns isCentered>
      <Column>
        <Table isBordered isStriped>
          <thead>
            <tr>
              <th>User ID</th>
              <th>Display Name</th>
              <th>Email</th>
            </tr>
          </thead>
          <tbody>
            {users.map((user) => (
              <tr
                key={`usertable-${user.ID}`}
                className={userID === user.ID ? "is-selected" : ""}
                onClick={() => setUserID(user.ID)}
              >
                <td>{user.ID}</td>
                <td>{user.displayName}</td>
                <td>{user.email}</td>
              </tr>
            ))}
          </tbody>
        </Table>
      </Column>
      <Column>
        <Roles id={userID}></Roles>
      </Column>
    </Columns>
  );
};

export { UsersTable };
