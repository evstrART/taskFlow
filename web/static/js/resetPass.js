document.getElementById('year').textContent = new Date().getFullYear();

// Получаем токен из URL
const urlParams = new URLSearchParams(window.location.search);
const token = urlParams.get('token');

document.getElementById('reset-password-form').addEventListener('submit', async function(event) {
    event.preventDefault(); // Отменяем стандартное поведение формы

    const newPassword = document.getElementById('new-password').value;
    const confirmPassword = document.getElementById('confirm-password').value;

    // Проверяем совпадение паролей
    if (newPassword !== confirmPassword) {
        alert("Пароли не совпадают. Пожалуйста, попробуйте снова.");
        return;
    }

    // Получаем токен из URL
    const urlParams = new URLSearchParams(window.location.search);
    const token = urlParams.get('token');

    const formData = {
        new_password: newPassword,
        token: token // Добавляем токен в данные формы
    };

    try {
        const response = await fetch('/auth/reset-password', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText);
        }

        alert("Пароль успешно сброшен.");
        window.location.href = '/auth/sign-in'; // Перенаправление на страницу входа
    } catch (error) {
        console.error('Ошибка:', error);
        alert("Ошибка при отправке запроса: " + error.message);
    }
});

