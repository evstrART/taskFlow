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

    const user_id = parseJwt(token).UserId; // Извлекаем UserId из токена
    console.log("User ID:", user_id); // Проверка полученного UserId

    try {
        const userResponse = await fetch(`http://localhost:8080/api/users/${user_id}`, {
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
    document.getElementById('username').textContent = user.username || 'Not available';
    document.getElementById('email').textContent = user.email || 'Not available';
    document.getElementById('role').textContent = user.role || 'Роль не указана';
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

function deleteProfile() {
    alert("Delete Account");
}

function logout() {
    // Удаляем токен из localStorage
    localStorage.removeItem('token');

    // Перенаправляем пользователя на страницу входа
    window.location.href = '/'; // Замените на актуальный путь к странице входа
}
// Fetch user profile on window load
window.onload = fetchUserProfile;
// Открытие модального окна

async function changePassword() {
    const oldPasswordInput = document.getElementById('old-password');
    const newPasswordInput = document.getElementById('new-password');
    const confirmNewPasswordInput = document.getElementById('confirm-new-password');

    // Сбрасываем предыдущие ошибки
    oldPasswordInput.classList.remove('error');
    newPasswordInput.classList.remove('error');
    confirmNewPasswordInput.classList.remove('error');
    document.getElementById('old-password-error').style.display = 'none';
    document.getElementById('new-password-error').style.display = 'none';
    document.getElementById('confirm-new-password-error').style.display = 'none';

    const oldPassword = oldPasswordInput.value;
    const newPassword = newPasswordInput.value;
    const confirmNewPassword = confirmNewPasswordInput.value;

    // Проверка на пустые поля
    if (!oldPassword || !newPassword || !confirmNewPassword) {
        alert("Пожалуйста, заполните все поля.");
        return;
    }

    // Проверка на совпадение нового пароля и подтверждения
    if (newPassword !== confirmNewPassword) {
        // Выделяем поля с ошибками
        newPasswordInput.classList.add('error'); // Добавляем класс для красной рамки
        confirmNewPasswordInput.classList.add('error'); // Добавляем класс для красной рамки

        // Отображаем сообщения об ошибках
        document.getElementById('new-password-error').textContent = "Новые пароли не совпадают.";
        document.getElementById('new-password-error').style.display = 'block';
        document.getElementById('confirm-new-password-error').textContent = "Новые пароли не совпадают.";
        document.getElementById('confirm-new-password-error').style.display = 'block';

        return;
    }

    const token = localStorage.getItem('token');

    if (!token) {
        alert("Токен не найден. Пожалуйста, войдите в систему.");
        return;
    }

    const userId = parseJwt(token).UserId; // Предполагается, что у вас есть функция parseJwt

    try {
        const response = await fetch(`http://localhost:8080/api/users/${userId}/change-password`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ old_password: oldPassword, new_password: newPassword })
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Ошибка при смене пароля: ${errorText}`);
        }

        alert("Пароль успешно изменён!");
        closeChangePasswordModal(); // Закрываем модальное окно после успешной смены пароля
    } catch (error) {
        console.error('Ошибка:', error);

        // Выделяем поле со старым паролем красным цветом
        oldPasswordInput.classList.add('error'); // Добавляем класс для красной рамки

        // Отображаем сообщение об ошибке
        const errorMessage = document.getElementById('old-password-error');
        errorMessage.textContent = "Старый пароль неверный."; // Сообщение об ошибке
        errorMessage.style.display = 'block'; // Показываем сообщение

        // Сбрасываем значения полей ввода (по желанию)
        oldPasswordInput.value = '';
    }
}

// Функция для закрытия модального окна
function closeChangePasswordModal() {
    document.getElementById('change-password-modal').style.display = 'none';

    // Сбросить стиль и сообщения об ошибках при закрытии модального окна
    document.getElementById('old-password').classList.remove('error');
    document.getElementById('new-password').classList.remove('error');
    document.getElementById('confirm-new-password').classList.remove('error');

    // Скрыть все сообщения об ошибках
    document.getElementById('old-password-error').style.display = 'none';
    document.getElementById('new-password-error').style.display = 'none';
    document.getElementById('confirm-new-password-error').style.display = 'none';
}


function openChangePasswordModal() {
    document.getElementById('change-password-modal').style.display = 'flex';
}

