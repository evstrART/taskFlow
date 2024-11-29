// Display current year in the footer
document.getElementById('year').textContent = new Date().getFullYear();

function formatDate(isoDate) {
    if (!isoDate) return 'Not available'; // Проверка на наличие даты
    const date = new Date(isoDate);
    const options = { day: '2-digit', month: '2-digit', year: 'numeric' }; // Используем 'numeric' для полного года
    return date.toLocaleDateString('ru-RU', options).replace(/\./g, '.'); // Заменяем точки, если нужно
}

async function fetchProjects() {
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }

    try {
        console.log('Fetching projects...');  // Для отладки
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
        console.error('Error:', error);  // Для отладки
        alert(error.message);
    }
}

// Display project names in the sidebar
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

        projectItem.innerHTML = `
            <div class="project-info">
                <h3>${project.name}</h3>
                <p>${project.description || 'No description available.'}</p>
                <span class="project-date">Created on: ${formatDate(project.created_at)}</span>
            </div>
            <button class="view-button">View Project</button>
        `;

        // Добавление обработчика событий для кнопки
        const viewButton = projectItem.querySelector('.view-button');
        viewButton.addEventListener('click', () => {
            window.location.href = `/projects/${project.project_id}`;
        });

        projectsList.appendChild(projectItem);
    });
}
// Load projects on page load
window.onload = fetchProjects;

// Функция для редиректа на страницу профиля или на страницу входа
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