import axios from "axios";
import { getCookie } from "@/utils/auth";

const apiClient = axios.create({
  baseURL: "/api",
  timeout: 10000,
  headers: { "Content-Type": "application/json" },
});

apiClient.interceptors.request.use((config) => {
  const token = getCookie("authToken");
  if (token && !config.headers.Authorization) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

// -------------- USERS --------------

export const registerUser = async (userData) => {
  try {
    const response = await apiClient.post("v1/user/register", userData);
    return response.data;
  } catch (error) {
    console.log(error);
    throw (
      error.response || {
        data: { error: { message: "Unknown error occured" } },
      }
    );
  }
};

export const loginUser = async (userData) => {
  try {
    const response = await apiClient.post("v1/user/login", userData, {
      withCredentials: true,
    });
    return response.data;
  } catch (error) {
    console.log(error);
    throw (
      error.response || {
        data: { error: { message: "Unknown error occured" } },
      }
    );
  }
};

export const confirmUserEmail = async (confirmParams) => {
  try {
    const resp = await apiClient.get("v1/user/confirm", {
      params: confirmParams,
    });
    console.log("confirm emial response:", resp);
    return resp.data;
  } catch (error) {
    console.log("cant confirm account:", error);
    throw (
      error.response || {
        data: { error: { message: "Unknown error occured" } },
      }
    );
  }
};

// -------------- ACCOUNTS --------------

export const fetchAccounts = async (data) => {
  try {
    const response = await apiClient.get("v1/transfer/accounts", {
      params: data,
      withCredentials: true,
    });
    return response.data;
  } catch (error) {
    console.log("cant fetch accounts:", error);
    throw (
      error.response || {
        data: { error: { message: "Unknown error occured" } },
      }
    );
  }
};

export const createAccount = async (data) => {
  try {
    const response = await apiClient.post("v1/transfer/account", data, {
      withCredentials: true,
    });
    return response.data;
  } catch (error) {
    console.log("cant create account:", error);
    throw (
      error.response || {
        data: { error: { message: "Unknown error occured" } },
      }
    );
  }
};

export const deleteAccount = async (data) => {
  try {
    const response = await apiClient.delete("v1/transfer/account", {
      params: data,
      withCredentials: true,
    });
    return response.data;
  } catch (error) {
    console.log("cant delete account:", error);
    throw (
      error.response || {
        data: { error: { message: "Unknown error occured" } },
      }
    );
  }
};

// -------------- TRANSFERS --------------

export const fetchTransfers = async (data) => {
  try {
    const response = await apiClient.get("v1/transfer", {
      params: data,
      withCredentials: true,
    });
    return response.data;
  } catch (error) {
    console.log("cant fetch transfers:", error);
    throw (
      error.response || {
        data: { error: { message: "Unknown error occured" } },
      }
    );
  }
};

export const createTransfer = async (data) => {
  try {
    const response = await apiClient.post("v1/transfer", data, {
      withCredentials: true,
    });
    return response.data;
  } catch (error) {
    console.log("cant create trx:", error);
    throw (
      error.response || {
        data: { message: "Unknown error occured" },
      }
    );
  }
};

export const fetchCurrencies = async () => {
  try {
    const response = await apiClient.get("v1/transfer/currencies", {
      withCredentials: true,
    });
    console.log("fetched currencies:", response.data);
    return response.data;
  } catch (error) {
    console.log("failed to fetch currencies:", error);
    throw (
      error.response || {
        data: { error: { message: "Unknown error occured" } },
      }
    );
  }
};
