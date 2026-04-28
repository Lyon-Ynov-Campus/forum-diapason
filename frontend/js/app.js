// Configuration
const API_BASE_URL = 'http://localhost:8080/api';

// État global
let currentUser = null;
let currentPage = 'home';
let postsOffset = 0;
const POSTS_LIMIT = 20;

// Utilitaires
function showPage(pageId) {
    document.querySelectorAll('.page').forEach(page => page.classList.remove('active'));
    document.getElementById(pageId + 'Page').classList.add('active');
    currentPage = pageId;
}

function showLoading() {
    document.getElementById('loading').style.display = 'block';
}

function hideLoading() {
    document.getElementById('loading').style.display = 'none';
}

function showError(message) {
    alert('Erreur: ' + message);
}

function getAuthToken() {
    return localStorage.getItem('authToken');
}

function setAuthToken(token) {
    localStorage.setItem('authToken', token);
}

function clearAuthToken() {
    localStorage.removeItem('authToken');
}

function updateUserInterface() {
    const userMenu = document.getElementById('userMenu');
    const createPostBtn = document.getElementById('createPostBtn');
    const navLinks = document.querySelectorAll('.nav-link');

    if (currentUser) {
        userMenu.style.display = 'flex';
        document.getElementById('userPseudo').textContent = currentUser.pseudo;
        createPostBtn.style.display = 'inline-flex';

        // Masquer les liens de connexion/inscription
        navLinks.forEach(link => {
            if (link.dataset.page === 'login' || link.dataset.page === 'register') {
                link.style.display = 'none';
            }
        });
    } else {
        userMenu.style.display = 'none';
        createPostBtn.style.display = 'none';

        // Afficher les liens de connexion/inscription
        navLinks.forEach(link => {
            if (link.dataset.page === 'login' || link.dataset.page === 'register') {
                link.style.display = 'inline';
            }
        });
    }
}

// API calls
async function apiRequest(endpoint, options = {}) {
    const url = API_BASE_URL + endpoint;
    const token = getAuthToken();

    const defaultOptions = {
        headers: {
            'Content-Type': 'application/json',
            ...(token && { 'Authorization': `Bearer ${token}` })
        }
    };

    const response = await fetch(url, { ...defaultOptions, ...options });

    if (!response.ok) {
        const error = await response.text();
        throw new Error(error || `HTTP ${response.status}`);
    }

    return response.json();
}

async function loadPosts() {
    try {
        showLoading();
        const posts = await apiRequest(`/posts?limit=${POSTS_LIMIT}&offset=${postsOffset}`);
        displayPosts(posts || []);
        postsOffset += POSTS_LIMIT;

        if (posts && posts.length === POSTS_LIMIT) {
            document.getElementById('loadMoreBtn').style.display = 'block';
        } else {
            document.getElementById('loadMoreBtn').style.display = 'none';
        }
    } catch (error) {
        console.error('Erreur lors du chargement des posts:', error);
        document.getElementById('postsContainer').innerHTML = '<p>Impossible de charger les publications. Assurez-vous que le serveur backend est lancé.</p>';
    } finally {
        hideLoading();
    }
}

function displayPosts(posts) {
    const container = document.getElementById('postsContainer');

    if (!posts || posts.length === 0) {
        container.innerHTML = '<p>Aucune publication pour le moment.</p>';
        return;
    }

    posts.forEach(post => {
        const postElement = createPostElement(post);
        container.appendChild(postElement);
    });
}

function createPostElement(post) {
    const postDiv = document.createElement('div');
    postDiv.className = 'post-card';
    postDiv.innerHTML = `
        <div class="post-header">
            <div class="post-author">
                <img src="${post.author_photo || 'https://via.placeholder.com/32'}" alt="Avatar" onerror="this.src='https://via.placeholder.com/32'">
                <span>${post.author_pseudo}</span>
            </div>
            <div class="post-meta">${new Date(post.created_at).toLocaleDateString('fr-FR')}</div>
        </div>
        <h3 class="post-title">${post.title}</h3>
        <div class="post-content">${post.content}</div>
        <div class="post-actions">
            <div class="post-stats">
                <button class="like-btn" data-post-id="${post.id}">
                    <i class="fas fa-heart"></i>
                    <span>${post.like_count}</span>
                </button>
                <span><i class="fas fa-comment"></i> ${post.comment_count}</span>
            </div>
        </div>
    `;

    // Gestionnaire d'événement pour le like
    const likeBtn = postDiv.querySelector('.like-btn');
    likeBtn.addEventListener('click', () => toggleLike(post.id, likeBtn));

    return postDiv;
}

