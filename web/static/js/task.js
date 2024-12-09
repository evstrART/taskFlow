// Функция для получения ID проекта и задачи из URL
function getIdsFromUrl() {
    const pathSegments = window.location.pathname.split('/'); // Разбиваем URL на части
    const projectId = pathSegments[2]; // ID проекта находится на 3-й позиции (индекс 2)
    const taskId = pathSegments[4]; // ID задачи находится на 5-й позиции (индекс 4)
    return { projectId, taskId };
}
const {projectId, taskId} = getIdsFromUrl()

const today = new Date().toISOString().split('T')[0];
document.getElementById('taskDueDate').setAttribute('min', today);

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

function openTagModal() {
    document.getElementById('tag-modal').style.display = 'block';
}

// Функция для закрытия модального окна для тегов
function closeTagModal() {
    document.getElementById('tag-modal').style.display = 'none';
}
document.querySelector('.add-tag').addEventListener('click', openTagModal);

// Закрытие модального окна для тегов
    document.getElementById('close-tag-modal').addEventListener('click', closeTagModal);

// Закрытие модального окна при клике вне его
window.addEventListener('click', function(event) {
    const modal = document.getElementById('tag-modal');
    if (event.target === modal) {
        modal.style.display = 'none'; // Закрываем модальное окно
    }
});
document.getElementById('save-tag').addEventListener('click', function() {
    const tagInput = document.getElementById('tag-input');
    const tagColor = document.querySelector('.color.selected'); // Получаем выбранный цвет
    if (tagInput.value && tagColor) {
        const tagData = {
            name: tagInput.value,
            color: tagColor.getAttribute('data-color'),
        };

        // Используем глобальные переменные projectId и taskId
        fetch(`/api/projects/${projectId}/tasks/${taskId}/tags`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}` // Если требуется авторизация
            },
            body: JSON.stringify(tagData),
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Ошибка при добавлении тега: ' + response.statusText);
                }
                return response.json();
            })
            .then(data => {
                console.log('Response data:', data); // Проверяем, что возвращает сервер

                const tag = document.createElement('div');
                tag.style.backgroundColor = tagData.color; // Устанавливаем цвет фона
                tag.textContent = tagData.name;
                tag.className = 'tag'; // Добавляем класс для стилей
                tag.setAttribute('data-tag-id', data.tag_id); // Сохраняем ID тега в атрибуте

                document.querySelector('.task-header').appendChild(tag); // Добавляем тег в заголовок задачи

                // Создаем элемент крестика для удаления
                const deleteIcon = document.createElement('span');
                deleteIcon.textContent = '✖'; // Символ крестика
                deleteIcon.className = 'delete-icon'; // Класс для стилей

                // Обработчик события на крестик
                deleteIcon.addEventListener('click', function(event) {
                    event.stopPropagation(); // Останавливаем событие, чтобы не срабатывало на родителе
                    deleteTag(data.tag_id, tag); // Удаляем тег
                });

                tag.appendChild(deleteIcon); // Добавляем крестик в тег
                document.querySelector('.tag-container').appendChild(tag); // Добавляем тег в контейнер
                tagInput.value = ''; // Очищаем поле ввода
                tagColor.classList.remove('selected'); // Сбрасываем выбор цвета
                document.getElementById('tag-modal').style.display = 'none';// Закрываем модальное окно
                loadTags();
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
    }
});

// Обработчик для выбора цвета
document.querySelectorAll('.color').forEach(color => {
    color.addEventListener('click', function() {
        document.querySelectorAll('.color').forEach(c => c.classList.remove('selected')); // Сбрасываем выделение
        color.classList.add('selected'); // Выделяем выбранный цвет
    });
});

// Функция для получения и отображения тегов при загрузке страницы
function loadTags() {
    fetch(`/api/projects/${projectId}/tasks/${taskId}/tags`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}` // Если требуется авторизация
        }
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Ошибка при загрузке тегов: ' + response.statusText);
            }
            return response.json();
        })
        .then(tags => {
            const tagContainer = document.querySelector('.tag-container');
            tagContainer.innerHTML = ''; // Очищаем контейнер перед добавлением тегов

            tags.forEach(tag => {
                const tagElement = document.createElement('div');
                tagElement.style.backgroundColor = tag.color; // Устанавливаем цвет фона тега
                tagElement.textContent = tag.name;
                tagElement.className = 'tag'; // Добавляем класс для стилей

                // Создаем элемент крестика для удаления
                const deleteIcon = document.createElement('span');
                deleteIcon.textContent = '✖'; // Символ крестика
                deleteIcon.className = 'delete-icon'; // Класс для стилей

                // Обработчик события на крестик
                deleteIcon.addEventListener('click', function(event) {
                    event.stopPropagation(); // Останавливаем событие, чтобы не срабатывало на родителе
                    deleteTag(tag.tag_id, tagElement); // Удаляем тег
                });

                tagElement.appendChild(deleteIcon); // Добавляем крестик в тег
                tagContainer.appendChild(tagElement); // Добавляем тег в контейнер
            });
        })
        .catch(error => {
            console.error('Ошибка:', error);
        });
}

