function getProjectIdFromUrl() {
    const pathSegments = window.location.pathname.split('/'); // Разбиваем URL на части
    const projectId = pathSegments[2]; // ID проекта находится на 3-й позиции (индекс 2)
    return projectId; // Возвращаем только projectId
}

const projectId = getProjectIdFromUrl();

// Display current year in the footer
document.getElementById('year').textContent = new Date().getFullYear();

const today = new Date().toISOString().split('T')[0];
document.getElementById('projectStartDate').setAttribute('min', today);

document.getElementById('taskDueDate').setAttribute('min', today);
// Обновление минимальной даты окончания при выборе даты начала
document.getElementById('projectStartDate').addEventListener('change', function() {
    const startDate = this.value;
    document.getElementById('projectEndDate').setAttribute('min', startDate);
});
let tasks = []; // Хранит список задач
let sortOrder = {
    title: 1, // 1 для прямой сортировки, -1 для обратной
    status: 1,
    due_date: 1
};

function sortTasks(criteria) {
    tasks.sort((a, b) => {
        if (criteria === 'due_date') {
            return (new Date(a.due_date) - new Date(b.due_date)) * sortOrder[criteria];
        }
        return a[criteria].localeCompare(b[criteria]) * sortOrder[criteria];
    });

    // Переключаем порядок сортировки
    sortOrder[criteria] *= -1;

    // Обновляем отображение задач
    displayTasks(tasks);
}


// Fetch project details
async function fetchProjectDetails(projectId) {
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }

    try {
        // Fetch project details
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
        displayProjectDetails(project);

        // Fetch tasks for the project
        const tasksResponse = await fetch(`/api/projects/${projectId}/tasks`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });

        if (!tasksResponse.ok) {
            throw new Error('Error fetching tasks');
        }

        tasks = await tasksResponse.json(); // Сохраняем задачи в глобальную переменную
        displayTasks(tasks);

        // Fetch project members
        await fetchProjectMembers(projectId); // Добавлено для загрузки участников

    } catch (error) {
        console.error('Error:', error);
        alert(error.message);
    }
}
async function fetchAllUsers(){
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
function displayUsers(users){
    const addMembersSelect = document.getElementById('memberId');
    addMembersSelect.innerHTML = '<option value="">Выберите пользователя</option>'

    users.forEach(user => {
        const option = document.createElement('option');
        option.value = user.user_id; // Предполагается, что id участника доступен
        option.textContent = user.username; // Имя участника
        addMembersSelect.appendChild(option);
    });
}

// Fetch project members
async function fetchProjectMembers(projectId) {
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }

    try {
        const response = await fetch(`http://localhost:8080/api/projects/${projectId}/members`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Error fetching project members');
        }

        const members = await response.json();
        displayMembers(members);
    } catch (error) {
        console.error('Error fetching project members:', error);
        alert(error.message);
    }
}

// Display members in the team section
function displayMembers(members) {
    const teamList = document.getElementById('team-list');
    const assignedToSelect = document.getElementById('taskAssignedTo');

    // Очищаем список участников
    teamList.innerHTML = '';
    assignedToSelect.innerHTML = '<option value="">Выберите участника</option>'; // Сбрасываем предыдущие значения

    if (members.length === 0) {
        teamList.innerHTML = '<li>No members found</li>';
        return;
    }

    members.forEach(member => {
        // Создаем элемент списка
        const listItem = document.createElement('li');

        // Имя и роль участника
        listItem.textContent = `${member.username} - ${member.role}`;

        // Создаем значок корзины
        const deleteIcon = document.createElement('span');
        deleteIcon.textContent = '🗑️'; // Символ корзины
        deleteIcon.style.fontSize = '11px';
        deleteIcon.className = 'delete-icon'; // Присваиваем класс для стилей


        // Обработчик события для удаления участника
        deleteIcon.addEventListener('click', async () => {
            const confirmed = confirm(`Удалить участника ${member.username}?`);
            if (confirmed) {
                await deleteMember(member.user_id); // Вызов функции удаления
                fetchProjectMembers(currentProjectId); // Обновить список участников
            }
        });

        // Показываем корзину при наведении
        listItem.addEventListener('mouseenter', () => {
            deleteIcon.style.display = 'inline'; // Показываем значок
        });
        listItem.addEventListener('mouseleave', () => {
            deleteIcon.style.display = 'none'; // Скрываем значок
        });

        listItem.appendChild(deleteIcon); // Добавляем значок в элемент списка
        teamList.appendChild(listItem); // Добавляем элемент списка в контейнер

        // Добавляем участника в выпадающий список
        const option = document.createElement('option');
        option.value = member.user_id; // Предполагается, что id участника доступен
        option.textContent = member.username; // Имя участника
        assignedToSelect.appendChild(option);
    });
}

