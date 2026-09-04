const modal = document.getElementById('newProjectModal');
const newProjectBtn = document.querySelector('.n-btn');
const closeBtn = document.querySelector('.modal-close');
const createBtn = document.querySelector('.create-btn');

newProjectBtn.addEventListener('click', () => {
    modal.classList.add('active');
});

closeBtn.addEventListener('click', () => {
    modal.classList.remove('active');
});

modal.addEventListener('click', (e) => {
    if (e.target === modal) {
        modal.classList.remove('active');
    }
});