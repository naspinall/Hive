import React from "react";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import "./App.css";
import { Nav } from "./containers/Nav";
import { Devices } from "./containers/Devices";
import Alarms from "./containers/Alarms";

const Routes = [
  { displayName: "Home", path: "/" },
  { displayName: "Settings", path: "/settings" },
  {
    displayName: "Manage",
    routes: [
      { displayName: "Devices", path: "/devices" },
      { displayName: "Measurements", path: "/measurements" },
      { displayName: "Alarms", path: "/alarms" }
    ]
  }
];

const App = () => {
  return (
    <div className="App">
      <Router>
        <header className="App-header">
          <Nav routes={Routes} />
        </header>
        <Switch>
          <Route path="/devices">
            <Devices />
          </Route>
          <Route path="/alarms">
            <Alarms />
          </Route>
        </Switch>
      </Router>
    </div>
  );
};

export default App;