async function toggleLike(postId, button) {
    if (!currentUser) {
        showPage('login');
        return;
    }

    try {
        const isLiked = button.classList.contains('liked');
        const endpoint = `/posts/${postId}/like`;
        const method = isLiked ? 'DELETE' : 'POST';

        await apiRequest(endpoint, { method });

        const countSpan = button.querySelector('span');
        let count = parseInt(countSpan.textContent);

        if (isLiked) {
            button.classList.remove('liked');
            countSpan.textContent = count - 1;
        } else {
            button.classList.add('liked');
            countSpan.textContent = count + 1;
        }
    } catch (error) {
        showError('Impossible de gérer le like');
        console.error(error);
    }
}

// Gestion des formulaires
async function handleLogin(event) {
    event.preventDefault();

    const identifier = document.getElementById('loginIdentifier').value;
    const password = document.getElementById('loginPassword').value;

    try {
        const response = await apiRequest('/auth/login', {
            method: 'POST',
            body: JSON.stringify({ email_or_pseudo: identifier, password })
        });

        setAuthToken(response.token);
        currentUser = response.user;
        updateUserInterface();
        showPage('home');
        postsOffset = 0;
        document.getElementById('postsContainer').innerHTML = '';
        loadPosts();
        document.getElementById('loginForm').reset();
    } catch (error) {
        showError('Échec de la connexion: ' + error.message);
        console.error(error);
    }
}

async function handleRegister(event) {
    event.preventDefault();

    const email = document.getElementById('registerEmail').value;
    const pseudo = document.getElementById('registerPseudo').value;
    const password = document.getElementById('registerPassword').value;

    try {
        const response = await apiRequest('/auth/register', {
            method: 'POST',
            body: JSON.stringify({ email, pseudo, password })
        });

        setAuthToken(response.token);
        currentUser = response.user;
        updateUserInterface();
        showPage('home');
        postsOffset = 0;
        document.getElementById('postsContainer').innerHTML = '';
        loadPosts();
        document.getElementById('registerForm').reset();
    } catch (error) {
        showError('Échec de l\'inscription: ' + error.message);
        console.error(error);
    }
}

async function handleCreatePost(event) {
    event.preventDefault();

    const title = document.getElementById('postTitle').value;
    const content = document.getElementById('postContent').value;

    try {
        await apiRequest('/posts', {
            method: 'POST',
            body: JSON.stringify({ title, content })
        });

        document.getElementById('createPostForm').reset();
        showPage('home');
        // Recharger les posts
        postsOffset = 0;
        document.getElementById('postsContainer').innerHTML = '';
        loadPosts();
    } catch (error) {
        showError('Impossible de créer la publication: ' + error.message);
        console.error(error);
    }
}

function handleLogout() {
    clearAuthToken();
    currentUser = null;
    updateUserInterface();
    showPage('home');
    postsOffset = 0;
    document.getElementById('postsContainer').innerHTML = '';
    loadPosts();
}

// Initialisation
document.addEventListener('DOMContentLoaded', function() {
    // Vérifier si l'utilisateur est connecté
    const token = getAuthToken();
    if (token) {
        // TODO: Valider le token et récupérer les infos utilisateur
        // Pour l'instant, on suppose qu'il est valide
        currentUser = { pseudo: 'Utilisateur' }; // Placeholder
        updateUserInterface();
    }

    // Gestionnaires d'événements de navigation
    document.querySelectorAll('.nav-link').forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            const page = this.dataset.page;
            showPage(page);
        });
    });

    // Gestionnaires de formulaires
    document.getElementById('loginForm')?.addEventListener('submit', handleLogin);
    document.getElementById('registerForm')?.addEventListener('submit', handleRegister);
    document.getElementById('createPostForm')?.addEventListener('submit', handleCreatePost);

    // Boutons
    document.getElementById('createPostBtn')?.addEventListener('click', () => showPage('createPost'));
    document.getElementById('cancelPostBtn')?.addEventListener('click', () => showPage('home'));
    document.getElementById('logoutBtn')?.addEventListener('click', handleLogout);
    document.getElementById('loadMoreBtn')?.addEventListener('click', loadPosts);

    // Charger les posts initiaux
    loadPosts();
});