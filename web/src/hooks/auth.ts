import { createContext, Reducer, useReducer, Dispatch } from "react";

type AuthActionType = "LOGIN" | "LOGOUT" | "TOKEN_EXPIRED";

interface AuthState {
  username: string;
  displayName: string;
  isAuthenticated: boolean;
  message: string;
  token?: string;
}

interface AuthAction {
  type: AuthActionType;
  username?: string;
  displayName?: string;
  message?: string;
  token?: string;
}

const initialAuth = {
  username: "",
  displayName: "Nick Aspinall",
  isAuthenticated: false,
  message: "",
};

const reducer: Reducer<AuthState, AuthAction> = (
  state: AuthState,
  action: AuthAction
) => {
  switch (action.type) {
    case "LOGIN":
      localStorage.token = action.token;
      return {
        username: action.username ?? "", // Think for a better solution to this
        displayName: action.displayName ?? "",
        message: `Welcome ${action.displayName}`,
        isAuthenticated: true,
      };
    case "LOGOUT":
      localStorage.token = "";
      return {
        username: "",
        displayName: "",
        message: "Sucuessfully Logged Out",
        isAuthenticated: false,
      };
    case "TOKEN_EXPIRED":
      localStorage.token = "";
      return {
        username: "",
        displayName: "",
        message: "Token Expired",
        isAuthenticated: false,
      };
    default:
      throw new Error();
  }
};

export const useAuth = () => {
  const [state, dispatch] = useReducer(reducer, initialAuth);
  const Login = (dispatch: Dispatch<AuthAction>) => (
    username: string,
    displayName: string,
    token: string
  ) => {
    dispatch({
      type: "LOGIN",
      username,
      displayName,
      token,
    });
  };
  const Logout = (dispatch: Dispatch<AuthAction>) => (
    username: string,
    displayName: string
  ) => {
    dispatch({
      type: "LOGOUT",
      username,
      displayName,
    });
  };
  const ExpireToken = (dispatch: Dispatch<AuthAction>) => (
    username: string,
    displayName: string
  ) => {
    dispatch({
      type: "TOKEN_EXPIRED",
      username,
      displayName,
    });
  };

  return {
    Login: Login(dispatch),
    Logout: Logout(dispatch),
    ExpireToken: ExpireToken(dispatch),
    AuthState: state,
    AuthDispatch: dispatch,
  };
};

export const AuthContext = createContext({
  AuthState: initialAuth,
  Login: (username: string, displayName: string, token : string) => {},
  Logout: (username: string, displayName: string) => {},
  ExpireToken: (username: string, displayName: string) => {},
});
