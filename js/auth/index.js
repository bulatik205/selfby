let timerInterval;
let countdown = 10;

function showError(message) {
    const errorDiv = document.getElementById('errorMessage');
    const errorText = document.getElementById('errorText');
    const timerBox = document.getElementById('timerBox');

    errorText.textContent = message;
    errorDiv.style.display = 'flex';
    timerBox.style.display = 'flex';

    clearInterval(timerInterval);
    countdown = 10;
    document.getElementById('timerCount').textContent = countdown;

    timerInterval = setInterval(() => {
        countdown--;
        document.getElementById('timerCount').textContent = countdown;

        if (countdown <= 0) {
            clearInterval(timerInterval);
            hideAllMessages();
        }
    }, 1000);
}

function hideError() {
    clearInterval(timerInterval);
    hideAllMessages();
}

function hideAllMessages() {
    const errorDiv = document.getElementById('errorMessage');
    const timerBox = document.getElementById('timerBox');

    errorDiv.style.display = 'none';
    timerBox.style.display = 'none';

    clearInterval(timerInterval);
}

const urlParams = new URLSearchParams(window.location.search);
const error = urlParams.get('error');

if (error) {
    showError(error);

    window.history.replaceState({}, document.title, window.location.pathname);
}