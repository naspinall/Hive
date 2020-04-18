import React from "react";
import { Section, Table, Title, Box, Column } from "bloomer";
import { Device } from "../module";
import { useRequest } from "../hooks/useRequest";
import { useParams } from "react-router-dom";
import Loader from "./Loader";

interface Props {}

const SingleDevice: React.FC<Props> = () => {
  const { deviceId } = useParams();
  const { response, error, isLoading } = useRequest<Device>(
    `/api/devices/${deviceId}`
  );
  return (
    <Loader isLoading={isLoading} error={error.message}>
      <Section>
        <Column>
          <Box>
            <Title>{response.name}</Title>
            <Table>
              <tbody>
                <tr>
                  <th>ID</th>
                  <th>{response.ID}</th>
                </tr>
                <tr>
                  <th>Name</th> <th>{response.name}</th>
                </tr>
                <tr>
                  <th>IMEI</th> <th>{response.imei}</th>
                </tr>
                <tr>
                  <th>CreatedAt</th> <th>{response.CreatedAt}</th>
                </tr>
                <tr>
                  <th>Latitude</th> <th>{response.latitude}</th>
                </tr>
                <tr>
                  <th>Longitude</th> <th>{response.longitude}</th>
                </tr>
              </tbody>
            </Table>
          </Box>
        </Column>
      </Section>
    </Loader>
  );
};

export default SingleDevice;
