document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('login-form').addEventListener('submit', async function(e) {
        e.preventDefault();

        const login = document.getElementById('login-username').value.trim();
        const password = document.getElementById('login-password').value.trim();

        if (!login || !password) {
            alert('Пожалуйста, заполните все поля');
            return;
        }

        try {
            const response = await fetch('http://localhost:81/api/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    login,
                    password
                })
            });

            const data = await response.json();
            
            if (data.role && data.role !== "unauthorized by user" && data.role !== "unauthorized by password") {
                localStorage.setItem('userRole', data.role);
                alert('Вход выполнен успешно!');
                window.location.href = 'index.html';
            } else {
                if (data.role === "unauthorized by user") {
                    alert('Пользователь не найден');
                } else if (data.role === "unauthorized by password") {
                    alert('Неверный пароль');
                } else {
                    alert('Ошибка при входе');
                }
            }
        } catch (error) {
            console.error('Ошибка:', error);
            alert('Произошла ошибка при попытке входа');
        }
    });
});

