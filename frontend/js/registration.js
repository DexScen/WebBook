document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('register-form');
    const usernameInput = document.getElementById('register-username');
    const emailInput = document.getElementById('register-email');
    const nameInput = document.getElementById('register-name');
    const passwordInput = document.getElementById('register-password');
    const roleSelect = document.getElementById('register-role');

    const style = document.createElement('style');
    style.textContent = `
        .input-valid {
            border: 2px solid green !important;
        }
        .input-invalid {
            border: 2px solid red !important;
        }
        .validation-message {
            font-size: 12px;
            margin-top: -10px;
            margin-bottom: 10px;
        }
        .error-message {
            color: red;
        }
        .success-message {
            color: green;
        }
    `;
    document.head.appendChild(style);

    // Функция для добавления сообщения валидации
    function addValidationMessage(input, message, isValid) {
        let messageDiv = input.nextElementSibling;
        if (!messageDiv || !messageDiv.classList.contains('validation-message')) {
            messageDiv = document.createElement('div');
            messageDiv.className = 'validation-message';
            input.parentNode.insertBefore(messageDiv, input.nextSibling);
        }
        messageDiv.textContent = message;
        messageDiv.className = `validation-message ${isValid ? 'success-message' : 'error-message'}`;
        input.className = isValid ? 'input-valid' : 'input-invalid';
    }

    function isValidEmail(email) {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    }

    function isValidPassword(password) {
        return password.length >= 6;
    }

    // Проверка уникальности логина через AJAX
    let loginCheckTimeout;
    async function checkLoginAvailability(login) {
        try {
            const response = await fetch(`http://localhost:81/api/auth/check-login?login=${encodeURIComponent(login)}`);
            const data = await response.json();
            return data.available;
        } catch (error) {
            console.error('Ошибка при проверке логина:', error);
            return false;
        }
    }

    // Обработчики событий для валидации в реальном времени
    emailInput.addEventListener('input', function() {
        const email = this.value.trim();
        const isValid = isValidEmail(email);
        addValidationMessage(
            this,
            isValid ? 'Email корректен' : 'Введите корректный email адрес',
            isValid
        );
    });

    passwordInput.addEventListener('input', function() {
        const password = this.value;
        const isValid = isValidPassword(password);
        addValidationMessage(
            this,
            isValid ? 'Пароль подходит' : 'Пароль должен содержать минимум 6 символов',
            isValid
        );
    });

    usernameInput.addEventListener('input', function() {
        const login = this.value.trim();
        
        clearTimeout(loginCheckTimeout);
        
        loginCheckTimeout = setTimeout(async () => {
            if (login.length < 3) {
                addValidationMessage(
                    this,
                    'Логин должен содержать минимум 3 символа',
                    false
                );
                return;
            }

            const isAvailable = await checkLoginAvailability(login);
            addValidationMessage(
                this,
                isAvailable ? 'Логин доступен' : 'Логин уже занят',
                isAvailable
            );
        }, 500); // Задержка 500мс перед проверкой
    });

    // Обработчик отправки формы
    form.addEventListener('submit', async function(e) {
        e.preventDefault();

        const login = usernameInput.value.trim();
        const email = emailInput.value.trim();
        const name = nameInput.value.trim();
        const password = passwordInput.value;
        const role = roleSelect.value;

        if (!login || !email || !name || !password || !role) {
            alert('Пожалуйста, заполните все поля');
            return;
        }

        if (!isValidEmail(email)) {
            alert('Пожалуйста, введите корректный email');
            return;
        }

        if (!isValidPassword(password)) {
            alert('Пароль должен содержать минимум 6 символов');
            return;
        }

        // Финальная проверка доступности логина
        const isLoginAvailable = await checkLoginAvailability(login);
        if (!isLoginAvailable) {
            alert('Этот логин уже занят. Пожалуйста, выберите другой.');
            return;
        }

        try {
            const response = await fetch('http://localhost:81/api/auth/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    login,
                    email,
                    name,
                    password,
                    role
                })
            });

            const data = await response.json();
            
            if (data === "success") {
                alert('Регистрация успешна!');
                window.location.href = 'login.html';
            } else {
                alert('Ошибка при регистрации: ' + data);
            }
        } catch (error) {
            console.error('Ошибка:', error);
            alert('Произошла ошибка при регистрации');
        }
    });
});
