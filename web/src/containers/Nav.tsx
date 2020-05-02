import * as React from "react";
import {
  Navbar,
  NavbarBrand,
  NavbarItem,
  Icon,
  NavbarBurger,
  NavbarMenu,
  NavbarStart,
  NavbarDropdown,
  NavbarEnd,
  Button,
} from "bloomer";
import { Link } from "react-router-dom";
import { AuthContext } from "../hooks/auth";
import { useContext } from "react";

interface RouteDescription {
  path: string;
  displayName: string;
  //component: React.FunctionComponent;
}

interface RouteMenu {
  displayName: string;
  routes: RouteDescription[];
}

interface NaveRouteMenuProps {
  menu: RouteMenu;
}

const NavRouteMenu: React.FunctionComponent<NaveRouteMenuProps> = (props) => {
  return (
    <NavbarItem hasDropdown isHoverable>
      <NavbarItem>{props.menu.displayName}</NavbarItem>
      <NavbarDropdown>
        {props.menu.routes.map((menuItem) => (
          <NavbarItem
            style={{ cursor: "pointer" }}
            key={`goto-${menuItem.path}`}
          >
            <Link to={menuItem.path}>{menuItem.displayName}</Link>
          </NavbarItem>
        ))}
      </NavbarDropdown>
    </NavbarItem>
  );
};

interface NavProps {
  routes: Array<RouteMenu | RouteDescription>;
}

const Nav: React.FunctionComponent<NavProps> = (props) => {
  const { AuthState } = useContext(AuthContext);
  return (
    <Navbar style={{ borderBottom: "solid 1px #FFFFFF", margin: "0" }}>
      <NavbarBrand>
        <NavbarItem>
          <img
            src="/images/logo.svg"
            style={{ marginRight: 5 }}
            alt="Hive Logo"
          />{" "}
          Hive
        </NavbarItem>
        <NavbarItem isHidden="desktop">
          <Icon className="fa fa-github" />
        </NavbarItem>
        <NavbarItem isHidden="desktop">
          <Icon className="fa fa-twitter" style={{ color: "#55acee" }} />
        </NavbarItem>
        <NavbarBurger
        //isActive={this.state.isActive}
        //onClick={this.onClickNav}
        />
      </NavbarBrand>
      <NavbarMenu hasTextColor="dark">
        <NavbarStart>
          {props.routes.map((route) => {
            return "routes" in route ? (
              <NavRouteMenu
                menu={route}
                key={`nav-key-for(${route.displayName})`}
              />
            ) : (
              <NavbarItem key={`nav-key-for(${route.displayName})`}>
                <Link to={route.path}>{route.displayName}</Link>
              </NavbarItem>
            );
          })}
        </NavbarStart>
        <NavbarEnd>
          {AuthState.isAuthenticated ? (
            <NavbarItem>{AuthState.displayName}</NavbarItem>
          ) : (
            <Link to="/signup">
              <Button isColor="info" style={{ margin: "1em" }}>
                New User
              </Button>
            </Link>
          )}
          <NavbarItem href="https://github.com/naspinall/hive" isHidden="touch">
            <Icon className="fa fa-github" />
          </NavbarItem>
        </NavbarEnd>
      </NavbarMenu>
    </Navbar>
  );
};

export { Nav };
