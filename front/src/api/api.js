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
        throw error.response || { data: { error: { messagee: `Unknown error occured` } } }
    }
}

export default {
    registerUser,
}