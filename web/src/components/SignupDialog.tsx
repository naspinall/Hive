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
import { NotificationContext } from "../hooks/notification";

interface Props {}

const SignupDialog = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [displayName, setDisplayName] = useState("");

  const { Set } = useContext(NotificationContext);

  const loginButton = async () => {
    try {
      await Axios.post("/api/users", {
        username: username,
        displayName: displayName,
        password: password,
      });
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
              <Label isSize="medium">Display Name</Label>
              <Control>
                <Input
                  type="text"
                  placeholder=""
                  value={displayName}
                  onChange={(e: FormEvent<HTMLInputElement>) =>
                    setDisplayName(e.currentTarget.value)
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
                <Button onClick={loginButton}>Create User</Button>
              </Control>
            </Field>
          </Box>
        </Column>
      </Columns>
    </Container>
  );
};

export default SignupDialog;
