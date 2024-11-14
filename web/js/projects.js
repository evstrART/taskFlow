// Display current year in the footer
document.getElementById('year').textContent = new Date().getFullYear();

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
        projectItem.textContent = project.name;
        projectItem.addEventListener('click', () => {
            window.location.href = `/projects/${project.project_id}`;
        });
        projectsList.appendChild(projectItem);
    });
}

// Load projects on page load
window.onload = fetchProjects;