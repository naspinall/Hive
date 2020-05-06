import React, { useState, FormEvent, useContext } from "react";
import {
  Box,
  Field,
  Control,
  Input,
  Label,
  Button,
  Container,
  Column,
} from "bloomer";
import { Columns } from "bloomer/lib/grid/Columns";
import Axios from "axios";
import { AuthContext } from "../hooks/auth";
import { NotificationContext } from "../hooks/notification";

interface Props {}

const LoginDialog = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const { Set } = useContext(NotificationContext);
  const {Login} = useContext(AuthContext)

  const loginButton = async () => {
    try {
      const { data } = await Axios.post("/api/login", {
        email: username,
        password: password,
      });
      Login(data.username, data.displayName);
    } catch (error) {
      Set(error.message, "danger");
    }
  };

  return (
    <Container isFluid style={{ marginTop: 10 }}>
      <Columns isCentered isVCentered>
        <Column isSize="1/3">
          <Box>
            <Field>
              <Label isSize="medium">Username</Label>
              <Control>
                <Input
                  type="text"
                  placeholder=""
                  value={username}
                  onChange={(e: FormEvent<HTMLInputElement>) =>
                    setUsername(e.currentTarget.value)
                  }
                ></Input>
              </Control>
            </Field>
            <Field>
              <Label isSize="medium">Password</Label>
              <Control>
                <Input
                  type="password"
                  placeholder=""
                  value={password}
                  onChange={(e: FormEvent<HTMLInputElement>) =>
                    setPassword(e.currentTarget.value)
                  }
                />
              </Control>
            </Field>
            <Field>
              <Control>
                <Button onClick={loginButton}>Login</Button>
              </Control>
            </Field>
          </Box>
        </Column>
      </Columns>
    </Container>
  );
};

export default LoginDialog;
