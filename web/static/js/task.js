// Функция для получения ID проекта и задачи из URL
function getIdsFromUrl() {
    const pathSegments = window.location.pathname.split('/'); // Разбиваем URL на части
    const projectId = pathSegments[2]; // ID проекта находится на 3-й позиции (индекс 2)
    const taskId = pathSegments[4]; // ID задачи находится на 5-й позиции (индекс 4)
    return { projectId, taskId };
}

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
        updateTaskTitle(taskId, task.title )
        displayTaskDescription(task.description);

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