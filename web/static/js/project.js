// Display current year in the footer
document.getElementById('year').textContent = new Date().getFullYear();

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

        const tasks = await tasksResponse.json();
        displayTasks(tasks);

        // Fetch project members
        await fetchProjectMembers(projectId); // Добавлено для загрузки участников

    } catch (error) {
        console.error('Error:', error);
        alert(error.message);
    }
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
    teamList.innerHTML = ''; // Очищаем список перед добавлением новых участников

    if (members.length === 0) {
        teamList.innerHTML = '<li>No members found</li>';
        return;
    }

    members.forEach(member => {
        const listItem = document.createElement('li');
        listItem.textContent = `${member.username} - ${member.role}`; // Выводим имя и роль
        teamList.appendChild(listItem);
    });
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

    projectStartDate.textContent = project.start_date || 'Not available';
    projectEndDate.textContent = project.end_date || 'Not available';
    projectStatus.textContent = project.status || 'Not available';
}

// Display tasks list in table
function displayTasks(tasks) {
    const tasksList = document.getElementById('tasks-list').getElementsByTagName('tbody')[0];
    tasksList.innerHTML = '';

    if (tasks.length === 0) {
        tasksList.innerHTML = '<tr><td colspan="3">No tasks found</td></tr>';
        return;
    }

    tasks.forEach(task => {
        const taskRow = document.createElement('tr');

        const taskTitleCell = document.createElement('td');
        const taskLink = document.createElement('a');
        taskLink.textContent = task.title;
        taskLink.href = `/projects/${task.project_id}/tasks/${task.task_id}`;
        taskTitleCell.appendChild(taskLink);

        const taskStatusCell = document.createElement('td');
        taskStatusCell.textContent = task.status || 'Not available';

        const taskDueDateCell = document.createElement('td');
        taskDueDateCell.textContent = task.due_date || 'Not available';

        taskRow.appendChild(taskTitleCell);
        taskRow.appendChild(taskStatusCell);
        taskRow.appendChild(taskDueDateCell);

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
        const response = await fetch('/api/projects', {
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

// Edit Project Function
// Функция для открытия модального окна
function openModal() {
    document.getElementById('editProjectModal').style.display = 'flex';
}

// Функция для закрытия модального окна
function closeModal() {
    document.getElementById('editProjectModal').style.display = 'none';
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
document.querySelector('.edit-button').onclick = openModal;

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
function addMembers(){
    alert("add members")
}

function addTask(){
    alert("add task")
}

