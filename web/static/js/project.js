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
function editProject() {
    alert("Edit functionality is not yet implemented.");
}

// Delete Project Function
function deleteProject() {
    alert("Delete functionality is not yet implemented.");
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

