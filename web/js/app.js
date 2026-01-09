document.addEventListener('DOMContentLoaded', function() {
    const helloBtn = document.getElementById('hello-btn');
    const messageDiv = document.getElementById('message');

    helloBtn.addEventListener('click', function() {
        fetch('/api/hello')
            .then(response => response.json())
            .then(data => {
                messageDiv.textContent = data.message;
                messageDiv.className = 'message success';
            })
            .catch(error => {
                messageDiv.textContent = 'Error: ' + error.message;
                messageDiv.className = 'message error';
            });
    });
});
