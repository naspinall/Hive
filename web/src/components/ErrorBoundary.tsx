import React, { Component } from "react";
import { NotificationContext } from "../hooks/notification";
import { Icon } from "bloomer";

interface Props {
  error: Error;
  isLoading: boolean;
}

interface State {
  error: Error;
  hasError: boolean;
  isLoading: boolean;
}

class ErrorBoundary extends Component<Props, State> {
  static contextType = NotificationContext;

  state = {
    error: new Error("Error!"),
    hasError: false,
    isLoading: false,
  };

  componentDidCatch(error: Error) {
    this.setState({ error, hasError: true });
    console.error(error);
  }

  render() {
    console.log(this.state.isLoading)
    if (this.state.hasError) {
      //this.context.Set(this.state.error.message);
      return null;
    } else if (this.props.isLoading) {
      return <Icon isSize="large" className="fa fa-spinner fa-3x" />;
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
