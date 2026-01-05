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

export const createTransfer = async (createTrxData) => {
  try {
    const response = await apiClient.post("v1/transfer/create", createTrxData, {
      withCredentials: true,
    });
    console.log("create trx response:", response);
    return response.data;
  } catch (error) {
    console.log("cant create trx:", error);
    throw (
      error.response || {
        data: { error: { message: "Unknown error occured" } },
      }
    );
  }
};

// export const searchEntries = async (trxSearchData) => {
//     try {
//         const response = await apiClient.get('v1/transfer/', trxSearchData, {
//             withCredentials: true,
//         });
//         console.log('search trx response:', response);
//         return response.data;
//     } catch(error) {
//         // console.log('cant search trx:', error.response?.data)
//         // throw error.response || { data: { error: { message: 'Unknown error occured' }  } }
//     }
// };

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
