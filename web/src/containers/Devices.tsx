import React from "react";
import { DeviceView } from "../components/DeviceView";
import { useRouteMatch, Switch, Route } from "react-router-dom";
import Device from "../components/SingleDevice";

interface DevicesProps {}

const Devices: React.FunctionComponent<DevicesProps> = () => {
  const match = useRouteMatch();

  return (
    <Switch>
      <Route path={`${match.path}/:deviceId`}>
        <Device />
      </Route>
      <Route path={match.path}>
        <DeviceView />
      </Route>
    </Switch>
  );
};

export { Devices };
