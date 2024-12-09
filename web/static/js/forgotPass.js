document.getElementById('year').textContent = new Date().getFullYear();

document.getElementById('reset-password-form').addEventListener('submit', async function(event) {
    event.preventDefault(); // Отменяем стандартное поведение формы

    const formData = {
        username: document.getElementById('username').value,
        email: document.getElementById('email').value
    };

    try {
        const response = await fetch('/auth/forgot-password', {
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

        alert("Ссылка для сброса пароля отправлена на вашу электронную почту.");
    } catch (error) {
        console.error('Ошибка:', error);
        alert("Ошибка при отправке запроса: " + error.message);
    }
});