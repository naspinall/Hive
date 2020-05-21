import React from "react";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import "./App.css";
import { Nav } from "./containers/Nav";
import { Devices } from "./containers/Devices";
import Alarms from "./containers/Alarms";
import { useAuth, AuthContext } from "./hooks/auth";
import LoginDialog from "./components/LoginDialog";
import { NotificationContext, useNotification } from "./hooks/notification";
import TopNotification from "./components/TopNotification";
import SignupDialog from "./components/SignupDialog";
import Users from "./containers/Users";

const Routes = [
  { displayName: "Home", path: "/" },
  { displayName: "Settings", path: "/settings" },
  {
    displayName: "Manage",
    routes: [
      { displayName: "Devices", path: "/devices" },
      { displayName: "Measurements", path: "/measurements" },
      { displayName: "Alarms", path: "/alarms" },
      { displayName: "Users", path : "/users"},
    ],
  },
];

const App = () => {
  const { AuthState, Login, Logout, ExpireToken } = useAuth();
  const { NotificationState, Set, Reset } = useNotification();

  return (
    <AuthContext.Provider value={{ AuthState, Login, Logout, ExpireToken }}>
      <NotificationContext.Provider value={{ NotificationState, Set, Reset }}>
        <div className="App">
          <Router>
            <header className="App-header">
              <Nav routes={Routes} />
            </header>
            <TopNotification />
            <Switch>
              <Route path="/" exact={true}>
                {AuthState.isAuthenticated ? <p> Welcome </p> : <LoginDialog />}
              </Route>
              <Route path="/devices">
                <Devices />
              </Route>
              <Route path="/alarms">
                <Alarms />
              </Route>
              <Route path="/users">
                <Users />
              </Route>
              <Route path="/signup">
                <SignupDialog />
              </Route>
            </Switch>
          </Router>
        </div>
      </NotificationContext.Provider>
    </AuthContext.Provider>
  );
};

export default App;
