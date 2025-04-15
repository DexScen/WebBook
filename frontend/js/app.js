
const booksPerPage = 10;  
let allBooks = [];  
let currentPage = 1; 

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

    booksToDisplay.forEach(book => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${book.title}</td>
            <td>${book.author}</td>
            <td>${book.price}</td>
            <td>
                <button onclick="editBook(${book.id}, '${encodeURIComponent(book.price)}', '${encodeURIComponent(book.author)}', '${encodeURIComponent(book.title)}')">Изменить</button>
                <button onclick="deleteBook(${book.id})">Удалить</button>
            </td>
        `;
        tableBody.appendChild(row);
    });
}


function updatePagination() {
    const paginationContainer = document.getElementById('pagination');
    paginationContainer.innerHTML = '';  

    const totalPages = Math.ceil(allBooks.length / booksPerPage); 

    if (currentPage > 1) {
        const prevPage = document.createElement('button');
        prevPage.innerText = 'Назад';
        prevPage.onclick = () => changePage(currentPage - 1);
        paginationContainer.appendChild(prevPage);
    }

    const pageNumber = document.createElement('span');
    pageNumber.innerText = `Страница ${currentPage} из ${totalPages}`;
    paginationContainer.appendChild(pageNumber);

    if (currentPage < totalPages) {
        const nextPage = document.createElement('button');
        nextPage.innerText = 'Вперёд';
        nextPage.onclick = () => changePage(currentPage + 1);
        paginationContainer.appendChild(nextPage);
    }
}


function changePage(page) {
    if (page < 1 || page > Math.ceil(allBooks.length / booksPerPage)) return;
    currentPage = page;
    renderBooks(currentPage);  
    updatePagination();  
}


function displayBooks(books) {
    const tableBody = document.getElementById('book-list');
    tableBody.innerHTML = '';  

    books.forEach(book => {
        const row = document.createElement('tr');

        const titleCell = document.createElement('td');
        titleCell.textContent = book.title;
        row.appendChild(titleCell);

        const authorCell = document.createElement('td');
        authorCell.textContent = book.author;
        row.appendChild(authorCell);

        const priceCell = document.createElement('td');
        priceCell.textContent = book.price;
        row.appendChild(priceCell);

        const actionsCell = document.createElement('td');

        const editButton = document.createElement('button');
        editButton.textContent = 'Изменить';
        editButton.onclick = () => editBook(book.id, book.price, book.author, book.title);
        actionsCell.appendChild(editButton);

        const deleteButton = document.createElement('button');
        deleteButton.textContent = 'Удалить';
        deleteButton.onclick = () => deleteBook(book.id);
        actionsCell.appendChild(deleteButton);

        row.appendChild(actionsCell);
        tableBody.appendChild(row);
    });
}

function editBook(bookId, bookPrice, bookAuthor, bookTitle) {
    window.location.href = `edit-book.html?id=${bookId}&price=${bookPrice}&author=${bookAuthor}&title=${bookTitle}`; 
}

async function deleteBook(bookId) {
    const id = {
        id : bookId
    };
    const confirmDelete = confirm('Вы уверены, что хотите удалить эту книгу?');
    if (confirmDelete) {
        try {
            const response = await fetch(`http://localhost:8080/books`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(id),
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
}

document.getElementById('add-book-form').addEventListener('submit', async (e) => {
    e.preventDefault();  

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

fetchBooks();