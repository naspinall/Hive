import React, { useContext } from "react";
import { Delete, Notification } from "bloomer";
import { NotificationContext } from "../hooks/notification";

interface Props {}

const TopNotification = (props: Props) => {
  const { NotificationState, Reset } = useContext(NotificationContext);
  if (NotificationState.active)
    return (
      <Notification isColor={NotificationState.color} style={{ margin: "1em" }}>
        <Delete onClick={Reset} />
        <p>{NotificationState.message}</p>
      </Notification>
    );
  else {
    return <div />;
  }
};

export default TopNotification;
