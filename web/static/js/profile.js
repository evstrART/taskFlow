// Display current year in the footer
document.getElementById('year').textContent = new Date().getFullYear();

// Fetch user profile
async function fetchUserProfile() {
    const token = localStorage.getItem('token');
    console.log("Token:", token); // Проверка наличия токена
    if (!token) {
        alert("No token found. Please log in.");
        return;
    }

    const userId = parseJwt(token).UserId; // Извлекаем UserId из токена
    console.log("User ID:", userId); // Проверка полученного UserId

    try {
        const userResponse = await fetch(`http://localhost:8080/api/users/${userId}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });

        if (!userResponse.ok) {
            const errorText = await userResponse.text(); // Получаем текст ошибки
            throw new Error(`Error fetching user profile: ${errorText}`);
        }

        const user = await userResponse.json();
        displayUserProfile(user);
    } catch (error) {
        console.error('Error:', error);
        alert(error.message);
    }
}

// Display user profile
function displayUserProfile(user) {
    console.log("User Data:", user); // Проверяем, что пришло от API
    document.getElementById('username').textContent = user.username;
    document.getElementById('email').textContent = user.email;
    document.getElementById('role').textContent = user.role || 'Not available';
}

// Parse JWT token
function parseJwt(token) {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
}

// Edit profile function
function editProfile() {
    alert("Edit profile functionality is not yet implemented.");
}
function deleteProfile(){
    alert("Delete Account")
}

// Fetch user profile on window load
window.onload = fetchUserProfile;