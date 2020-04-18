import React from "react";
import { Notification, Delete, Icon } from "bloomer";

interface Props {
  error: string;
  isLoading: boolean;
}

const Loader: React.FunctionComponent<Props> = ({
  error,
  isLoading,
  children
}) => {
  console.log(children);
  if (error) {
    return (
      <Notification>
        <Delete />
        {error}
      </Notification>
    );
  } else if (isLoading) {
    return <Icon className="fa fa-cog fa-spin fa-3x fa-fw"></Icon>;
  } else {
    return <div> {children} </div>;
  }
};

export default Loader;
