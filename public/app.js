document.addEventListener('DOMContentLoaded', loadTasks);

// Obsługa formularza dodawania zadania
document.getElementById('addTaskForm').addEventListener('submit', function(event) {
  event.preventDefault();

  const title = document.getElementById('title').value;
  const dueDate = document.getElementById('dueDate').value;

  const task = {
    title: title,
    due_date: dueDate
  };

  fetch('/tasks', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(task)
  })
    .then(res => {
      loadTasks(); // Przeładuj listę zadań
      document.getElementById('addTaskForm').reset(); // Wyczyść formularz
    })
    .catch(err => console.error('Błąd przy dodawaniu zadania', err));
});

function loadTasks() {
  fetch('/tasks')
    .then(res => res.json())
    .then(tasks => {
      const container = document.getElementById('taskList');
      container.innerHTML = '';

      tasks.forEach(task => {
        if (window.location.pathname.includes('completed.html') && !task.is_completed) return;
        if (window.location.pathname.includes('index.html') && task.is_completed) return;

        const card = document.createElement('div');
        card.className = 'p-4 bg-white rounded shadow-md flex flex-col gap-2';

        const title = document.createElement('div');
        title.className = 'text-lg font-semibold';
        title.textContent = task.title;

        const date = document.createElement('div');
        date.className = 'text-sm text-gray-500';
        date.textContent = `Do wykonania: ${new Date(task.due_date).toLocaleString()}`;

        const btns = document.createElement('div');
        btns.className = 'flex gap-2';

        const completeBtn = document.createElement('button');
        completeBtn.className = 'flex-1 bg-green-500 text-white p-2 rounded';
        completeBtn.textContent = task.is_completed ? 'Cofnij' : 'Ukończ';
        completeBtn.onclick = () => toggleComplete(task.id, task.is_completed);

        const deleteBtn = document.createElement('button');
        deleteBtn.className = 'flex-1 bg-red-500 text-white p-2 rounded';
        deleteBtn.textContent = 'Usuń';
        deleteBtn.onclick = () => deleteTask(task.id);

        btns.appendChild(completeBtn);
        btns.appendChild(deleteBtn);

        card.appendChild(title);
        card.appendChild(date);
        card.appendChild(btns);

        container.appendChild(card);
      });
    });
}

function toggleComplete(id, isCompleted) {
  fetch(`/tasks/${id}/${isCompleted ? 'uncomplete' : 'complete'}`, { method: 'PUT' })
    .then(() => loadTasks());
}

function deleteTask(id) {
  fetch(`/tasks/${id}`, { method: 'DELETE' })
    .then(() => loadTasks());
}
