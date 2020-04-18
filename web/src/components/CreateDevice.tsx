import React, { useState } from "react";
import {
  ModalContent,
  Modal,
  ModalBackground,
  ModalClose,
  Field,
  Label,
  Control,
  Input,
  Box,
  Button,
  Subtitle
} from "bloomer";
import { DeviceLike } from "../module";
import axios from "axios";
import cogoToast from "cogo-toast";

interface Props {
  isActive: boolean;
  setActive: (state: boolean) => void;
}

const CreateDevice: React.FunctionComponent<Props> = ({
  isActive,
  setActive
}: Props) => {
  const [name, setName] = useState("");

  const PostDevice = async (device: DeviceLike) => {
    try {
      await axios.post("/api/device", device);
      setActive(false);
    } catch (error) {
      cogoToast.error(error.message);
    }
  };

  return (
    <Modal isActive={isActive}>
      <ModalBackground />
      <ModalContent>
        <Box>
          <Subtitle>Create Device</Subtitle>
          <Field>
            <Label>Name</Label>
            <Control>
              <Input
                type="text"
                placeholder="Enter Name"
                value={name}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                  setName(e.target.value)
                }
              />
            </Control>
          </Field>
          <Field isGrouped>
            <Control>
              <Button isColor="primary" onClick={() => PostDevice({ name })}>
                Submit
              </Button>
            </Control>
            <Control>
              <Button isColor="warning" onClick={() => setActive(false)}>
                Cancel
              </Button>
            </Control>
          </Field>
        </Box>
      </ModalContent>
      <ModalClose />
    </Modal>
  );
};

export { CreateDevice };
