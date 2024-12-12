let currentUserId = null; // Переменная для хранения ID текущего пользователя

// Функция для получения всех пользователей
async function fetchAllUsers() {
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }
    try {
        const response = await fetch(`http://localhost:8080/api/users`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Error fetching users');
        }

        const users = await response.json();
        displayUsers(users);
    } catch (error) {
        console.error('Error fetching users:', error);
        alert(error.message);
    }
}

// Функция для отображения пользователей в таблице
function displayUsers(users) {
    const tbody = document.querySelector('#users-table tbody');
    tbody.innerHTML = ''; // Очистить предыдущие данные

    users.forEach(user => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${user.username}</td>
            <td>${user.email}</td>
            <td>${user.role}</td>
            <td>
                <i class='fas fa-edit' onclick='openEditModal(${user.user_id}, "${user.username}", "${user.email}", "${user.role}")'></i> 
                | 
                <i class='fas fa-trash' onclick='openDeleteConfirmation(${user.user_id})'></i>
            </td>`;
        tbody.appendChild(row);
    });
}

// Функция для открытия модального окна редактирования
function openEditModal(userId, username, email, role) {
    currentUserId = userId; // Сохраняем ID пользователя
    document.getElementById('edit-name').value = username;
    document.getElementById('edit-email').value = email;
    document.getElementById('edit-role').value = role;

    document.getElementById('edit-modal').style.display = 'block'; // Показываем модальное окно
}

// Функция для сохранения изменений пользователя
async function saveChanges() {
    const name = document.getElementById('edit-name').value;
    const email = document.getElementById('edit-email').value;
    const role = document.getElementById('edit-role').value;

    // Проверка на пустые поля
    if (!name || !email || !role) {
        alert("Пожалуйста, заполните все поля.");
        return;
    }

    const token = localStorage.getItem('token');

    if (!token) {
        alert("Токен не найден. Пожалуйста, войдите в систему.");
        return;
    }

    console.log("User ID:", currentUserId); // Выводим ID пользователя для отладки

    try {
        const response = await fetch(`http://localhost:8080/api/users/${currentUserId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ username: name, email: email, role: role }) // Отправляем данные для обновления
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Ошибка при обновлении пользователя: ${errorText}`);
        }

        alert("Информация о пользователе успешно обновлена!");
        closeModal(); // Закрываем модальное окно после успешного обновления
        fetchAllUsers(); // Обновляем список пользователей после изменения
    } catch (error) {
        console.error('Ошибка:', error);
        alert(error.message);
    }
}

// Функция для открытия подтверждения удаления
function openDeleteConfirmation(userId) {
    currentUserId = userId; // Сохраняем ID пользователя для удаления
    document.getElementById('delete-confirmation').style.display = 'block'; // Показываем окно подтверждения
}

// Функция для подтверждения удаления пользователя
async function confirmDelete() {
    const token = localStorage.getItem('token');

    if (!token) {
        alert("Токен не найден. Пожалуйста, войдите в систему.");
        return;
    }

    try {
        const response = await fetch(`http://localhost:8080/api/users/${currentUserId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Ошибка при удалении пользователя: ${errorText}`);
        }

        alert("Пользователь успешно удален!");
        closeDeleteConfirmation(); // Закрываем окно подтверждения после успешного удаления
        fetchAllUsers(); // Обновляем список пользователей после удаления
    } catch (error) {
        console.error('Ошибка:', error);
        alert(error.message);
    }
}

// Функция для закрытия модального окна редактирования
function closeModal() {
    document.getElementById('edit-modal').style.display = 'none'; // Скрываем модальное окно
}

// Функция для закрытия окна подтверждения удаления
function closeDeleteConfirmation() {
    document.getElementById('delete-confirmation').style.display = 'none'; // Скрываем окно подтверждения
}

// Вызов функции для получения пользователей при загрузке страницы
document.addEventListener('DOMContentLoaded', fetchAllUsers);
