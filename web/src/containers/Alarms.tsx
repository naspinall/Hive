import React, { ReactElement } from "react";
import { useArrayRequest } from "../hooks/useRequest";
import { Title, Container, Level, LevelLeft, LevelItem, Icon } from "bloomer";
import { Alarm } from "../module";
import { AlarmTable } from "../components/AlarmTable";
import ErrorBoundary from "../components/ErrorBoundary";

interface Props {}

export default function Alarms(): ReactElement {
  const { response, isLoading, error } = useArrayRequest("/api/alarms");
  return (
    <Container isFluid style={{ marginTop: 10 }}>
      <Level>
        <LevelLeft>
          <LevelItem>
            <Title>Alarms</Title>
          </LevelItem>
        </LevelLeft>
      </Level>
      <ErrorBoundary isLoading={isLoading} error={error}>
        <AlarmTable alarms={response as Alarm[]} />
      </ErrorBoundary>
    </Container>
  );
}
