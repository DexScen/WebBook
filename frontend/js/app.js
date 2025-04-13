async function fetchBooks() {
    try {
        const response = await fetch('http://localhost:8080/books'); 
        if (!response.ok) {
            throw new Error('Не удалось загрузить данные');
        }
        const books = await response.json();
        displayBooks(books);
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Ошибка при загрузке данных');
    }
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

        const yearCell = document.createElement('td');
        yearCell.textContent = book.year;
        row.appendChild(yearCell);

        tableBody.appendChild(row);  
    });
}

fetchBooks(); 