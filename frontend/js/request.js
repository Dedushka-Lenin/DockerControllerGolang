async function loadData(method, url, data) {
    try {
        const requestMethod = method.toUpperCase();

        const options = {
            method: requestMethod,
            headers: {
                'Content-Type': 'application/json'
            }
        };

        if (requestMethod !== 'GET' && requestMethod !== 'HEAD' && data !== undefined && data !== null && data !== '') {
            options.body = JSON.stringify(data);
        }

        const response = await fetch(url, options);

        let result = null;
        const contentType = response.headers.get("content-type");
        if (contentType && contentType.includes("application/json")) {
            result = await response.json();
        }

        if (!response.ok) {
            const errorMessage = (result && result.error) 
                ? result.error 
                : `HTTP error! status: ${response.status}`;
            throw new Error(errorMessage);
        }

        if (result && result.error) {
            throw new Error(result.error);
        }

        return result;
    } catch (error) {
        console.error('Request error:', error);
        throw error;
    }
}