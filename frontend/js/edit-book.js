// Функция для получения роли пользователя
function getUserRole() {
    return localStorage.getItem('userRole') || 'guest';
}

// Проверяем права доступа - только администратор может редактировать книги
document.addEventListener('DOMContentLoaded', function() {
    if (getUserRole() !== 'admin') {
        alert('У вас нет прав для редактирования книг');
        window.location.href = 'index.html';
        return;
    }
});

const urlParams = new URLSearchParams(window.location.search);
const bookId = urlParams.get('id');  
const bookPrice = urlParams.get('price'); 
const bookAuthor = urlParams.get('author'); 
const bookTitle = urlParams.get('title'); 

const editForm = document.getElementById('edit-book-form');

async function loadBookData() {
    document.getElementById('title').value = bookTitle;
    document.getElementById('author').value = bookAuthor;
    document.getElementById('price').value = bookPrice;
}

async function updateBook(event) {
    event.preventDefault();
    
    // Повторная проверка прав перед отправкой запроса
    if (getUserRole() !== 'admin') {
        alert('У вас нет прав для редактирования книг');
        window.location.href = 'index.html';
        return;
    }

    const title = document.getElementById('title').value;
    const author = document.getElementById('author').value;
    const price = parseInt(document.getElementById('price').value, 10); 
    const id = parseInt(bookId, 10)

    const updatedBook = { 
        id : id,
        title : title, 
        author : author,
        price: price 
    };

    try {
        const response = await fetch(`http://localhost:8080/books`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(updatedBook),
        });

        if (response.ok) {
            alert('Книга успешно обновлена');
            window.location.href = 'index.html';  
        } else {
            throw new Error('Не удалось обновить книгу');
        }
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Ошибка при обновлении книги');
    }
}

editForm.addEventListener('submit', updateBook);

loadBookData();