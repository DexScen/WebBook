document.addEventListener('DOMContentLoaded', () => {
    updateAuthButtons();
    fetchBooks();
});

const booksPerPage = 10;
let allBooks = [];
let currentPage = 1;

// Функция для получения роли пользователя
function getUserRole() {
    return localStorage.getItem('userRole') || 'guest';
}

function fetchBooks() {
    const url = "http://localhost:8080/books";
    fetch(url)
        .then(response => response.json())
        .then(data => {
            allBooks = data;  
            renderBooks(currentPage);  
            updatePagination();  
        })
        .catch(error => console.error("Ошибка загрузки книг:", error));
}

function renderBooks(page) {
    const tableBody = document.getElementById('book-table-body');
    tableBody.innerHTML = '';

    const start = (page - 1) * booksPerPage;
    const end = start + booksPerPage;
    const booksToDisplay = allBooks.slice(start, end);
    const userRole = getUserRole();

    booksToDisplay.forEach(book => {
        const row = document.createElement('tr');
        
        // Базовая информация о книге (доступна всем)
        row.innerHTML = `
            <td>${book.title}</td>
            <td>${book.author}</td>
            <td>${book.price}</td>
            <td></td>
        `;

        // Добавляем кнопки действий в зависимости от роли
        if (userRole === 'admin') {
            const actionsCell = row.querySelector('td:last-child');
            
            // Создаем кнопку редактирования
            const editButton = document.createElement('button');
            editButton.textContent = 'Изменить';
            editButton.addEventListener('click', function() {
                if (getUserRole() !== 'admin') {
                    alert('У вас нет прав для редактирования книг');
                    return;
                }
                window.location.href = `edit-book.html?id=${book.id}&price=${encodeURIComponent(book.price)}&author=${encodeURIComponent(book.author)}&title=${encodeURIComponent(book.title)}`;
            });
            actionsCell.appendChild(editButton);

            // Создаем кнопку удаления
            const deleteButton = document.createElement('button');
            deleteButton.textContent = 'Удалить';
            deleteButton.addEventListener('click', async function() {
                if (getUserRole() !== 'admin') {
                    alert('У вас нет прав для удаления книг');
                    return;
                }

                const confirmDelete = confirm('Вы уверены, что хотите удалить эту книгу?');
                if (confirmDelete) {
                    try {
                        const response = await fetch(`http://localhost:8080/books`, {
                            method: 'DELETE',
                            headers: {
                                'Content-Type': 'application/json',
                            },
                            body: JSON.stringify({ id: book.id }),
                        });

                        if (response.ok) {
                            alert('Книга удалена');
                            fetchBooks();
                        } else {
                            throw new Error('Не удалось удалить книгу');
                        }
                    } catch (error) {
                        console.error('Ошибка:', error);
                        alert('Ошибка при удалении книги');
                    }
                }
            });
            actionsCell.appendChild(deleteButton);
        }

        tableBody.appendChild(row);
    });

    // Обновляем видимость формы добавления книг
    const addBookForm = document.getElementById('add-book-form');
    if (addBookForm) {
        addBookForm.style.display = (userRole !== 'guest') ? 'block' : 'none';
    }

    // Обновляем пагинацию
    updatePagination();
}


function updatePagination() {
    const paginationContainer = document.getElementById('pagination');
    paginationContainer.innerHTML = '';

    const totalPages = Math.ceil(allBooks.length / booksPerPage);

    if (totalPages > 1) {  // Показываем пагинацию только если есть больше одной страницы
        if (currentPage > 1) {
            const prevPage = document.createElement('button');
            prevPage.innerText = 'Назад';
            prevPage.onclick = () => changePage(currentPage - 1);
            paginationContainer.appendChild(prevPage);
        }

        const pageNumber = document.createElement('span');
        pageNumber.innerText = ` Страница ${currentPage} из ${totalPages} `;
        pageNumber.style.margin = '0 10px';  // Добавим отступы для лучшей читаемости
        paginationContainer.appendChild(pageNumber);

        if (currentPage < totalPages) {
            const nextPage = document.createElement('button');
            nextPage.innerText = 'Вперёд';
            nextPage.onclick = () => changePage(currentPage + 1);
            paginationContainer.appendChild(nextPage);
        }
    } else {
        // Даже если всего одна страница, все равно показываем информацию
        const pageNumber = document.createElement('span');
        pageNumber.innerText = `Страница ${currentPage} из ${totalPages || 1}`;
        paginationContainer.appendChild(pageNumber);
    }
}


function changePage(page) {
    if (page < 1 || page > Math.ceil(allBooks.length / booksPerPage)) return;
    currentPage = page;
    renderBooks(currentPage);  
    updatePagination();  
}

document.getElementById('add-book-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const userRole = getUserRole();
    if (userRole === 'guest' || !['user', 'admin'].includes(userRole)) {
        alert('У вас нет прав для добавления книг');
        return;
    }

    const title = document.getElementById('title').value;
    const author = document.getElementById('author').value;
    const price = parseInt(document.getElementById('price').value, 10);

    const newBook = {
        title: title,
        author: author,
        price: price
    };

    try {
        const response = await fetch('http://localhost:8080/books', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(newBook),
        });

        if (response.ok) {
            alert('Книга добавлена');
            fetchBooks();  
            document.getElementById('add-book-form').reset();  
        } else {
            throw new Error('Не удалось добавить книгу');
        }
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Ошибка при добавлении книги');
    }
});


// Обновление кнопок авторизации
function updateAuthButtons() {
    const buttonsBody = document.getElementById('buttons');
    buttonsBody.innerHTML = '';

    const userRole = getUserRole();
    
    if (userRole === 'guest') {
        // Кнопки для неавторизованного пользователя
        const registerButton = document.createElement('button');
        registerButton.textContent = 'Регистрация';
        registerButton.onclick = () => window.location.href = 'register.html';
        buttonsBody.appendChild(registerButton);

        const loginButton = document.createElement('button');
        loginButton.textContent = 'Войти';
        loginButton.onclick = () => window.location.href = 'login.html';
        buttonsBody.appendChild(loginButton);
    } else {
        // Показываем роль пользователя
        const roleSpan = document.createElement('span');
        roleSpan.textContent = `Роль: ${userRole} | `;
        buttonsBody.appendChild(roleSpan);

        // Кнопка выхода
        const logoutButton = document.createElement('button');
        logoutButton.textContent = 'Выйти';
        logoutButton.onclick = () => {
            localStorage.removeItem('userRole');
            window.location.reload();
        };
        buttonsBody.appendChild(logoutButton);
    }

    // Показываем/скрываем форму добавления книг
    const addBookForm = document.getElementById('add-book-form');
    addBookForm.style.display = (userRole === 'guest') ? 'none' : 'block';
}
