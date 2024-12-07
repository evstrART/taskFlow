// Функция для получения ID проекта и задачи из URL
function getIdsFromUrl() {
    const pathSegments = window.location.pathname.split('/'); // Разбиваем URL на части
    const projectId = pathSegments[2]; // ID проекта находится на 3-й позиции (индекс 2)
    const taskId = pathSegments[4]; // ID задачи находится на 5-й позиции (индекс 4)
    return { projectId, taskId };
}
const {projectId, taskId} = getIdsFromUrl()

// Функция для обновления Breadcrumb
function updateBreadcrumb(projectId, projectTitle, taskTitle) {
    const breadcrumb = document.querySelector('.breadcrumb ul');
    breadcrumb.innerHTML = `
        <li><a href="/">Home</a></li>
        <li><a href="/projects">Projects</a></li>
        <li><a href="/projects/${projectId}">${projectTitle}</a></li>
        <li><span aria-current="page">${taskTitle}</span></li>
    `;
}

// Функция для получения деталей проекта и обновления Breadcrumb
async function fetchProjectAndTask() {
    const { projectId, taskId } = getIdsFromUrl();
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }

    try {
        // Запрос для получения деталей проекта
        const projectResponse = await fetch(`/api/projects/${projectId}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });

        if (!projectResponse.ok) {
            throw new Error('Error fetching project details');
        }

        const project = await projectResponse.json();

        const taskResponse = await fetch(`/api/projects/${projectId}/tasks/${taskId}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });
        if (!taskResponse.ok) {
            throw new Error('Error fetching task details');
        }
        const task = await taskResponse.json();

        updateBreadcrumb(projectId, project.name, task.title);
        updateTaskTitle(taskId, task.title );
        displayTaskDescription(task.description);
        displayStatus(task.status);
        fetchUser(task.assigned_to);


    } catch (error) {
        console.error('Error:', error);
        alert(error.message);
    }
}

// Вызов функции при загрузке страницы
window.onload = function() {
    fetchProjectAndTask(); // Загружаем проект и обновляем Breadcrumb
};

function redirectToProfile() {
    const token = localStorage.getItem('token');
    if (token) {
        // Если есть токен, перенаправить на страницу профиля
        window.location.href = '/profile';  // Замените на актуальный путь к профилю
    } else {
        // Если нет токена, перенаправить на страницу входа
        window.location.href = '/auth/sign-in';
    }
}
function updateTaskTitle(taskId, taskTitle) {
    const titleElement = document.getElementById('task-title');
    titleElement.textContent = `#${taskId} ${taskTitle}`; // Обновляем текст заголовка
}
function displayTaskDescription(taskDescription) {
    const taskDescriptionElement = document.querySelector('.task-description p'); // Находим элемент <p> внутри .task-description
    if (taskDescriptionElement) {
        taskDescriptionElement.textContent = taskDescription; // Устанавливаем текст описания задачи
    } else {
        console.error('Element with class "task-description" not found.');
    }
}

