import axios from 'axios';

const apiClient = axios.create({
    baseURL: 'http://localhost:80/api',
    timeout: 10000,
    headers: {'Content-Type': 'application/json'}
});

export const registerUser = async (userData) => {
    try {
        const response = await apiClient.post('v1/user/register', userData);
        return response.data;
    } catch (error) {
        console.log(error)
        throw error.response || { data: { error: { message: 'Unknown error occured'} } };
    }
};

export const loginUser = async (userData) => {
    try {
        const response = await apiClient.post('v1/user/login', userData, {
            withCredentials: true
        });
        return response.data
    } catch (error) {
        console.log(error)
        throw error.response || { data: { error: { message: `Unknown error occured` } } }
    }
};

export const fetchTransactions = async (data) => {
    try {
        const response = await apiClient.get('v1/transaction/search', {
            params: data,
            withCredentials: true,
        });
        return response.data
    } catch(error) {
        // console.log(error)
        // throw error.response || {data: {error: { message: `Unknown error occured ` } } }
    }
};

export const getUserBalance = async () => {
    try {
        const response = await apiClient.get('v1/user/balance', {
            withCredentials: true,
        });
        console.log("balance respons:", response.data?.balance)
        return response.data?.balance;
    } catch (error) {
        console.log('cant get balance:', error);
        throw error.response || { data: { error: { message: 'Unknown error occured' }  } }
    }
};

export const createTransaction = async (createTrxData) => {
    try {
        const response = await apiClient.post('v1/transaction/create', createTrxData, {
            withCredentials: true,
        });
        console.log('create trx response:', response)
        return response.data;
    } catch (error) {
        console.log('cant create trx:', error)
        throw error.response || { data: { error: { message: 'Unknown error occured' }  } }
    }
};

// export const searchEntries = async (trxSearchData) => {
//     try {
//         const response = await apiClient.get('v1/transaction/', trxSearchData, {
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
        const resp = await apiClient.post('v1/user/confirm', confirmParams);
        console.log('confirm emial response:', resp)
        return resp.data;
    } catch(error) {
        console.log('cant confirm account:', error)
        throw error.response || { data: { error: { message: 'Unknown error occured' }  } }
    }
}

export default {
    registerUser,
}