async function deleteMember(memberId) {
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }

    try {
        const response = await fetch(`http://localhost:8080/api/projects/${projectId}/members/${memberId}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Error deleting member');
        }

        alert('Member deleted successfully'); // Уведомление об успешном удалении
        location.reload()
    } catch (error) {
        console.error('Error deleting member:', error);
        alert(error.message);
    }
}

function formatDate(isoDate) {
    if (!isoDate) return 'Not available'; // Проверка на наличие даты
    const date = new Date(isoDate);
    const options = { day: '2-digit', month: '2-digit', year: 'numeric' }; // Используем 'numeric' для полного года
    return date.toLocaleDateString('ru-RU', options).replace(/\./g, '.'); // Заменяем точки, если нужно
}

// Display project details
function displayProjectDetails(project) {
    const projectTitle = document.getElementById('project-title');
    const projectDescription = document.getElementById('project-description');
    const projectNameInBreadcrumb = document.getElementById('project-name');
    const projectStartDate = document.getElementById('project-start-date');
    const projectEndDate = document.getElementById('project-end-date');
    const projectStatus = document.getElementById('project-status');

    projectTitle.textContent = project.name;
    projectDescription.textContent = project.description || 'No description available';
    projectNameInBreadcrumb.textContent = project.name;

    projectStartDate.textContent = formatDate(project.start_date);
    projectEndDate.textContent = formatDate(project.end_date);
    projectStatus.textContent = project.status || 'Not available';
}

// Display tasks list in table
function displayTasks(tasks) {
    const tasksList = document.getElementById('tasks-list').getElementsByTagName('tbody')[0];
    tasksList.innerHTML = ''; // Очищаем текущее содержимое

    if (!Array.isArray(tasks) || tasks.length === 0) {
        tasksList.innerHTML = '<tr><td colspan="3">Задачи не найдены</td></tr>';
        return;
    }

    tasks.forEach(task => {
        const taskRow = document.createElement('tr');

        // Создаем ячейку для названия задачи
        const taskTitleCell = document.createElement('td');
        const taskLink = document.createElement('a');
        taskLink.textContent = task.title;
        taskLink.href = `/projects/${task.project_id}/tasks/${task.task_id}`;
        taskLink.style.textDecoration = 'none'; // Убираем подчеркивание ссылки
        taskLink.style.color = 'inherit'; // Наследуем цвет текста от родителя
        taskTitleCell.appendChild(taskLink);

        // Создаем ячейку для статуса задачи
        const taskStatusCell = document.createElement('td');
        taskStatusCell.textContent = task.status || 'Не доступен';

        // Создаем ячейку для даты выполнения задачи
        const taskDueDateCell = document.createElement('td');
        taskDueDateCell.textContent = formatDate(task.due_date); // Используем функцию formatDate

        // Добавляем ячейки в строку задачи
        taskRow.appendChild(taskTitleCell);
        taskRow.appendChild(taskStatusCell);
        taskRow.appendChild(taskDueDateCell);

        // Добавляем строку в таблицу задач
        tasksList.appendChild(taskRow);
    });
}

// Display list of projects in the sidebar
async function fetchProjects() {
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }

    try {
        const response = await fetch('/api/projects/user', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Error fetching projects');
        }

        const projects = await response.json();
        displayProjects(projects);
    } catch (error) {
        console.error('Error:', error);
        alert(error.message);
    }
}

// Display projects in sidebar
function displayProjects(projects) {
    const projectsList = document.getElementById('projects-list');
    projectsList.innerHTML = '';

    if (projects.length === 0) {
        projectsList.innerHTML = '<li>No projects found</li>';
        return;
    }

    projects.forEach(project => {
        const projectItem = document.createElement('li');
        projectItem.classList.add('project-item');
        projectItem.textContent = project.name;
        projectItem.addEventListener('click', () => {
            window.location.href = `/projects/${project.project_id}`;
        });
        projectsList.appendChild(projectItem);
    });
}

// Extract project ID from URL and fetch data
window.onload = function() {
    const projectId = window.location.pathname.split('/').pop();  // Get project ID from URL
    if (projectId) {
        fetchProjectDetails(projectId);  // Load project and tasks
        fetchProjects(); // Load the list of projects in the sidebar
    } else {
        alert('Project ID is missing in the URL.');
    }
};

function openModal(modalId) {
    document.getElementById(modalId).style.display = 'flex';
}

function closeModal() {
    const modals = document.getElementsByClassName('modal');
    for (let modal of modals) {
        modal.style.display = 'none';
    }
}
function openAddMemberModal(){
    fetchAllUsers();
    openModal('addMemberModal');
}

function openAddTaskModal() {
    openModal('addTaskModal');
}
function openEditProjectModal() {
    openModal('editProjectModal');
}

// Функция для редактирования проекта
async function editProject() {
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }
    const urlParts = window.location.href.split('/');
    const idString = urlParts[urlParts.length - 1];
    const id = parseInt(idString, 10);

    if (isNaN(id)) {
        alert('Invalid project ID');
        return;
    }

    const name = document.getElementById('projectName').value;
    const description = document.getElementById('projectDescription').value;
    const startDate = document.getElementById('projectStartDate').value;
    const endDate = document.getElementById('projectEndDate').value;
    const status = document.getElementById('projectStatus').value;

    const input = {
        name: name ? name : null,
        description: description ? description : null,
        start_date: startDate ? startDate : null,
        end_date: endDate ? endDate : null,
        status: status ? status : null,
    };

    const url = `http://localhost:8080/api/projects/${id}`;

    try {
        const response = await fetch(url, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(input)
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(`Error: ${errorData.message} (Status: ${response.status})`);
        }

        const successData = await response.json();
        console.log(`Project ${id} updated successfully:`, successData);
        alert(`Project ${id} updated successfully.`);
        closeModal(); // Закрываем модальное окно
        window.location.href = `/projects/${id}`; // Перенаправление на страницу проекта
    } catch (error) {
        console.error(`Failed to update project: ${error.message}`);
        alert(`Failed to update project: ${error.message}`);
    }
}

// Добавляем обработчик события для кнопки редактирования
document.querySelector('.edit-button').onclick = openEditProjectModal;
document.querySelector('#add-button').onclick = openAddTaskModal;
document.querySelector('.add-members-button').onclick= openAddMemberModal;

// Delete Project Function
async function deleteProject(projectId) {
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }
    // Получаем ID проекта из URL
    const urlParts = window.location.href.split('/');
    const idString = urlParts[urlParts.length - 1]; // Последний элемент в URL
    const id = parseInt(idString, 10); // Преобразование в целое число

    // Проверка, является ли id действительным целым числом
    if (isNaN(id)) {
        alert('Invalid project ID');
        return;
    }

    const url = `http://localhost:8080/api/projects/${id}`;

    try {
        const response = await fetch(url, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            // Обработка ошибок
            const errorData = await response.json();
            throw new Error(`Error: ${errorData.message} (Status: ${response.status})`);
        }

        // Успешное удаление
        const successData = await response.json();
        console.log(`Project ${id} deleted successfully:`, successData);
        alert(`Project ${id} deleted successfully.`);
        // Перенаправление или обновление страницы после удаления
        window.location.href = '/projects'; // Замените на нужный вам адрес
    } catch (error) {
        // Обработка ошибки
        console.error(`Failed to delete project: ${error.message}`);
        alert(`Failed to delete project: ${error.message}`);
    }
}

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
function addMembers() {
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }

    const urlParts = window.location.href.split('/');
    const idString = urlParts[urlParts.length - 1];
    const projectId = parseInt(idString, 10); // Преобразование в целое число

    const memberId = parseInt(document.getElementById('memberId').value, 10);
    const role = document.getElementById('role').value;

    console.log("Selected Member ID:", memberId);
    console.log("Selected Role:", role);

    const memberData = {
        project_id: projectId,
        user_id: memberId,
        role: role
    };

    fetch(`http://localhost:8080/api/projects/${projectId}/members`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json', // Указываем, что отправляем JSON
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(memberData) // Исправлено на memberData
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            return response.json();
        })
        .then(data => {
            alert("Member added successfully!");
            closeModal(); // Закрываем модальное окно
            location.reload(); // Перезагружаем страницу
        })
        .catch(error => {
            console.error('There was a problem with the fetch operation:', error);
            alert("Failed to add member. Please try again.");
        });
}

