// Display current year in the footer
document.getElementById('year').textContent = new Date().getFullYear();

const today = new Date().toISOString().split('T')[0];
document.getElementById('startDate').setAttribute('min', today);
document.getElementById('endDate').setAttribute('min', today);

document.getElementById('startDate').addEventListener('change', function() {
    const startDate = this.value;
    document.getElementById('endDate').setAttribute('min', startDate);
});

function formatDate(isoDate) {
    if (!isoDate) return 'Not available'; // Проверка на наличие даты
    const date = new Date(isoDate);
    const options = { day: '2-digit', month: '2-digit', year: 'numeric' }; // Используем 'numeric' для полного года
    return date.toLocaleDateString('ru-RU', options).replace(/\./g, '.'); // Заменяем точки, если нужно
}
let allProjects = []; // Хранит все проекты для поиска

async function fetchProjects() {
    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }

    try {
        console.log('Fetching projects...');  // Для отладки
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

        allProjects = await response.json(); // Сохраняем все проекты
        displayProjects(allProjects); // Отображаем все проекты
    } catch (error) {
        console.error('Error:', error);  // Для отладки
        alert(error.message);
    }
}

function searchProjects() {
    const searchInput = document.getElementById('searchInput').value.toLowerCase(); // Получаем текст поиска
    const filteredProjects = allProjects.filter(project => {
        // Проверяем наличие свойств title и description
        const title = project.name ? project.name.toLowerCase() : '';
        const description = project.description ? project.description.toLowerCase() : '';

        return title.includes(searchInput) || description.includes(searchInput); // Фильтруем по заголовку и описанию
    });

    displayProjects(filteredProjects); // Отображаем отфильтрованные проекты
}

// Добавляем обработчик события для кнопки поиска
document.getElementById('searchButton').addEventListener('click', searchProjects);

// Также можно добавить обработчик события для ввода текста
document.getElementById('searchInput').addEventListener('keyup', (event) => {
    if (event.key === 'Enter') {
        searchProjects(); // Запускаем поиск при нажатии Enter
    }
});

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
            <button class="view-button">Посмотреть проект</button>
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

// Открытие модального окна
function openModal(modalId) {
    document.getElementById(modalId).style.display = 'flex';
}

// Закрытие модального окна
function closeModal() {
    const modals = document.getElementsByClassName('modal');
    for (let modal of modals) {
        modal.style.display = 'none';
    }
}

// Открытие модального окна для добавления проекта
const addProjectButton = document.querySelector('button[type="add-project"]');
addProjectButton.onclick = function() {
    openModal('addProjectModal');
}

// Закрытие модального окна при нажатии на крестик
const closeModalButton = document.getElementById('closeModal');
closeModalButton.onclick = closeModal;

// Закрытие модального окна при клике вне содержимого
window.onclick = function(event) {
    const modal = document.getElementById('addProjectModal');
    if (event.target === modal) {
        closeModal();
    }
}
// Обработка отправки формы
document.getElementById('addProjectForm').onsubmit = async function(event) {
    event.preventDefault(); // Предотвращаем перезагрузку страницы

    const token = localStorage.getItem('token');
    const name = document.getElementById('projectName').value;
    const description = document.getElementById('projectDescription').value;
    const status = document.getElementById('projectStatus').value;
    const startDate = document.getElementById('startDate').value + 'T00:00:00Z'; // Добавляем временную метку
    const endDate = document.getElementById('endDate').value + 'T00:00:00Z'; // Добавляем временную метку

    try {
        const response = await fetch('/api/projects', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                name: name,
                description: description,
                status: status,
                start_date: startDate, // Добавляем дату начала
                end_date: endDate      // Добавляем дату окончания
            })
        });

        if (!response.ok) {
            const errorData = await response.json(); // Получаем данные об ошибке
            throw new Error(errorData.message || 'Failed to create project');
        }

        closeModal(); // Закрываем модальное окно
        document.getElementById('addProjectForm').reset(); // Очищаем форму
        fetchProjects(); // Обновляем список проектов
    } catch (error) {
        console.error('Error:', error);
        alert(error.message);
    }
}