function loadComments() {
    const { projectId, taskId } = getIdsFromUrl(); // Получаем ID из URL
    const url = `/api/projects/${projectId}/tasks/${taskId}/comments`;
    const token = localStorage.getItem('token'); // Получаем токен

    if (!token) {
        console.error("No token found. Cannot load comments.");
        return;
    }

    console.log("Loading comments from:", url); // Логируем URL

    fetch(url, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}` // Включаем токен
        }
    })
        .then(response => {
            if (!response.ok) {
                throw new Error(`Error fetching comments: ${response.status} ${response.statusText}`);
            }
            return response.json(); // Преобразуем ответ в JSON
        })
        .then(data => {
            const commentList = document.getElementById('comment-list');
            commentList.innerHTML = ''; // Очищаем предыдущие комментарии

            if (data.length === 0) {
                commentList.innerHTML = '<p>Нет комментариев.</p>'; // Сообщение, если нет комментариев
                return;
            }

            data.forEach(comment => {
                const commentDiv = document.createElement('div');
                commentDiv.className = 'comment';
                commentDiv.dataset.commentId = comment.comment_id;

                // Добавляем содержимое комментария с учетом структуры Comment
                commentDiv.innerHTML = `
                <div class="comment-avatar">🙂</div>
                <div class="comment-content">
                    <p>${comment.content}</p>
                    <div class="comment-meta">
                        <div class="comment-username">${comment.username}</div>
                        <div class="comment-time">${new Date(comment.created_at).toLocaleString()}</div>   
                    </div>
                </div>
                <div class="comment-actions">
                    <button class="edit-comment" title="Edit">✏️</button>
                    <button class="delete-comment" title="Delete">🗑️</button>
                </div>
            `;

                commentList.appendChild(commentDiv);
            });
            document.querySelectorAll('.edit-comment').forEach(button => {
                button.addEventListener('click', handleEditComment);
            });

            document.querySelectorAll('.delete-comment').forEach(button => {
                button.addEventListener('click', handleDeleteComment);
            });
        })
        .catch(error => {
            console.error('Ошибка при загрузке комментариев:', error);
        });
}
function handleEditComment(event) {
    const commentDiv = event.target.closest('.comment');
    const commentId = commentDiv.dataset.commentId;
    const newContent = prompt("Введите новый текст комментария:");

    if (newContent) {
        const url = `/api/projects/${projectId}/tasks/${taskId}/comments/${commentId}`;
        const token = localStorage.getItem('token');

        fetch(url, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ comment: newContent }) // Форматируем данные для отправки
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Ошибка при редактировании комментария: ' + response.statusText);
                }
                return response.json();
            })
            .then(() => {
                loadComments(); // Обновляем комментарии после редактирования
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
    }
}

function handleDeleteComment(event) {
    const commentDiv = event.target.closest('.comment');
    const commentId = commentDiv.dataset.commentId;
    const url = `/api/projects/${projectId}/tasks/${taskId}/comments/${commentId}`;
    const token = localStorage.getItem('token');

    if (confirm("Вы уверены, что хотите удалить этот комментарий?")) {
        fetch(url, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Ошибка при удалении комментария: ' + response.statusText);
                }
                loadComments(); // Обновляем список комментариев после удаления
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
    }
}
// Вызов функции загрузки комментариев при загрузке страницы
document.addEventListener('DOMContentLoaded', loadComments);

// Обработчик события для отправки нового комментария
document.getElementById('submit-comment').addEventListener('click', function() {
    const commentText = document.getElementById('comment-text').value; // Получаем текст комментария
    if (commentText) {
        const url = `/api/projects/${projectId}/tasks/${taskId}/comments`; // URL для добавления комментария
        const token = localStorage.getItem('token'); // Получаем токен

        if (!token) {
            console.error("No token found. Cannot send comment.");
            return;
        }

        const commentInput = { comment: commentText }; // Создаем объект для отправки

        fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json', // Указываем тип содержимого
                'Authorization': `Bearer ${token}` // Включаем токен
            },
            body: JSON.stringify(commentInput) // Преобразуем объект в JSON
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Ошибка при отправке комментария: ' + response.statusText);
                }
                return response.json(); // Преобразуем ответ в JSON
            })
            .then(data => {
                console.log("Комментарий добавлен с ID:", data.id); // Логируем ID нового комментария
                loadComments(); // Обновляем список комментариев после успешной отправки
                document.getElementById('comment-text').value = ''; // Очистить поле ввода
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
    } else {
        alert("Пожалуйста, введите комментарий."); // Предупреждение, если поле пустое
    }
});


// Вызов функции загрузки комментариев при загрузке страницы
document.addEventListener('DOMContentLoaded', loadComments);

async function fetchUser(user_id){
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }
    try {
        const response = await fetch(`http://localhost:8080/api/users/${user_id}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Error fetching users');
        }

        const user = await response.json();
        displayUser(user);
    } catch (error) {
        console.error('Error fetching users:', error);
        alert(error.message);
    }

}
function displayUser(user){
    const assignedUserElement = document.getElementById('assigned-user'); // Находим элемент по ID
    if (assignedUserElement) {
        assignedUserElement.textContent = user.username; // Обновляем текст элемента
    } else {
        console.error('Элемент с ID "assigned-user" не найден.');
    }
}

function displayStatus(taskStatus) {
    const statusElement = document.getElementById("task-status");
    if (statusElement) {
        statusElement.textContent = taskStatus; // Обновляем текст элемента
    } else {
        console.error('Статус не найден');
    }
}