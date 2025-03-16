function getCookie(name) {
    const cookies = document.cookie.split('; ');
    console.log("cookies:", cookies)
    for (const cookie of cookies) {
        const [key, value] = cookie.split('=');
        if (key === name) return decodeURIComponent(value);
    }

    return null;
}

export function getPasetoFooter(token) {
    try {
        const parts = token.split('.');
        if (parts.length !== 4) {
            throw new Error('Invalid PASETO token format');
        }

        console.log("token parts:", parts)
        const footer = atob(parts[3]); 
        console.log("footer:", footer)
        return footer.replaceAll('"', '');
    } catch (error) {
        console.error("Cant decode paseto token:", error);
        return null;
    }
}

export function getUsernameFromToken() {
    const token = getCookie('authToken');
    if (!token) return null;
    
    const usrname = getPasetoFooter(token);    
    return usrname;
}