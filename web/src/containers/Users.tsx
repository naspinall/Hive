import React, { ReactElement } from "react";
import { useArrayRequest } from "../hooks/useRequest";
import { Title, Container, Level, LevelLeft, LevelItem, Icon } from "bloomer";
import { User } from "../module";
import { UsersTable } from "../components/UsersTable";

export default function Users(): ReactElement {
  const { response, isLoading } = useArrayRequest("/api/users");
  return (
    <Container isFluid style={{ marginTop: 10 }}>
      <Level>
        <LevelLeft>
          <LevelItem>
            <Title>Users</Title>
          </LevelItem>
        </LevelLeft>
      </Level>
      {isLoading ? (
        <Icon isSize="large" className="fa fa-spinner fa-3x" />
      ) : (
        <UsersTable users={response as User[]} />
      )}
    </Container>
  );
}