// Функция для удаления тега
function deleteTag(tag_id, tagElement) {
    console.log(`Удаление тега с ID: ${tag_id}`);
    fetch(`/api/projects/${projectId}/tasks/${taskId}/tags/${tag_id}`, {
        method: 'DELETE',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}` // Если требуется авторизация
        }
    })
        .then(response => {
            console.log('Ответ сервера:', response);
            if (!response.ok) {
                throw new Error('Ошибка при удалении тега: ' + response.statusText);
            }
            tagElement.remove(); // Удаляем элемент из DOM
        })
        .catch(error => {
            console.error('Ошибка:', error);
        });
}

// Вызываем функцию для загрузки тегов при загрузке страницы
document.addEventListener('DOMContentLoaded', loadTags);

function openEditTaskModal() {
    document.getElementById('editTaskModal').style.display = 'block';
}

function closeEditTaskModal() {
    document.getElementById('editTaskModal').style.display = 'none';
}

async function editTask() {
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }

    // Получаем данные задачи
    const taskUrl = `http://localhost:8080/api/projects/${projectId}/tasks/${taskId}`; // Замените projectId на актуальный ID проекта

    try {
        const response = await fetch(taskUrl, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(`Error fetching task: ${errorData.message}`);
        }

        const taskData = await response.json();
        document.getElementById('taskTitle').value = taskData.title || '';
        document.getElementById('taskDescription').value = taskData.description || '';
        document.getElementById('taskStatus').value = taskData.status || '';
        document.getElementById('taskPriority').value = taskData.priority || '';
        document.getElementById('taskDueDate').value = taskData.dueDate || '';

        // Заполнение выпадающего списка "Assigned To"
        await loadAssignedToOptions(taskData.assignedTo); // Загружаем исполнителей

        openEditTaskModal(); // Открываем модальное окно
    } catch (error) {
        console.error(`Failed to fetch task: ${error.message}`);
        alert(`Failed to fetch task: ${error.message}`);
    }
}
async function loadAssignedToOptions(selectedId) {
    const token = localStorage.getItem('token');
    const usersUrl = 'http://localhost:8080/api/users'; // URL для получения списка пользователей

    try {
        const response = await fetch(usersUrl, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Failed to fetch users');
        }

        const users = await response.json();
        const assignedToSelect = document.getElementById('assignedTo');
        assignedToSelect.innerHTML = ''; // Очистка предыдущих опций

        users.forEach(user => {
            const option = document.createElement('option');
            option.value = user.id; // Предполагается, что у пользователя есть поле `id`
            option.textContent = user.username; // Предполагается, что у пользователя есть поле `username`
            if (user.id === selectedId) {
                option.selected = true; // Устанавливаем выбранным текущего исполнителя
            }
            assignedToSelect.appendChild(option);
        });
    } catch (error) {
        console.error(`Failed to load users: ${error.message}`);
        alert(`Failed to load users: ${error.message}`);
    }
}

async function updateTask() {
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }

    const input = {
        title: document.getElementById('taskTitle').value || null,
        description: document.getElementById('taskDescription').value || null,
        assigned_to: parseInt(document.getElementById('assignedTo').value) || null,
        status: document.getElementById('taskStatus').value || null,
        priority: document.getElementById('taskPriority').value || null,
        due_date: document.getElementById('taskDueDate').value || null
    };

    const taskUrl = `http://localhost:8080/api/projects/${projectId}/tasks/${taskId}`; // Замените projectId на актуальный ID проекта

    try {
        const response = await fetch(taskUrl, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(input)
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(`Error updating task: ${errorData.message}`);
        }

        alert(`Task ${taskId} updated successfully.`);
        closeEditTaskModal(); // Закрываем модальное окно
        // Обновление страницы или перенаправление, если необходимо
        window.location.reload(); // Перезагружаем страницу, чтобы увидеть изменения
    } catch (error) {
        console.error(`Failed to update task: ${error.message}`);
        alert(`Failed to update task: ${error.message}`);
    }
}

async function deleteTask() {
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }

    const taskUrl = `http://localhost:8080/api/projects/${projectId}/tasks/${taskId}`; // Замените projectId на актуальный ID проекта

    if (confirm("Вы уверены, что хотите удалить эту задачу?")) {
        try {
            const response = await fetch(taskUrl, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(`Error deleting task: ${errorData.message}`);
            }

            alert(`Task ${taskId} deleted successfully.`);
            // Обновление страницы или перенаправление, если необходимо
            window.location.href = `http://localhost:8080/projects/${projectId}/`;
        } catch (error) {
            console.error(`Failed to delete task: ${error.message}`);
            alert(`Failed to delete task: ${error.message}`);
        }
    }
}