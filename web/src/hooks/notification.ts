import { createContext, Reducer, useReducer, Dispatch } from "react";

export type NotificationColor =
  | "white"
  | "light"
  | "dark"
  | "black"
  | "primary"
  | "info"
  | "success"
  | "warning"
  | "danger";

type NotificationActionType = "SHOW" | "HIDE";

interface NotificationState {
  message: string;
  active: boolean;
  color: NotificationColor;
}

interface NotificationAction {
  type: NotificationActionType;
  message?: string;
  color?: NotificationColor;
}

const initialNotification = {
  color: "white" as NotificationColor,
  message: "",
  active: false,
};

const reducer: Reducer<NotificationState, NotificationAction> = (
  state: NotificationState,
  action: NotificationAction
) => {
  switch (action.type) {
    case "SHOW":
      return {
        message: action.message ?? "", // Think for a better solution to this
        color: action.color ?? "white",
        active: true,
      };
    case "HIDE":
      return initialNotification;
    default:
      throw new Error();
  }
};

export const useNotification = () => {
  const [state, dispatch] = useReducer(reducer, initialNotification);

  const Set = (dispatch: Dispatch<NotificationAction>) => (
    message: string,
    color: NotificationColor,
    timeout?: number
  ) => {
  
    // Automatically hiding the notification.
    setTimeout(() => {
      dispatch({ type: "HIDE" });
    }, timeout || 3000);

    dispatch({
      type: "SHOW",
      message,
      color,
    });
  };

  const Reset = (dispatch: Dispatch<NotificationAction>) => () => {
    dispatch({
      type: "HIDE",
    });
  };

  return {
    Set: Set(dispatch),
    Reset: Reset(dispatch),
    NotificationState: state,
  };
};

export const NotificationContext = createContext({
    NotificationState : initialNotification,
    Set : (message: string, color: NotificationColor, timeout?: number | undefined) => {},
    Reset : () => {}
});
