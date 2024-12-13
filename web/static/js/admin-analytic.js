async function fetchCompletedTasksByProject() {
    const token = localStorage.getItem('token');
    if (!token) {
        alert("No token found. Please log in.");
        return;
    }
    try {
        const response = await fetch('/api/completed-tasks-by-project', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        });

        if (!response.ok) throw new Error('Network response was not ok');

        const data = await response.json();

        console.log('Data received for completed tasks by project:', data);

        if (!data.projects) throw new Error('No projects data found');

        const projects = data.projects.map(project => project.project_name);
        const completedTasks = data.projects.map(project => project.completed_tasks);

        createCompletedTasksChart(projects, completedTasks);
    } catch (error) {
        console.error('Ошибка при получении данных завершенных задач по проектам:', error);
    }
}

async function fetchCreatedTasksByUser() {
    const token = localStorage.getItem('token');
    if (!token) {
        alert("No token found. Please log in.");
        return;
    }
    try {
        const response = await fetch('/api/created-tasks-by-user', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        });

        // Логируем объект response
        console.log('Response:', response);

        if (!response.ok) throw new Error('Network response was not ok');

        // Преобразуем ответ в JSON
        const data = await response.json();

        // Логируем данные после преобразования в JSON
        console.log('Data received for created tasks by user:', data);

        if (!data.users) throw new Error('No users data found');

        const users = data.users.map(user => user.username);
        const createdTasks = data.users.map(user => user.created_tasks);

        createCreatedTasksChart(users, createdTasks);
    } catch (error) {
        console.error('Ошибка при получении данных созданных задач по пользователям:', error);
    }
}



function createCompletedTasksChart(projects, completedTasks) {
    const ctx = document.getElementById('completedTasksChart').getContext('2d');
    new Chart(ctx, {
        type: 'bar', // Тип графика: столбчатый
        data: {
            labels: projects, // Метки по оси X (названия проектов)
            datasets: [{
                label: 'Завершенные задачи', // Подпись для легенды графика
                data: completedTasks, // Данные по оси Y (количество завершенных задач)
                backgroundColor: 'rgba(75, 192, 192, 0.2)', // Цвет фона столбцов
                borderColor: 'rgba(75, 192, 192, 1)', // Цвет границы столбцов
                borderWidth: 1 // Толщина границы столбцов
            }]
        },
        options: {
            scales: {
                y: {
                    beginAtZero: true // Начинать ось Y с нуля
                }
            }
        }
    });
}
function createCreatedTasksChart(users, createdTasks) {
    const ctx = document.getElementById('createdTasksChart').getContext('2d');
    new Chart(ctx, {
        type: 'bar', // Тип графика: столбчатый
        data: {
            labels: users, // Метки по оси X (имена пользователей)
            datasets: [{
                label: 'Созданные задачи', // Подпись для легенды графика
                data: createdTasks, // Данные по оси Y (количество созданных задач)
                backgroundColor: 'rgba(153, 102, 255, 0.2)', // Цвет фона столбцов
                borderColor: 'rgba(153, 102, 255, 1)', // Цвет границы столбцов
                borderWidth: 1 // Толщина границы столбцов
            }]
        },
        options: {
            scales: {
                y: {
                    beginAtZero: true // Начинать ось Y с нуля
                }
            }
        }
    });
}
