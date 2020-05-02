import {
  Container,
  Level,
  LevelLeft,
  LevelItem,
  Title,
  LevelRight,
  Button,
  Subtitle,
  Field,
  Control,
  Input,
  Icon
} from "bloomer";
import React, { useState } from "react";
import { useArrayRequest } from "../hooks/useRequest";
import { CreateDevice } from "./CreateDevice";
import { DeviceTable } from "./DeviceTable";
import { Device } from "../module";

interface Props {}

const DeviceView: React.FunctionComponent<Props> = () => {
  const [params, setParams] = useState({});
  const [search, setSearch] = useState("");
  const { response, isLoading } = useArrayRequest<Device>(
    "/api/devices",
    params
  );
  const [createActive, setCreateActive] = useState(false);
  const onSearch = (term: string, search: string) => {
    setParams({ [term]: search });
  };

  return (
    <Container isFluid style={{ marginTop: 10 }}>
      <Level>
        <LevelLeft>
          <LevelItem>
            <Title>Devices</Title>
          </LevelItem>
        </LevelLeft>
        <LevelRight>
          <LevelItem>
            <Button isColor="info" onClick={() => setCreateActive(true)}>
              Create
            </Button>
          </LevelItem>
        </LevelRight>
      </Level>
      <CreateDevice isActive={createActive} setActive={setCreateActive} />
      <Level>
        <LevelLeft>
          <LevelItem>
            <Subtitle>{response.length} Devices</Subtitle>
          </LevelItem>
          <LevelItem>
            <Field hasAddons>
              <Control>
                <Input
                  placeholder="Find a device"
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                    setSearch(e.target.value)
                  }
                ></Input>
              </Control>
              <Control>
                <Button onClick={() => onSearch("name", search)}>Search</Button>
              </Control>
            </Field>
          </LevelItem>
        </LevelLeft>
      </Level>
      {isLoading ? (
        <Icon isSize="large" className="fa fa-spinner fa-3x" />
      ) : (
        <DeviceTable devices={response} />
      )}
    </Container>
  );
};

export { DeviceView };
