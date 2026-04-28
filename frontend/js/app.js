function showPage(pageId) {
    document.querySelectorAll('.page').forEach((page) => page.classList.remove('active'));
    const target = document.getElementById(pageId + 'Page');
    if (target) target.classList.add('active');
}

document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.nav-link').forEach((link) => {
        link.addEventListener('click', (event) => {
            event.preventDefault();
            const page = link.dataset.page;
            if (page) showPage(page);
        });
    });
});