function addTask() {
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }

    // Получаем ID проекта из URL
    const urlParts = window.location.href.split('/');
    const idString = urlParts[urlParts.length - 1];
    const projectId = parseInt(idString, 10); // Преобразование в целое число

    // Сбор данных из формы
    const taskTitle = document.getElementById('taskTitle').value;
    const taskDescription = document.getElementById('taskDescription').value;

    // Получаем ID пользователя и преобразуем в целое число
    const taskAssignedTo = parseInt(document.getElementById('taskAssignedTo').value, 10);

    const taskStatus = document.getElementById('taskStatus').value;
    const taskPriority = document.getElementById('taskPriority').value;
    const taskDueDate = new Date(document.getElementById('taskDueDate').value).toISOString()

    // Создание объекта задачи
    const taskData = {
        title: taskTitle,
        description: taskDescription,
        assigned_to: taskAssignedTo, // Теперь это целое число
        status: taskStatus,
        priority: taskPriority,
        due_date: taskDueDate
    };

    // Отправка POST-запроса
    fetch(`http://localhost:8080/api/projects/${projectId}/tasks`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json', // Указываем, что отправляем JSON
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(taskData) // Преобразуем объект в строку JSON
    })
        .then(response => {
            if (!response.ok) {
                return response.json().then(err => {
                    throw new Error(`Ошибка ${response.status}: ${err.message || 'Unknown error'}`);
                });
            }
            return response.json();
        })
        .then(data => {
            console.log('Task added successfully:', data);
            closeModal(); // Закрываем модальное окно
            location.reload(); // Перезагружаем страницу
        })
        .catch(error => {
            console.error('Error adding task:', error);
        });
}

async function completeProject() {
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }

    const url = `http://localhost:8080/api/projects/${projectId}/complete`; // URL для завершения проекта

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(`Error completing project: ${errorData.message}`);
        }

        alert("Project completed successfully!"); // Успешное завершение проекта
        window.location.href = "http://localhost:8080/projects";
    } catch (error) {
        console.error(`Failed to complete project: ${error.message}`);
        alert(`Failed to complete project: ${error.message}`);
    }
}

