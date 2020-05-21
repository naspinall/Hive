import { useState, useEffect, useContext } from "react";
import { useAuth, AuthContext } from "./auth";
import axios, { AxiosRequestConfig } from "axios";

interface Response<T> {
  response: T;
  error: Error;
  isLoading: boolean;
}

const useRequest = <T>(
  path: string,
  params?: AxiosRequestConfig
): Response<T> => {
  const [response, setResponse] = useState({} as T);
  const [error, setError] = useState({} as Error);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    setIsLoading(true);
    const doRequest = async () => {
      try {
        const res = await axios.get(path, {
          params,
          headers: {
            Authorization: `Bearer ${localStorage.token}`,
          },
        });
        setResponse(res.data);
        setIsLoading(false);
      } catch (error) {
        setError(error);
        setIsLoading(false);
      }
    };
    doRequest();
  }, [path, params]);

  return { response, error, isLoading };
};

const useArrayRequest = <T>(
  path: string,
  params?: any
): { response: T[]; error: Error; isLoading: boolean } => {
  const { response, error, isLoading } = useRequest<T[]>(path, params);
  if (Array.isArray(response)) {
    return { response, error, isLoading };
  } else {
    return { response: [], error, isLoading };
  }
};

export { useRequest, useArrayRequest };
