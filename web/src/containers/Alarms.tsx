import React, { ReactElement } from "react";
import { useArrayRequest } from "../hooks/useRequest";
import { Title, Container, Level, LevelLeft, LevelItem, Icon } from "bloomer";
import { Alarm } from "../module";
import { AlarmTable } from "../components/AlarmTable";

interface Props {}

export default function Alarms(): ReactElement {
  const { response, isLoading } = useArrayRequest("/api/alarms");
  return (
    <Container isFluid style={{ marginTop: 10 }}>
      <Level>
        <LevelLeft>
          <LevelItem>
            <Title>Alarms</Title>
          </LevelItem>
        </LevelLeft>
      </Level>
      {isLoading ? (
        <Icon isSize="large" className="fa fa-spinner fa-3x" />
      ) : (
        <AlarmTable alarms={response as Alarm[]} />
      )}
    </Container>
  );
}
