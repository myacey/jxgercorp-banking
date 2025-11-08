function getCookie(name) {
    const cookies = document.cookie.split('; ');
    console.log("cookies:", cookies)
    for (const cookie of cookies) {
        const [key, value] = cookie.split('=');
        if (key === name) return decodeURIComponent(value);
    }

    return null;
}

export function getUsernameFromToken() {
    const token = getCookie('authToken');
    if (!token) return null;

    try {
        const payloadBase64 = token.split('.')[1];
        if(!payloadBase64) {
            console.error('JWT format invalid');
            return null;
        }

        const payloadJson = atob(payloadBase64.replace(/-/g, '+').replace(/_/g, '/'));
        const payload = JSON.parse(payloadJson);

        return payload.username || null;
    } catch (err) {
        console.error('Cannot decode JWT:', err)
        return null;
    }
}