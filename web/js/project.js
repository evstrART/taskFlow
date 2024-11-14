
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
} catch (error) {
    console.error('Error:', error);
    alert(error.message);
}
}

    // Display project details
    function displayProjectDetails(project) {
    const projectTitle = document.getElementById('project-title');
    const projectDescription = document.getElementById('project-description');
    const projectNameInBreadcrumb = document.getElementById('project-name');

    projectTitle.textContent = project.name;
    projectDescription.textContent = project.description || 'No description available';
    projectNameInBreadcrumb.textContent = project.name;
}

    // Display tasks list
    function displayTasks(tasks) {
    const tasksList = document.getElementById('tasks-list');
    tasksList.innerHTML = '';

    if (tasks.length === 0) {
    tasksList.innerHTML = '<li>No tasks found</li>';
    return;
}

    tasks.forEach(task => {
    const taskItem = document.createElement('li');
    taskItem.textContent = task.title;
    tasksList.appendChild(taskItem);
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